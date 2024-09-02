resource "kubernetes_stateful_set" "postgresql-sts" {
  metadata {
    name      = "postgresql-sts"
    namespace = "default"
  }

  spec {
    pod_management_policy  = "OrderedReady"
    replicas               = "1"
    revision_history_limit = "10"

    selector {
      match_labels = {
        app = "postgresql-sts"
      }
    }

    service_name = "postgresql"

    template {
      metadata {
        labels = {
          app = "postgresql-sts"
        }
      }

      spec {
        automount_service_account_token = "false"

        container {
          env {
            name = "POSTGRES_PASSWORD"

            value_from {
              secret_key_ref {
                key      = "POSTGRES_PASSWORD"
                name     = "postgresql-secret"
                optional = "false"
              }
            }
          }

          env_from {
            config_map_ref {
              name     = kubernetes_config_map.postgresql_config.metadata.0.name
              optional = "false"
            }
          }

          image             = "postgres:13"
          image_pull_policy = "IfNotPresent"
          name              = "postgresql-db"

          port {
            container_port = "5432"
            host_port      = "5432"
            protocol       = "TCP"
          }

          stdin                      = "false"
          stdin_once                 = "false"
          termination_message_path   = "/dev/termination-log"
          termination_message_policy = "File"
          tty                        = "false"

          volume_mount {
            mount_path = "/var/lib/postgresql/data"
            name       = "postgresdb"
            read_only  = "false"
          }
        }

        dns_policy                       = "ClusterFirst"
        enable_service_links             = "false"
        host_ipc                         = "false"
        host_network                     = "false"
        host_pid                         = "false"
        restart_policy                   = "Always"
        share_process_namespace          = "false"
        termination_grace_period_seconds = "30"

        volume {
          name = "postgresdb"

          persistent_volume_claim {
            claim_name = kubernetes_persistent_volume_claim.postgresql_pv_claim.metadata.0.name
            read_only  = "false"
          }
        }
      }
    }
  }
}
