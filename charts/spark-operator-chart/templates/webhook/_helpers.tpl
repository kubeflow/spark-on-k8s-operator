{{/*
Create the name of the cert-manager issuer
*/}}
{{- define "cert-manager.issuerName" -}}
{{ include "spark-operator.fullname" . }}-issuer
{{- end -}}

{{/*
Create the name of the cert-manager certificate
*/}}
{{- define "cert-manager.certificateName" -}}
{{ include "spark-operator.fullname" . }}-cert
{{- end -}}

{{/*
Create the name of the secret to be used by webhook
*/}}
{{- define "spark-operator.webhookSecretName" -}}
{{- include "spark-operator.fullname" . }}-webhook-tls
{{- end -}}

{{/*
Create the name of the service to be used by webhook
*/}}
{{- define "spark-operator.webhookServiceName" -}}
{{- include "spark-operator.fullname" . }}-webhook-svc
{{- end -}}

{{/*
Create the name of webhook configuration
*/}}
{{- define "spark-operator.webhookName" -}}
{{- include "spark-operator.fullname" . }}-webhook
{{- end -}}

{{/*
Create the name of mutating webhook configuration
*/}}
{{- define "spark-operator.mutatingWebhookConfigurationName" -}}
webhook.sparkoperator.k8s.io
{{- end -}}

{{/*
Create the name of mutating webhook configuration
*/}}
{{- define "spark-operator.validatingWebhookConfigurationName" -}}
quotaenforcer.sparkoperator.k8s.io
{{- end -}}
