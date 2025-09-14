# neo4j.tf

resource "helm_release" "neo4j" {
  name             = "neo4j-database"
  repository       = "https://helm.neo4j.com/neo4j"
  chart            = "neo4j"
  version          = "5.12.0"
  force_update     = true
  wait             = true

  values = [
    file("${path.module}/values/neo4j-values.yaml")
  ]
}
