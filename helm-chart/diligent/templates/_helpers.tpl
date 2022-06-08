{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "diligent.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "diligent.commonLabels" -}}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ include "diligent.chart" . }}
{{- end }}

{{/*
Define name for the boss app
*/}}
{{- define "diligent.boss.name" -}}
{{- printf "%s-diligent-boss" .Release.Name | trunc 63 }}
{{- end }}

{{/*
Boss selector labels
*/}}
{{- define "diligent.boss.selectorLabels" -}}
app.kubernetes.io/name: {{ include "diligent.boss.name" . }}
{{- end }}

{{/*
Define name for the minion app
*/}}
{{- define "diligent.minion.name" -}}
{{- printf "%s-diligent-minion" .Release.Name | trunc 63 }}
{{- end }}

{{/*
Minion selector labels
*/}}
{{- define "diligent.minion.selectorLabels" -}}
app.kubernetes.io/name: {{ include "diligent.minion.name" . }}
{{- end }}

{{/*
Define name for the prometheus app
*/}}
{{- define "diligent.prometheus.name" -}}
{{- printf "%s-diligent-prometheus" .Release.Name | trunc 63 }}
{{- end }}

{{/*
Prometheus selector labels
*/}}
{{- define "diligent.prometheus.selectorLabels" -}}
app.kubernetes.io/name: {{ include "diligent.prometheus.name" . }}
{{- end }}

{{/*
Define name for the grafana app
*/}}
{{- define "diligent.grafana.name" -}}
{{- printf "%s-diligent-grafana" .Release.Name | trunc 63 }}
{{- end }}

{{/*
Grafana selector labels
*/}}
{{- define "diligent.grafana.selectorLabels" -}}
app.kubernetes.io/name: {{ include "diligent.grafana.name" . }}
{{- end }}
