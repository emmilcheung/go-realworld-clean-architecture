resource "kubernetes_persistent_volume_claim" "postgresql_pv_claim" {
  metadata {
    labels = {
      app = "postgresql"
    }

    name      = "postgresql-pv-claim"
    namespace = "default"
  }

  spec {
    access_modes = ["ReadWriteMany"]

    resources {
      requests = {
        storage = "1Gi"
      }
    }

    storage_class_name = "manual"
    volume_name        = "postgresql-pv-volume"
  }
}
