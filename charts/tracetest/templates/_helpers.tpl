{{/*
Expand the name of the chart.
*/}}
{{- define "tracetest.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "tracetest.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "tracetest.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "tracetest.labels" -}}
helm.sh/chart: {{ include "tracetest.chart" . }}
{{ include "tracetest.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Create a default fully qualified postgresql name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "tracetest.postgresql.fullname" -}}
{{- include "common.names.dependency.fullname" (dict "chartName" "postgresql" "chartValues" .Values.postgresql "context" $) -}}
{{- end -}}


{{- define "tracetest.database.encryptedPassword" -}}
  {{- include "tracetest.database.rawPassword" . | b64enc | quote -}}
{{- end -}}

{{- define "tracetest.database.rawPassword" -}}
{{- if .Values.postgresql.enabled }}
    {{- if .Values.global.postgresql }}
        {{- if .Values.global.postgresql.auth }}
            {{- coalesce .Values.global.postgresql.auth.postgresPassword .Values.postgresql.auth.postgresPassword -}}
        {{- else -}}
            {{- .Values.postgresql.auth.postgresPassword -}}
        {{- end -}}
    {{- else -}}
        {{- .Values.postgresql.auth.postgresPassword -}}
    {{- end -}}
{{- else -}}
    {{- .Values.externalDatabase.password -}}
{{- end -}}
{{- end -}}


{{/*
Selector labels
*/}}
{{- define "tracetest.selectorLabels" -}}
app.kubernetes.io/name: {{ include "tracetest.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "tracetest.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "tracetest.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Define a valid pathPrefix
*/}}
{{- define "tracetest.pathPrefix" -}}
    {{- if .Values.server.pathPrefix }}
        {{- .Values.server.pathPrefix }}
    {{- else }}
        {{- "/" }}
    {{- end }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "otelCollector.labels" -}}
helm.sh/chart: {{ include "tracetest.chart" . }}
{{ include "otelCollector.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "otelCollector.selectorLabels" -}}
app.kubernetes.io/name: {{ default "otel-collector" .Values.otelCollector.name }}
{{- end -}}