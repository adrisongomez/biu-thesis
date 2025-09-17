# main.tf

resource "kubernetes_namespace" "neo4j_ns" {
  metadata {
    name = "neo4j"
  }
}

# Helm release: Neo4j
resource "helm_release" "neo4j" {
  name       = "neo4j-database"
  repository = "https://helm.neo4j.com/neo4j"
  chart      = "neo4j-standalone"

  namespace = kubernetes_namespace.neo4j_ns.metadata[0].name

  wait = false
  wait_for_jobs = false
  values = [
    file("${path.module}/values/neo4j-values.yaml")
  ]

  depends_on = [kubernetes_namespace.neo4j_ns]
}

# --- Outputs ---
output "neo4j_namespace" {
  value = kubernetes_namespace.neo4j_ns.metadata[0].name
}

output "neo4j_release_name" {
  value = helm_release.neo4j.name
}

output "neo4j_connection_info" {
  value = "Run 'kubectl port-forward svc/neo4j-database-ui -n ${kubernetes_namespace.neo4j_ns.metadata[0].name} 7474:7474' to connect."
}
