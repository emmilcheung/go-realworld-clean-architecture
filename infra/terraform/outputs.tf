# deployments
output "kubernetes_deployment_adminer_id" {
  value = "${kubernetes_deployment.adminer.id}"
}

output "kubernetes_deployment_api-server_id" {
  value = "${kubernetes_deployment.api-server.id}"
}

output "kubernetes_deployment_jaeger_all_in_one_id" {
  value = "${kubernetes_deployment.jaeger_all_in_one.id}"
}

output "kubernetes_deployment_redis_id" {
  value = "${kubernetes_deployment.redis.id}"
}

# services
output "kubernetes_service_adminer_id" {
  value = "${kubernetes_service.adminer.id}"
}

output "kubernetes_service_jaeger_all_in_one_id" {
  value = "${kubernetes_service.jaeger_all_in_one.id}"
}

output "kubernetes_service_postgresql_id" {
  value = "${kubernetes_service.postgresql.id}"
}

output "kubernetes_service_redis_id" {
  value = "${kubernetes_service.redis.id}"
}


# statefulsets
output "kubernetes_stateful_set_postgresql-sts_id" {
  value = "${kubernetes_stateful_set.postgresql-sts.id}"
}


# pv
output "kubernetes_persistent_volume_pv_volume_id" {
  value = "${kubernetes_persistent_volume.pv_volume.id}"
}

# pvc
output "kubernetes_persistent_volume_claim_postgresql_pv_claim_id" {
  value = "${kubernetes_persistent_volume_claim.postgresql_pv_claim.id}"
}


# config map
output "kubernetes_config_map_postgresql_config_id" {
  value = "${kubernetes_config_map.postgresql_config.id}"
}

# secrets
output "kubernetes_secret_postgresql-secret_id" {
  value = "${kubernetes_secret.postgresql-secret.id}"
}

# ingress - load balancer
# Display load balancer hostname (typically present in AWS)
output "load_balancer_hostname" {
  value = kubernetes_ingress_v1.nginx_proxy.status.0.load_balancer.0.ingress.0.hostname
}

# Display load balancer IP (typically present in GCP, or using Nginx ingress controller)
output "load_balancer_ip" {
  value = kubernetes_ingress_v1.nginx_proxy.status.0.load_balancer.0.ingress.0.ip
}