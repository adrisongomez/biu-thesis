# prometheus.tf

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
}