
resource "kubernetes_config_map" "postgresql_config" {
  data = {
    POSTGRES_DB   = "realworld"
    POSTGRES_USER = "postgres"
  }

  metadata {
    labels = {
      app = "postgresql"
    }

    name      = "postgresql-config"
    namespace = "default"
  }
}
