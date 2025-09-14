resource "kubernetes_namespace" "neo4j_ns" {
  metadata {
    name = "neo4j"
  }
}

resource "helm_release" "neo4j" {
  name       = "neo4j-database"
  repository = "https://helm.neo4j.com/neo4j"
  chart      = "neo4j"
  version    = "5.12.0"
  namespace = kubernetes_namespace.neo4j_ns.metadata[0].name
  atomic   = true

  values = [
    file("${path.module}/values/neo4j-values.base.yaml"),
    file("${path.module}/values/neo4j-values.minikube.yaml"),
  ]

  depends_on = [

    kubernetes_namespace.neo4j_ns
  ]
}

