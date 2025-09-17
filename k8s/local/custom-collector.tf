resource "kubernetes_namespace_v1" "custom_collector_ns" {
  metadata {
    name = "custom-collector"
  }
}

resource "kubernetes_deployment_v1" "custom_collector_deployment" {
  depends_on = [kubernetes_namespace_v1.custom_collector_ns]

  metadata {
    namespace = "custom-collector"
    name      = "custom-collector"
    labels = {
      app = "custom-collector"
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "custom-collector"
      }
    }
    template {
      metadata {
        labels = {
          app = "custom-collector"
        }
      }
      spec {
        container {
          name  = "custom-collector"
          image = "adrisongomez/custom-collector:latest"
          port {
            container_port = 4317
            host_port      = 4317
          }
        }
      }
    }
  }
}

resource "kubernetes_service_v1" "custom_collector_svc" {
  metadata {
    name      = "custom-collector"
    namespace = "custom-collector"
  }
  spec {
    type = "ClusterIP"
    selector = {
      app = "custom-collector"
    }
    port {
      port        = 4317
      target_port = 4317
      protocol    = "TCP"
    }
  }
  depends_on = [kubernetes_namespace_v1.custom_collector_ns]
}
