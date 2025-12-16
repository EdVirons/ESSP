{{/*
Expand the name of the chart.
*/}}
{{- define "essp.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "essp.fullname" -}}
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
{{- define "essp.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "essp.labels" -}}
helm.sh/chart: {{ include "essp.chart" . }}
{{ include "essp.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "essp.selectorLabels" -}}
app.kubernetes.io/name: {{ include "essp.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "essp.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "essp.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Database connection string
*/}}
{{- define "essp.databaseUrl" -}}
{{- printf "postgresql://%s:%s@%s:%d/%s?sslmode=%s" .Values.secrets.database.username .Values.secrets.database.password .Values.database.host (.Values.database.port | int) .Values.database.name .Values.database.sslMode }}
{{- end }}

{{/*
Redis connection string
*/}}
{{- define "essp.redisUrl" -}}
{{- if .Values.secrets.redis.password }}
{{- printf "redis://:%s@%s:%d/%d" .Values.secrets.redis.password .Values.redis.host (.Values.redis.port | int) (.Values.redis.db | int) }}
{{- else }}
{{- printf "redis://%s:%d/%d" .Values.redis.host (.Values.redis.port | int) (.Values.redis.db | int) }}
{{- end }}
{{- end }}

{{/*
MinIO endpoint URL
*/}}
{{- define "essp.minioEndpoint" -}}
{{- if .Values.minio.useSSL }}
{{- printf "https://%s" .Values.minio.endpoint }}
{{- else }}
{{- printf "http://%s" .Values.minio.endpoint }}
{{- end }}
{{- end }}

{{/*
Common environment variables for all services
*/}}
{{- define "essp.commonEnv" -}}
- name: DB_HOST
  value: {{ .Values.database.host | quote }}
- name: DB_PORT
  value: {{ .Values.database.port | quote }}
- name: DB_NAME
  value: {{ .Values.database.name | quote }}
- name: DB_USER
  valueFrom:
    secretKeyRef:
      name: {{ include "essp.fullname" . }}-secrets
      key: db-username
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "essp.fullname" . }}-secrets
      key: db-password
- name: REDIS_HOST
  value: {{ .Values.redis.host | quote }}
- name: REDIS_PORT
  value: {{ .Values.redis.port | quote }}
- name: REDIS_DB
  value: {{ .Values.redis.db | quote }}
{{- if .Values.secrets.redis.password }}
- name: REDIS_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "essp.fullname" . }}-secrets
      key: redis-password
{{- end }}
- name: NATS_URL
  value: {{ .Values.nats.url | quote }}
{{- if .Values.secrets.nats.username }}
- name: NATS_USERNAME
  valueFrom:
    secretKeyRef:
      name: {{ include "essp.fullname" . }}-secrets
      key: nats-username
- name: NATS_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "essp.fullname" . }}-secrets
      key: nats-password
{{- end }}
- name: MINIO_ENDPOINT
  value: {{ .Values.minio.endpoint | quote }}
- name: MINIO_BUCKET
  value: {{ .Values.minio.bucket | quote }}
- name: MINIO_USE_SSL
  value: {{ .Values.minio.useSSL | quote }}
- name: MINIO_ACCESS_KEY
  valueFrom:
    secretKeyRef:
      name: {{ include "essp.fullname" . }}-secrets
      key: minio-access-key
- name: MINIO_SECRET_KEY
  valueFrom:
    secretKeyRef:
      name: {{ include "essp.fullname" . }}-secrets
      key: minio-secret-key
{{- end }}

{{/*
Image pull policy helper
*/}}
{{- define "essp.imagePullPolicy" -}}
{{- if eq .tag "latest" }}
{{- "Always" }}
{{- else }}
{{- default "IfNotPresent" .pullPolicy }}
{{- end }}
{{- end }}
