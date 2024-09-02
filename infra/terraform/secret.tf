resource "kubernetes_secret" "postgresql-secret" {
  data = {
    POSTGRES_PASSWORD = "password"
  }

  immutable = "false"

  metadata {
    name      = "postgresql-secret"
    namespace = "default"
  }

  type = "Opaque"
}