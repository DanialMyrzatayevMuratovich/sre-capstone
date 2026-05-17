# Docker Network — изолированная сеть для сервисов
resource "docker_network" "app_network" {
  name = "${var.app_name}_network"
}

# Backend Image
resource "docker_image" "backend" {
  name = "${var.app_name}_backend:latest"
  build {
    context    = "${path.module}/../backend"
    dockerfile = "Dockerfile"
  }
  triggers = {
    dir_sha = sha1(join("", [for f in fileset("${path.module}/../backend", "**/*.go") : filesha1("${path.module}/../backend/${f}")]))
  }
}

# Frontend Image
resource "docker_image" "frontend" {
  name = "${var.app_name}_frontend:latest"
  build {
    context    = "${path.module}/../frontend"
    dockerfile = "Dockerfile"
  }
  triggers = {
    dir_sha = sha1(join("", [for f in fileset("${path.module}/../frontend/src", "**/*") : filesha1("${path.module}/../frontend/src/${f}")]))
  }
}

# Backend Container
resource "docker_container" "backend" {
  name  = "${var.app_name}_backend"
  image = docker_image.backend.image_id

  ports {
    internal = 8080
    external = var.backend_port
  }

  networks_advanced {
    name = docker_network.app_network.name
  }

  env = [
    "PORT=8080",
    "ENV=production",
    "MONGO_URI=${var.mongo_uri}",
    "MONGO_DATABASE=${var.mongo_database}",
    "JWT_SECRET=${var.jwt_secret}",
    "JWT_EXPIRATION=24h"
  ]

  restart = "unless-stopped"

  healthcheck {
    test         = ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/api/health"]
    interval     = "30s"
    timeout      = "10s"
    retries      = 3
    start_period = "10s"
  }
}

# Frontend Container
resource "docker_container" "frontend" {
  name  = "${var.app_name}_frontend"
  image = docker_image.frontend.image_id

  ports {
    internal = 80
    external = var.frontend_port
  }

  networks_advanced {
    name = docker_network.app_network.name
  }

  restart = "unless-stopped"

  depends_on = [docker_container.backend]
}
