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
		Aliases: []string{"bs"},
	}
	grumbleApp.AddCommand(reportCmd)

	reportSaveCmd := &grumble.Command{
		Name: "save",
		Help: "save the report for a job",
		Run:  reportSave,
	}
	reportCmd.AddCommand(reportSaveCmd)
}

func reportSave(c *grumble.Context) error {
	promAddr := c.Flags.String("prom")

	c.App.Printf("Discovering timespan for current job...")
	startTime, endTime, stepSize, err := getTimes(c)
	if err != nil {
		return err
	}
	c.App.Printf("Proceeding with:\n")
	c.App.Printf("start-time: %s\n", startTime.Format(time.UnixDate))
	c.App.Printf("end-time: %s\n", endTime.Format(time.UnixDate))
	c.App.Printf("step-duration: %s\n", stepSize.String())

	chs := make([]components.Charter, 0)
	for _, p := range panels {
		ch := newLineChart(p, startTime, endTime)
		for _, q := range p.queries {
			metrics, err := promQuery(promAddr, q.query, startTime, endTime, stepSize)
			if err != nil {
				return err
			}
			err = plotData(ch, metrics, q)
			if err != nil {
				return err
			}
		}
		chs = append(chs, ch)
	}

	page := components.NewPage()
	page.AddCharts(chs...)
	fileName := fmt.Sprintf("report.html")
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

func getTimes(c *grumble.Context) (startTime, endTime time.Time, stepSize time.Duration, err error) {
	bossAddr := c.Flags.String("boss")
	bossClient, err := getBossClient(bossAddr)
	if err != nil {
		return
	}

	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), bossRequestTimeoutSecs*time.Second)
	res, err := bossClient.QueryJob(grpcCtx, &proto.BossQueryJobRequest{})
	grpcCancel()
	if err != nil {
		return
	}

	unreachableCount := 0
	endedCount := 0
	notSuccess := 0
	withoutJobInfoCount := 0
	for _, ji := range res.GetMinionJobInfos() {
		if !ji.GetStatus().GetIsOk() {
			unreachableCount++
			continue
		}
		if ji.GetJobInfo() == nil {
			withoutJobInfoCount++
			continue
		}
		switch ji.GetJobInfo().GetJobState() {
		case proto.JobState_ENDED_SUCCESS:
			endedCount++
		case proto.JobState_ENDED_FAILURE:
			endedCount++
			notSuccess++
		case proto.JobState_ENDED_ABORTED:
			endedCount++
			notSuccess++
		case proto.JobState_ENDED_NEVER_RAN:
			endedCount++
			notSuccess++
		}
	}
	totalMinions := len(res.GetMinionJobInfos())
	notEnded := totalMinions - (unreachableCount + endedCount)
	fmt.Printf("Total minions: %d\n", totalMinions)
	fmt.Printf("Unreachable minions: %d\n", unreachableCount)
	fmt.Printf("Ended minions: %d\n", endedCount)
	fmt.Printf("Failed minions: %d\n", notSuccess)
	fmt.Printf("Minions without valid job info: %d\n", withoutJobInfoCount)
	fmt.Printf("Minions not yet ended: %d\n", notEnded)

	if notEnded > 0 {
		err = fmt.Errorf("unable to generate report. %d minions have not clearly eneded", notEnded)
		return
	}
	if withoutJobInfoCount != 0 {
		err = fmt.Errorf("some minions do not have a valid job info")
		return
	}

	startTime = time.UnixMilli(math.MaxInt64)
	endTime = time.UnixMilli(0)
	for _, ji := range res.GetMinionJobInfos() {
		if ji.GetJobInfo() == nil {
			continue
		}

		st := time.UnixMilli(ji.GetJobInfo().GetRunTime())
		et := time.UnixMilli(ji.GetJobInfo().GetEndTime())

		if st.Before(startTime) {
			startTime = st
		}
		if et.After(endTime) {
			endTime = et
		}
	}
	// Capture metrics from 2 mins before and after with 1min rounding
	startTime = startTime.Add(-2 * time.Minute)
	startTime = startTime.Round(1 * time.Minute)
	endTime = endTime.Add(2 * time.Minute)
	endTime = endTime.Round(1 * time.Minute)
	if startTime.After(endTime) {
		err = fmt.Errorf("calculated startTime %s is before end time %s", startTime.Format(time.UnixDate), endTime.Format(time.UnixDate))
		return
	}
	stepSize = 10 * time.Second
	return
}

func promQuery(server, query string, startTime, endTime time.Time, step time.Duration) (model.Matrix, error) {
	client, err := api.NewClient(api.Config{Address: server})
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
			TriggerOn: "mouseclick",
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
