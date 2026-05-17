output "backend_url" {
  description = "Backend API URL"
  value       = "http://localhost:${var.backend_port}/api/health"
}

output "frontend_url" {
  description = "Frontend application URL"
  value       = "http://localhost:${var.frontend_port}"
}

output "network_name" {
  description = "Docker network name"
  value       = docker_network.app_network.name
}

output "backend_container_id" {
  description = "Backend container ID"
  value       = docker_container.backend.id
}

output "frontend_container_id" {
  description = "Frontend container ID"
  value       = docker_container.frontend.id
}
