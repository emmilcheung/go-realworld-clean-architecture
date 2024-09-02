resource "kubernetes_persistent_volume" "pv_volume" {
  metadata {
    labels = {
      app  = "postgresql"
      type = "local"
    }

    name = "postgresql-pv-volume"
  }

  spec {
    access_modes = ["ReadWriteMany"]

    capacity = {
      storage = "1Gi"
    }

    claim_ref {
      name      = kubernetes_persistent_volume_claim.postgresql_pv_claim.metadata.0.name
      namespace = "default"
    }

    persistent_volume_reclaim_policy = "Delete"

    persistent_volume_source {
      host_path {
        path = "/mnt/data"
      }
    }

    storage_class_name = "manual"
    volume_mode        = "Filesystem"
  }
}
