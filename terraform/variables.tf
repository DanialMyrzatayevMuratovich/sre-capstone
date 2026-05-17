variable "app_name" {
  description = "Application name used as prefix for all resources"
  type        = string
  default     = "cinemahub"
}

variable "backend_port" {
  description = "Port exposed by the backend container"
  type        = number
  default     = 8080
}

variable "frontend_port" {
  description = "Port exposed by the frontend container"
  type        = number
  default     = 3000
}

variable "mongo_uri" {
  description = "MongoDB Atlas connection string"
  type        = string
  sensitive   = true
}

variable "jwt_secret" {
  description = "JWT secret key for authentication"
  type        = string
  sensitive   = true
}

variable "mongo_database" {
  description = "MongoDB database name"
  type        = string
  default     = "cinema_booking_db"
}
