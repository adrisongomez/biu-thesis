# opentelemetry.tf

resource "helm_release" "opentelemetry_collector_custom" {
  name             = "opentelemetry-collector-custom"
  repository       = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart            = "opentelemetry-collector"
  version          = "0.85.0"
  namespace        = "monitoring"
  create_namespace = true
  cleanup_on_fail  = true

  values = [
    file("${path.module}/values/opentelemetry-collector-values.custom.yaml")
  ]

  # The Collector needs Tempo to be running before it starts
  depends_on = [
    kubernetes_deployment_v1.custom_collector_deployment
  ]
}
