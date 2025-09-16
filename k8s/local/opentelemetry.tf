# opentelemetry.tf

resource "helm_release" "opentelemetry_collector" {
  name             = "opentelemetry-collector"
  repository       = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  chart            = "opentelemetry-collector"
  version          = "0.85.0"
  namespace        = "monitoring"
  create_namespace = true

  values = [
    file("${path.module}/values/opentelemetry-collector-values.yaml")
  ]

  # The Collector needs Tempo to be running before it starts
  depends_on = [
    helm_release.tempo
  ]
}
