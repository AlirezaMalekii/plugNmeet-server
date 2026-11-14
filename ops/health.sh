#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
docker compose -f compose.yml ps
curl --fail --silent --show-error https://meet.ineshast.ir/healthCheck
printf '\n'
df -h /var/lib/docker
free -h
if (( $(df --output=pcent /var/lib/docker | tail -1 | tr -dc '0-9') > 80 )); then
    echo 'Docker filesystem usage exceeds 80 percent' >&2
    exit 1
fi
