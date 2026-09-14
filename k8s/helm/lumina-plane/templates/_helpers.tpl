{{- define "lumina-plane.name" -}}
{{- default .Chart.Name .Values.nameOverride -}}
{{- end -}}

{{- define "lumina-plane.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride }}
{{- else }}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end -}}

{{- define "lumina-plane.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "." "-" }}
app.kubernetes.io/name: {{ include "lumina-plane.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "lumina-plane.selectorLabels" -}}
app.kubernetes.io/name: {{ include "lumina-plane.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
