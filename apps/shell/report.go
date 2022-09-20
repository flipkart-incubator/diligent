package main

import (
	"context"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/flipkart-incubator/diligent/pkg/proto"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"io"
	"math"
	"os"
	"strings"
	"text/template"
	"time"
)

func init() {
	reportCmd := &grumble.Command{
		Name:    "report",
		Help:    "work with reports",
		Aliases: []string{"re"},
	}
	grumbleApp.AddCommand(reportCmd)

	reportJobCmd := &grumble.Command{
		Name: "job",
		Help: "generate a report for the current job",
		Run:  reportJob,
	}
	reportCmd.AddCommand(reportJobCmd)

	reportExpCmd := &grumble.Command{
		Name: "experiment",
		Help: "generate a report for the current experiment",
		Run:  reportExperiment,
	}
	reportCmd.AddCommand(reportExpCmd)
}

func reportJob(c *grumble.Context) error {
	jobName, startTime, endTime, stepSize, err := getJobParams(c)
	if err != nil {
		return err
	}

	chs, err := makePanels(c, startTime, endTime, stepSize)
	if err != nil {
		return err
	}
	now := time.Now()
	fileName := fmt.Sprintf("%s-job-report-%04d-%02d-%02dT%02d%02d%02d.html", jobName,
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	err = saveReport(c, chs, fileName)

	return err
}

func reportExperiment(c *grumble.Context) error {
	expName, startTime, endTime, stepSize, err := getExperimentParams(c)
	if err != nil {
		return err
	}

	chs, err := makePanels(c, startTime, endTime, stepSize)
	if err != nil {
		return err
	}

	now := time.Now()
	fileName := fmt.Sprintf("%s-experiment-report-%04d-%02d-%02dT%02d%02d%02d.html", expName,
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	err = saveReport(c, chs, fileName)

	return err
}

func getJobParams(c *grumble.Context) (jobName string, startTime, endTime time.Time, stepSize time.Duration, err error) {
	c.App.Printf("Getting info about current job...\n")

	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	res, err := bossClient.GetJobInfo(grpcCtx, &proto.BossGetJobInfoRequest{})
	grpcCancel()
	if err != nil {
		return
	}

	if res.GetJobInfo() == nil {
		err = fmt.Errorf("no current job. unable to generate report")
		return
	}

	jobName = res.GetJobInfo().GetJobSpec().GetJobName()

	switch res.GetJobInfo().GetJobState() {
	case proto.JobState_NEW, proto.JobState_PREPARED, proto.JobState_RUNNING:
		err = fmt.Errorf("unable to generate report. job has not yet ended")
		return
	case proto.JobState_ENDED_SUCCESS:
		c.App.Printf("Job ended successfully\n")
	case proto.JobState_ENDED_FAILURE:
		c.App.Printf("Warning: job ended with failure\n")
	case proto.JobState_ENDED_ABORTED:
		c.App.Printf("Warning: job was aborted\n")
	}

	// Capture metrics from 2 mins before and after with 1min rounding
	startTime = time.UnixMilli(res.GetJobInfo().GetRunTime()).Add(-2 * time.Minute).Round(1 * time.Minute)
	endTime = time.UnixMilli(res.GetJobInfo().GetEndTime()).Add(2 * time.Minute).Round(1 * time.Minute)

	if startTime.After(endTime) {
		err = fmt.Errorf("calculated startTime %s is before end time %s", startTime.Format(time.UnixDate), endTime.Format(time.UnixDate))
		return
	}
	stepSize = 10 * time.Second

	c.App.Printf("job-name     : %s\n", jobName)
	c.App.Printf("start-time   : %s\n", startTime.Format(time.UnixDate))
	c.App.Printf("end-time     : %s\n", endTime.Format(time.UnixDate))
	c.App.Printf("step-duration: %s\n", stepSize.String())

	return
}

func getExperimentParams(c *grumble.Context) (expName string, startTime, endTime time.Time, stepSize time.Duration, err error) {
	c.App.Printf("Getting info about current experiment...\n")

	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	res, err := bossClient.GetExperimentInfo(grpcCtx, &proto.BossGetExperimentInfoRequest{})
	grpcCancel()
	if err != nil {
		return
	}

	if res.GetExperimentInfo() == nil {
		err = fmt.Errorf("no current experiment. unable to generate report")
		return
	}

	expName = res.GetExperimentInfo().GetName()

	switch res.GetExperimentInfo().GetState() {
	case proto.ExperimentState_NEW_EXPERIMENT, proto.ExperimentState_STARTED:
		err = fmt.Errorf("unable to generate report. experiment has not yet ended")
		return
	case proto.ExperimentState_STOPPED:
		c.App.Printf("Experiment has stopped\n")
	}

	// Capture metrics from 2 mins before and after with 1min rounding
	startTime = time.UnixMilli(res.GetExperimentInfo().GetStartTime()).Add(-2 * time.Minute).Round(1 * time.Minute)
	endTime = time.UnixMilli(res.GetExperimentInfo().GetStopTime()).Add(2 * time.Minute).Round(1 * time.Minute)

	if startTime.After(endTime) {
		err = fmt.Errorf("calculated startTime %s is before end time %s", startTime.Format(time.UnixDate), endTime.Format(time.UnixDate))
		return
	}
	stepSize = 10 * time.Second

	c.App.Printf("experiment-name: %s\n", expName)
	c.App.Printf("start-time     : %s\n", startTime.Format(time.UnixDate))
	c.App.Printf("end-time       : %s\n", endTime.Format(time.UnixDate))
	c.App.Printf("step-duration  : %s\n", stepSize.String())

	return
}

func makePanels(c *grumble.Context, startTime, endTime time.Time, stepSize time.Duration) ([]components.Charter, error) {
	c.App.Printf("Getting data for panels and creating charts...\n")
	chs := make([]components.Charter, 0)
	for _, p := range panels {
		ch := newLineChart(p, startTime, endTime)
		for _, q := range p.queries {
			metrics, err := promQuery(c, q.query, startTime, endTime, stepSize)
			if err != nil {
				return nil, err
			}
			err = plotData(ch, metrics, q)
			if err != nil {
				return nil, err
			}
		}
		chs = append(chs, ch)
	}
	return chs, nil
}

func promQuery(c *grumble.Context, query string, startTime, endTime time.Time, step time.Duration) (model.Matrix, error) {
	promAddr := c.Flags.String("prom")

	client, err := api.NewClient(api.Config{Address: promAddr})
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus api client: %v", err)
	}

	promAPI := v1.NewAPI(client)

	value, _, err := promAPI.QueryRange(context.Background(), query, v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  step,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query prometheus api: %v", err)
	}

	metrics, ok := value.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("unsupported result format: %s", value.Type().String())
	}

	return metrics, nil
}

func newLineChart(panel Panel, startTime, endTime time.Time) *charts.Line {
	lineChart := charts.NewLine()
	lineChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: panel.title}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Time",
			Type: "time",
			Min:  float64(startTime.UnixMilli()),
			Max:  float64(endTime.UnixMilli()),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: panel.yAxisLabel,
			Type: "value",
			Min:  0,
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "axis",
			TriggerOn: "mousemove|click",
		}),
	)
	return lineChart
}

func plotData(chart *charts.Line, metrics model.Matrix, q Query) error {
	xScale := 1.0
	yScale := 1.0
	if q.xScale != 0 {
		xScale = q.xScale
	}
	if q.yScale != 0 {
		yScale = q.yScale
	}

	for _, series := range metrics {
		items := make([]opts.LineData, 0)
		for _, s := range series.Values {
			t := float64(s.Timestamp) * xScale
			v := float64(s.Value) * yScale
			v = math.Round(v*100) / 100
			if math.IsNaN(v) {
				continue
			}
			items = append(items, opts.LineData{Value: []float64{t, v}})
		}

		t := template.Must(template.New("t").Parse(q.legendText))
		var b strings.Builder
		err := t.Execute(&b, toMap(series.Metric))
		if err != nil {
			return err
		}
		chart.AddSeries(b.String(), items)
	}
	return nil
}

func saveReport(c *grumble.Context, chs []components.Charter, fileName string) error {
	c.App.Printf("Generating and saving report...\n")
	page := components.NewPage()
	page.AddCharts(chs...)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = page.Render(io.MultiWriter(f))
	if err != nil {
		return err
	}
	c.App.Printf("Report saved: %s\n", fileName)
	return nil
}

type Query struct {
	query      string
	legendText string
	xScale     float64
	yScale     float64
}

type Panel struct {
	title      string
	yAxisLabel string
	queries    []Query
}

var panels = []Panel{
	{
		title:      "Statement Count",
		yAxisLabel: "Count",
		queries: []Query{
			{
				query:      "sum by (statement) (diligent_statement_duration_seconds_count)",
				legendText: "{{.statement}}",
			},
		},
	},
	{
		title:      "Statement Rate",
		yAxisLabel: "Count/Second",
		queries: []Query{
			{
				query:      "sum by (statement) (rate(diligent_statement_duration_seconds_count[1m]))",
				legendText: "{{.statement}}",
			},
		},
	},
	{
		title:      "Statement Failure Counts",
		yAxisLabel: "Count",
		queries: []Query{
			{
				query:      "sum by (statement) (diligent_statement_failure)",
				legendText: "{{.statement}}",
			},
		},
	},
	{
		title:      "Statement Duration Mean",
		yAxisLabel: "Duration(ms)",
		queries: []Query{
			{
				query:      "(sum by (statement) (rate(diligent_statement_duration_seconds_sum[1m]))) / (sum by (statement) (rate(diligent_statement_duration_seconds_count[1m])))",
				legendText: "{{.statement}}",
				yScale:     1000,
			},
		},
	},
	{
		title:      "Statement Duration p99",
		yAxisLabel: "Duration(ms)",
		queries: []Query{
			{
				query:      "histogram_quantile(0.99, sum by (le, statement) (rate(diligent_statement_duration_seconds_bucket[1m])))",
				legendText: "{{.statement}}",
				yScale:     1000,
			},
		},
	},
	{
		title:      "Statement Duration p999",
		yAxisLabel: "Duration(ms)",
		queries: []Query{
			{
				query:      "histogram_quantile(0.999, sum by (le, statement) (rate(diligent_statement_duration_seconds_bucket[1m])))",
				legendText: "{{.statement}}",
				yScale:     1000,
			},
		},
	},
	{
		title:      "Transaction Rate",
		yAxisLabel: "Count/Second",
		queries: []Query{
			{
				query:      "sum(rate(diligent_transaction_duration_seconds_count[1m]))",
				legendText: "TPS",
			},
		},
	},
	{
		title:      "Transaction Duration",
		yAxisLabel: "Duration(ms)",
		queries: []Query{
			{
				query:      "(sum (rate(diligent_transaction_duration_seconds_sum[1m]))) / (sum (rate(diligent_transaction_duration_seconds_count[1m])))",
				legendText: "Mean",
				yScale:     1000,
			},
			{
				query:      "histogram_quantile(0.99, sum by (le) (rate(diligent_transaction_duration_seconds_bucket[1m])))",
				legendText: "p99",
				yScale:     1000,
			},
		},
	},
	{
		title:      "Statements Per Transaction",
		yAxisLabel: "Count",
		queries: []Query{
			{
				query:      "sum(rate(diligent_statement_duration_seconds_count{statement=\"begin\"}[1m])) / sum(rate(diligent_transaction_duration_seconds_count[1m]))",
				legendText: "Begin",
			},
			{
				query:      "sum(rate(diligent_statement_duration_seconds_count{statement=\"commit\"}[1m])) / sum(rate(diligent_transaction_duration_seconds_count[1m]))",
				legendText: "Commit",
			},
			{
				query:      "sum(rate(diligent_statement_duration_seconds_count{statement=\"select\"}[1m])) / sum(rate(diligent_transaction_duration_seconds_count[1m]))",
				legendText: "Select",
			},
			{
				query:      "sum(rate(diligent_statement_duration_seconds_count{statement=\"insert\"}[1m])) / sum(rate(diligent_transaction_duration_seconds_count[1m]))",
				legendText: "Insert",
			},
			{
				query:      "sum(rate(diligent_statement_duration_seconds_count{statement=\"update\"}[1m])) / sum(rate(diligent_transaction_duration_seconds_count[1m]))",
				legendText: "Update",
			},
			{
				query:      "sum(rate(diligent_statement_duration_seconds_count{statement=\"delete\"}[1m])) / sum(rate(diligent_transaction_duration_seconds_count[1m]))",
				legendText: "Delete",
			},
			{
				query:      "sum(rate(diligent_statement_duration_seconds_count[1m])) / sum(rate(diligent_transaction_duration_seconds_count[1m]))",
				legendText: "All Types",
			},
		},
	},
	{
		title:      "Concurrency",
		yAxisLabel: "Count",
		queries: []Query{
			{
				query:      "sum(diligent_config_concurrency)",
				legendText: "Configured",
			},
			{
				query:      "sum(diligent_concurrency)",
				legendText: "Actual",
			},
		},
	},
	{
		title:      "Transactions Enabled",
		yAxisLabel: "True/False",
		queries: []Query{
			{
				query:      "diligent_config_transaction_enabled",
				legendText: "{{.instance}}",
			},
		},
	},
	{
		title:      "Batch Size",
		yAxisLabel: "Size",
		queries: []Query{
			{
				query:      "diligent_config_batch_size",
				legendText: "{{.instance}}",
			},
		},
	},
	{
		title:      "DB Connections",
		yAxisLabel: "Connections",
		queries: []Query{
			{
				query:      "sum by (label) (diligent_db_connections)",
				legendText: "{{.label}}",
			},
		},
	},
}

func toMap(metric model.Metric) map[string]string {
	m := make(map[string]string)
	for k, v := range metric {
		m[string(k)] = string(v)
	}
	return m
}
