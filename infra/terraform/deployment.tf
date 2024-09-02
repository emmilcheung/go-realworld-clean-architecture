resource "kubernetes_deployment" "adminer" {
  metadata {
    labels = {
      app   = "adminer"
      group = "db"
    }
    name      = "adminer"
    namespace = "default"
  }

  spec {
    min_ready_seconds         = "0"
    paused                    = "false"
    progress_deadline_seconds = "600"
    replicas                  = "1"
    revision_history_limit    = "10"

    selector {
      match_labels = {
        app = "adminer"
      }
    }

    strategy {
      rolling_update {
        max_surge       = "25%"
        max_unavailable = "25%"
      }

      type = "RollingUpdate"
    }

    template {
      metadata {
        labels = {
          app   = "adminer"
          group = "db"
        }
      }

      spec {
        automount_service_account_token = "false"

        container {
          env {
            name  = "ADMINER_DEFAULT_SERVER"
            value = "postgres"
          }

          env {
            name  = "ADMINER_DESIGN"
            value = "pepa-linha"
          }

          image             = "adminer:4.7.6-standalone"
          image_pull_policy = "IfNotPresent"
          name              = "adminer"

          port {
            container_port = "8080"
            protocol       = "TCP"
          }

          resources {
            limits = {
              cpu    = "500m"
              memory = "256Mi"
            }
          }

          stdin                      = "false"
          stdin_once                 = "false"
          termination_message_path   = "/dev/termination-log"
          termination_message_policy = "File"
          tty                        = "false"
        }

        dns_policy                       = "ClusterFirst"
        enable_service_links             = "false"
        host_ipc                         = "false"
        host_network                     = "false"
        host_pid                         = "false"
        restart_policy                   = "Always"
        share_process_namespace          = "false"
        termination_grace_period_seconds = "30"
      }
    }
  }
}

resource "kubernetes_deployment" "api-server" {
  metadata {
    labels = {
      app = "api-deployment"
    }

    name      = "api-deployment"
    namespace = "default"
  }

  spec {
    min_ready_seconds         = "0"
    paused                    = "false"
    progress_deadline_seconds = "600"
    replicas                  = "3"
    revision_history_limit    = "10"

    selector {
      match_labels = {
        app = "api-deployment"
      }

    }
    strategy {
      type = "Recreate"
    }
    template {
      metadata {
        labels = {
          app = "api-deployment"
        }
      }

      spec {
        automount_service_account_token = "false"

        container {
          env {
            name  = "PORT"
            value = "8080"
          }

          image             = "realworld"
          image_pull_policy = "Never"
          name              = "api"

          port {
            container_port = "8080"
            protocol       = "TCP"
          }

          stdin                      = "false"
          stdin_once                 = "false"
          termination_message_path   = "/dev/termination-log"
          termination_message_policy = "File"
          tty                        = "false"
        }

        dns_policy                       = "ClusterFirst"
        enable_service_links             = "false"
        host_ipc                         = "false"
        host_network                     = "false"
        host_pid                         = "false"
        restart_policy                   = "Always"
        share_process_namespace          = "false"
        termination_grace_period_seconds = "30"
      }
    }
  }
}

resource "kubernetes_deployment" "jaeger_all_in_one" {
  metadata {
    labels = {
      app = "jaeger-deployment"
    }

    name      = "jaeger-deployment"
    namespace = "default"
  }

  spec {
    min_ready_seconds         = "0"
    paused                    = "false"
    progress_deadline_seconds = "600"
    replicas                  = "1"
    revision_history_limit    = "10"

    selector {
      match_labels = {
        app = "jaeger-deployment"
      }
    }

    strategy {
      rolling_update {
        max_surge       = "25%"
        max_unavailable = "25%"
      }

      type = "RollingUpdate"
    }

    template {
      metadata {
        labels = {
          app = "jaeger-deployment"
        }
      }

      spec {
        automount_service_account_token = "false"

        container {
          env {
            name  = "COLLECTOR_ZIPKIN_HTTP_PORT"
            value = "9411"
          }

          image             = "jaegertracing/all-in-one:1"
          image_pull_policy = "IfNotPresent"
          name              = "jaeger-container"

          port {
            container_port = "14250"
            host_port      = "14250"
            protocol       = "TCP"
          }

          port {
            container_port = "14268"
            host_port      = "14268"
            protocol       = "TCP"
          }

          port {
            container_port = "16686"
            host_port      = "16686"
            protocol       = "TCP"
          }

          port {
            container_port = "5775"
            host_port      = "5775"
            protocol       = "UDP"
          }

          port {
            container_port = "5778"
            host_port      = "5778"
            protocol       = "TCP"
          }

          port {
            container_port = "6831"
            host_port      = "6831"
            protocol       = "UDP"
          }

          port {
            container_port = "6832"
            host_port      = "6832"
            protocol       = "UDP"
          }

          port {
            container_port = "9411"
            host_port      = "9411"
            protocol       = "TCP"
          }

          stdin                      = "false"
          stdin_once                 = "false"
          termination_message_path   = "/dev/termination-log"
          termination_message_policy = "File"
          tty                        = "false"
        }

        dns_policy                       = "ClusterFirst"
        enable_service_links             = "false"
        host_ipc                         = "false"
        host_network                     = "false"
        host_pid                         = "false"
        restart_policy                   = "Always"
        share_process_namespace          = "false"
        termination_grace_period_seconds = "30"
      }
    }
  }
}

resource "kubernetes_deployment" "redis" {
  metadata {
    labels = {
      app = "redis-deployment"
    }

    name      = "redis-deployment"
    namespace = "default"
  }

  spec {
    min_ready_seconds         = "0"
    paused                    = "false"
    progress_deadline_seconds = "600"
    replicas                  = "1"
    revision_history_limit    = "10"

    selector {
      match_labels = {
        app = "redis-deployment"
      }
    }

    strategy {
      rolling_update {
        max_surge       = "25%"
        max_unavailable = "25%"
      }

      type = "RollingUpdate"
    }

    template {
      metadata {
        labels = {
          app = "redis-deployment"
        }
      }

      spec {
        automount_service_account_token = "false"

        container {
          image             = "redis:latest"
          image_pull_policy = "Always"
          name              = "api-redis"

          port {
            container_port = "6379"
            host_port      = "6379"
            protocol       = "TCP"
          }

          stdin                      = "false"
          stdin_once                 = "false"
          termination_message_path   = "/dev/termination-log"
          termination_message_policy = "File"
          tty                        = "false"
        }

        dns_policy                       = "ClusterFirst"
        enable_service_links             = "false"
        host_ipc                         = "false"
        host_network                     = "false"
        host_pid                         = "false"
        restart_policy                   = "Always"
        share_process_namespace          = "false"
        termination_grace_period_seconds = "30"
      }
    }
  }
}


