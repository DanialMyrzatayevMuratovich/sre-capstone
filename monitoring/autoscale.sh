#!/bin/bash
# CinemaHub Auto-Scaler
# Reads http_requests_in_flight from Prometheus and scales the backend service.
# Run from the movie-app directory: bash monitoring/autoscale.sh

set -euo pipefail

PROMETHEUS_URL="http://localhost:9090"
SCALE_UP_THRESHOLD=5    # scale up  when in-flight > 5
SCALE_DOWN_THRESHOLD=2  # scale down when in-flight < 2
MAX_REPLICAS=3
MIN_REPLICAS=1
CHECK_INTERVAL=10       # seconds between polls

current_replicas=1

query_metric() {
    local raw
    raw=$(curl -sf "${PROMETHEUS_URL}/api/v1/query" \
        --data-urlencode "query=http_requests_in_flight" 2>/dev/null) || { echo 0; return; }
    echo "$raw" | python3 -c "
import sys, json
data = json.load(sys.stdin)
results = data.get('data', {}).get('result', [])
print(int(float(results[0]['value'][1])) if results else 0)
"
}

echo "======================================"
echo "  CinemaHub Auto-Scaler"
echo "======================================"
echo "  Prometheus : ${PROMETHEUS_URL}"
echo "  Scale UP   : in-flight > ${SCALE_UP_THRESHOLD}"
echo "  Scale DOWN : in-flight < ${SCALE_DOWN_THRESHOLD}"
echo "  Replicas   : ${MIN_REPLICAS} – ${MAX_REPLICAS}"
echo "  Interval   : ${CHECK_INTERVAL}s"
echo "======================================"
echo ""

while true; do
    in_flight=$(query_metric)
    ts=$(date '+%H:%M:%S')
    echo "[${ts}] in-flight=${in_flight}  replicas=${current_replicas}"

    if (( in_flight > SCALE_UP_THRESHOLD && current_replicas < MAX_REPLICAS )); then
        current_replicas=$(( current_replicas + 1 ))
        echo "  ▲ SCALE UP  → ${current_replicas} replicas  (in-flight=${in_flight} > ${SCALE_UP_THRESHOLD})"
        docker compose up -d --scale backend="${current_replicas}" --no-recreate 2>&1 | grep -E "Started|Running|Error" || true

    elif (( in_flight < SCALE_DOWN_THRESHOLD && current_replicas > MIN_REPLICAS )); then
        current_replicas=$(( current_replicas - 1 ))
        echo "  ▼ SCALE DOWN → ${current_replicas} replicas  (in-flight=${in_flight} < ${SCALE_DOWN_THRESHOLD})"
        docker compose up -d --scale backend="${current_replicas}" 2>&1 | grep -E "Started|Running|Error" || true
    fi

    sleep "${CHECK_INTERVAL}"
done
