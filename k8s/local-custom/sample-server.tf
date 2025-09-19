resource "kubernetes_namespace_v1" "sample_namespaces" {
  metadata {
    name = "sample"
  }
}

resource "kubernetes_deployment_v1" "sample_service_custom_deployment" {
  depends_on = [kubernetes_namespace_v1.sample_namespaces]
  metadata {
    namespace = "sample"
    name      = "sample-service-custom"
    labels = {
      app = "sample"
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "sample"
      }
    }
    template {
      metadata {
        labels = {
          app = "sample"
        }
      }
      spec {
        container {
          name  = "sample-service-custom"
          image = "adrisongomez/server-custom:latest"
          port {
            container_port = 5000
          }
          liveness_probe {
            http_get {
              path = "/api/healthcheck"
              port = 5000
            }
            initial_delay_seconds = 5
            period_seconds        = 20
            timeout_seconds       = 5
            failure_threshold     = 3
          }
          readiness_probe {
            http_get {
              path = "/api/healthcheck"
              port = 5000
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service_v1" "sample_service_custom_svc" {
  metadata {
    name      = "sample-service-custom"
    namespace = "sample"
  }
  spec {
    type = "ClusterIP"
    selector = {
      app = "sample"
    }
    port {
      port        = 5000
      target_port = 5000
      protocol    = "TCP"
    }
  }
  depends_on = [kubernetes_namespace_v1.sample_namespaces]
}
