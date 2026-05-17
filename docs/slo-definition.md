# CinemaHub — SLI/SLO Definitions

## Overview

CinemaHub is a cinema booking platform. The following SLIs and SLOs are defined
for the backend API service (`cinemahub-backend`).

---

## Service Level Indicators (SLIs)

| SLI | Description | Prometheus Query |
|-----|-------------|-----------------|
| **Availability** | Fraction of requests that return a non-5xx response | `1 - (rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]))` |
| **Latency p95** | 95th percentile request duration | `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))` |
| **Latency p99** | 99th percentile request duration | `histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))` |
| **Error Rate** | Fraction of requests resulting in 5xx errors | `rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])` |
| **Throughput** | Requests per second handled by the API | `rate(http_requests_total[5m])` |
| **In-Flight** | Number of requests currently being processed | `http_requests_in_flight` |

---

## Service Level Objectives (SLOs)

| SLO | Target | Measurement Window | Alert Threshold |
|-----|--------|--------------------|-----------------|
| **Availability** | ≥ 99.5% | 30-day rolling | < 99.0% triggers PagerDuty |
| **Latency p95** | ≤ 300 ms | 5-minute window | > 400 ms triggers warning |
| **Latency p99** | ≤ 500 ms | 5-minute window | > 500 ms triggers critical |
| **Error Rate** | ≤ 1% | 5-minute window | > 1% triggers HighErrorRate alert |
| **In-Flight Requests** | ≤ 50 concurrent | Real-time | > 50 triggers auto-scale |

---

## Error Budget

Based on the **99.5% Availability SLO** over a 30-day month:

| Metric | Value |
|--------|-------|
| Total minutes/month | 43,200 |
| Allowed downtime (0.5%) | **216 minutes (~3.6 hours)** |
| Error budget burn rate (normal) | 1× |
| Alert if burn rate exceeds | 5× (budget exhausted in < 6 days) |

---

## Auto-Scaling Policy

| Condition | Action |
|-----------|--------|
| `http_requests_in_flight > 5` | Scale backend UP (+1 replica) |
| `http_requests_in_flight < 2` | Scale backend DOWN (-1 replica) |
| Maximum replicas | 3 |
| Minimum replicas | 1 |
| Check interval | 10 seconds |

---

## SLO Rationale

- **Availability 99.5%** — Cinema booking is a user-facing transactional service. Brief maintenance windows are acceptable, but extended outages directly cause revenue loss and user churn.
- **p95 ≤ 300ms** — Users browsing movies and selecting seats expect near-instant responses. Above 300ms the UI feels sluggish and cart abandonment increases.
- **p99 ≤ 500ms** — Even worst-case requests (complex queries, seat locking) must complete within half a second to prevent booking timeouts.
- **Error Rate ≤ 1%** — Payment and booking flows require high reliability. A 1% error rate means roughly 1 in 100 booking attempts fail, which is the maximum tolerable threshold.

---

## Alert Rules Summary

| Alert Name | Condition | Severity |
|-----------|-----------|----------|
| `HighErrorRate` | 5xx rate > 1% for 5 min | critical |
| `HighLatency` | p99 > 500ms for 5 min | warning |
| `BackendDown` | backend scrape target = 0 | critical |
| `HighInFlightRequests` | in-flight > 50 for 2 min | warning |
