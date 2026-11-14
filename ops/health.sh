#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
docker compose -f compose.yml ps
while IFS= read -r service; do
    container="$(docker compose -f compose.yml ps -a -q "$service")"
    if [[ -z "$container" ]]; then
        echo "Missing service $service" >&2
        exit 1
    fi
    state="$(docker inspect "$container" --format '{{.State.Status}}')"
    if [[ "$state" != running ]]; then
        echo "Service $service is $state" >&2
        exit 1
    fi
    health="$(docker inspect "$container" --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}none{{end}}')"
    if [[ "$health" != none && "$health" != healthy ]]; then
        echo "Service $service health is $health" >&2
        exit 1
    fi
done < <(docker compose -f compose.yml config --services)
curl --fail --silent --show-error https://meet.ineshast.ir/healthCheck
printf '\n'
df -h /var/lib/docker
free -h
if (( $(df --output=pcent /var/lib/docker | tail -1 | tr -dc '0-9') > 80 )); then
    echo 'Docker filesystem usage exceeds 80 percent' >&2
    exit 1
fi
