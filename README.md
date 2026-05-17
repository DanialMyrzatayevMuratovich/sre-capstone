# CinemaHub

A full-stack cinema booking platform with a complete SRE implementation — Infrastructure as Code, CI/CD pipeline, observability stack, and load testing.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25, Gin, MongoDB Atlas |
| Frontend | Vue.js 3, Vite, Nginx |
| Database | MongoDB Atlas (cloud) |
| Containerization | Docker, Docker Compose |
| IaC | Terraform (Docker provider) |
| CI/CD | GitHub Actions, Docker Hub |
| Monitoring | Prometheus, Grafana, Alertmanager |
| Load Testing | Locust |

---

## Services

| Service | URL | Credentials |
|---------|-----|-------------|
| Frontend | http://localhost:3000 | — |
| Prometheus | http://localhost:9090 | — |
| Grafana | http://localhost:3001 | admin / admin123 |
| Alertmanager | http://localhost:9093 | — |
| Node Exporter | http://localhost:9101 | — |

---

## Quick Start

### Prerequisites

- Docker Desktop
- Docker Compose v2+
- MongoDB Atlas connection string

### 1. Clone the repository

```bash
git clone https://github.com/myrzatayev/sre-capstone.git
cd sre-capstone/movie-app
```

### 2. Set environment variables

Create a `.env` file in `movie-app/`:

```env
MONGO_URI=mongodb+srv://<user>:<password>@cluster.mongodb.net/
MONGO_DATABASE=cinemahub
JWT_SECRET=your-secret-key
```

### 3. Start all services

```bash
docker compose up -d --build
```

### 4. Verify everything is running

```bash
docker compose ps
```

---

## Project Structure

```
movie-app/
├── backend/                    # Go API server
│   ├── cmd/main.go
│   ├── handlers/               # HTTP handlers
│   ├── middleware/             # Auth, CORS, Prometheus metrics
│   ├── models/                 # MongoDB models
│   ├── routes/routes.go
│   └── Dockerfile
│
├── frontend/                   # Vue.js 3 SPA
│   ├── src/
│   │   ├── views/
│   │   ├── components/
│   │   └── services/api.js
│   ├── nginx.conf
│   └── Dockerfile
│
├── monitoring/                 # Observability stack
│   ├── prometheus.yml
│   ├── alert_rules.yml
│   ├── autoscale.sh            # Prometheus-based auto-scaler
│   ├── alertmanager/
│   └── grafana/
│       ├── dashboards/         # CinemaHub SRE dashboard
│       └── provisioning/
│
├── terraform/                  # Infrastructure as Code
│   ├── main.tf
│   ├── variables.tf
│   └── providers.tf
│
├── locust/                     # Load testing
│   ├── locustfile.py
│   └── locust.conf
│
├── docs/
│   └── slo-definition.md       # SLI/SLO definitions
│
├── .github/
│   └── workflows/
│       └── ci-cd.yml           # GitHub Actions pipeline
│
└── docker-compose.yml
```

---

## SLOs

| SLO | Target | Window |
|-----|--------|--------|
| Availability | ≥ 99.5% | 30-day rolling |
| Latency p95 | ≤ 300 ms | 5-minute |
| Latency p99 | ≤ 500 ms | 5-minute |
| Error Rate | ≤ 1% | 5-minute |

---

## CI/CD Pipeline

Every push to `main` triggers a 4-stage GitHub Actions workflow:

```
test-backend → build-backend → build-frontend → deploy
```

Docker images are published to Docker Hub:
- `myrzatayev/cinemahub-backend:latest`
- `myrzatayev/cinemahub-frontend:latest`

---

## Load Testing

```bash
pip3 install locust
cd locust && locust -f locustfile.py
```

Open http://localhost:8089 — set 50 users, spawn rate 5/s.

---

## Scaling

```bash
# Scale backend to 2 replicas
docker compose up -d --scale backend=2 --no-recreate

# Auto-scale based on Prometheus metrics
bash monitoring/autoscale.sh
```

---

## License

MIT
