# tempo.tf

resource "helm_release" "tempo" {
  name       = "tempo"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "tempo"
  version    = "1.5.0"
  namespace  = "monitoring"
  create_namespace = true
  depends_on = [
    helm_release.prometheus_stack
  ]
  values = [
    file("${path.module}/values/tempo-values.yaml")
  ]
}