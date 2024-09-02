resource "kubernetes_service" "adminer" {
  metadata {
    labels = {
      group = "db"
    }
    name      = "adminer"
    namespace = "default"
  }

  spec {
    ip_families            = ["IPv4"]
    ip_family_policy       = "SingleStack"

    port {
      port        = "8080"
      protocol    = "TCP"
      target_port = "8080"
    }

    publish_not_ready_addresses = "false"

    selector = {
      app = kubernetes_deployment.adminer.metadata.0.labels.app
    }

    session_affinity = "None"
    type             = "ClusterIP"
  }
}

resource "kubernetes_service" "api_server" {
  metadata {
    labels = {
      api = "api"
    }

    name      = "api"
    namespace = "default"
  }

  spec {
    ip_families            = ["IPv4"]
    ip_family_policy       = "SingleStack"

    port {
      name        = "8080"
      port        = "8080"
      protocol    = "TCP"
      target_port = "8080"
    }

    publish_not_ready_addresses = "false"

    selector = {
      app = kubernetes_deployment.api-server.metadata.0.labels.app
    }

    session_affinity = "None"
    type             = "ClusterIP"
  }
}

resource "kubernetes_service" "jaeger_all_in_one" {
  metadata {
    labels = {
      service = "jaeger"
    }

    name      = "jaeger"
    namespace = "default"
  }

  spec {
    ip_families            = ["IPv4"]
    ip_family_policy       = "SingleStack"

    port {
      name        = "14250"
      port        = "14250"
      protocol    = "TCP"
      target_port = "14250"
    }

    port {
      name        = "14268"
      port        = "14268"
      protocol    = "TCP"
      target_port = "14268"
    }

    port {
      name        = "16686"
      port        = "16686"
      protocol    = "TCP"
      target_port = "16686"
    }

    port {
      name        = "5775"
      port        = "5775"
      protocol    = "UDP"
      target_port = "5775"
    }

    port {
      name        = "5778"
      port        = "5778"
      protocol    = "TCP"
      target_port = "5778"
    }

    port {
      name        = "6831"
      port        = "6831"
      protocol    = "UDP"
      target_port = "6831"
    }

    port {
      name        = "6832"
      port        = "6832"
      protocol    = "UDP"
      target_port = "6832"
    }

    port {
      name        = "9411"
      port        = "9411"
      protocol    = "TCP"
      target_port = "9411"
    }

    publish_not_ready_addresses = "false"

    selector = {
      app = kubernetes_deployment.jaeger_all_in_one.metadata.0.labels.app
    }

    session_affinity = "None"
    type             = "ClusterIP"
  }
}

resource "kubernetes_service" "postgresql" {
  metadata {
    name      = "postgresql"
    namespace = "default"
  }

  spec {
    ip_families            = ["IPv4"]
    ip_family_policy       = "SingleStack"

    port {
      port        = "5432"
      protocol    = "TCP"
      target_port = "5432"
    }

    publish_not_ready_addresses = "false"

    selector = {
      app = kubernetes_stateful_set.postgresql-sts.spec.0.selector.0.match_labels.app
    }

    session_affinity = "None"
    type             = "ClusterIP"
  }
}

resource "kubernetes_service" "redis" {
  metadata {
    labels = {
      app = "redis"
    }

    name      = "redis"
    namespace = "default"
  }

  spec {
    ip_families            = ["IPv4"]
    ip_family_policy       = "SingleStack"

    port {
      name        = "6379"
      port        = "6379"
      protocol    = "TCP"
      target_port = "6379"
    }

    publish_not_ready_addresses = "false"

    selector = {
      app = kubernetes_deployment.redis.metadata.0.labels.app
    }

    session_affinity = "None"
    type             = "ClusterIP"
  }
}
