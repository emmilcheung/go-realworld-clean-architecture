resource "kubernetes_ingress_v1" "nginx_proxy" {
  wait_for_load_balancer = true
  metadata {
    name = "ingress-service"
    annotations = {
      "kubernetes.io/ingress.class"           = "nginx"
      "nginx.ingress.kubernetes.io/use-regex" = "true"
    }
  }
  spec {

    rule {
      http {
        path {
          path = "/api"
          path_type = "Prefix"
          backend {
            service {
              name = kubernetes_service.api_server.metadata.0.name
              port {
                number = 8080
              }
            }
          }
        }
      }
    }

    rule {
      host = "adminer.localhost"
      http {
        path {
          path = "/"
          path_type = "Prefix"
          backend {
            service {
              name = kubernetes_service.adminer.metadata.0.name
              port {
                number = 8080
              }
            }
          }
        }
      }
    }
    rule {
      host = "jaeger.localhost"
      http {
        path {
          path = "/"
          path_type = "Prefix"
          backend {
            service {
              name = kubernetes_service.jaeger_all_in_one.metadata.0.name
              port {
                number = 16686
              }
            }
          }
        }
      }
    }
  }
}

