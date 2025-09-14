# prometheus.tf

variable "grafana_admin_user" {
  description = "The admin username for Grafana."
  type        = string
}

variable "grafana_admin_password" {
  description = "The admin password for Grafana."
  type        = string
  sensitive   = true # This prevents the value from being shown in logs
}

resource "kubernetes_secret" "grafana_credentials" {
  metadata {
    name      = "grafana-admin-credentials"
    namespace = "monitoring"
  }

  # The 'data' field will be automatically base64 encoded by Terraform.
  data = {
    admin_user     = var.grafana_admin_user
    admin_password = var.grafana_admin_password
  }

}

resource "helm_release" "prometheus_stack" {
  name             = "prometheus-stack"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  version          = "51.6.0"
  namespace        = "monitoring"
  create_namespace = true

  values = [
    file("${path.module}/values/prometheus-values.yaml")
  ]
  # Add this depends_on block
  depends_on = [
    kubernetes_secret.grafana_credentials
  ]
}
