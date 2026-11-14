#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
if [[ $# != 1 || ! "$1" =~ ^[a-f0-9]{40}$ ]]; then
    echo 'Usage: rollback.sh <previous successfully deployed full SHA>' >&2
    exit 2
fi
previous="$1"
if [[ ! -f "releases/$previous.images.json" ]]; then
    echo 'Requested SHA has no recorded healthy deployment receipt' >&2
    exit 1
fi
umask 077
if [[ ! -d "releases/$previous.runtime" ]]; then
    echo 'Requested release has no matching runtime configuration snapshot' >&2
    exit 1
fi
cp "releases/$previous.env" .env
cp -a "releases/$previous.runtime/." runtime/
cp "releases/$previous.compose.yml" compose.yml
cp "releases/$previous.haproxy.cfg" /etc/haproxy/haproxy.cfg
haproxy -c -f /etc/haproxy/haproxy.cfg
systemctl reload haproxy
docker compose -f compose.yml config --quiet
docker compose -f compose.yml up -d --wait --wait-timeout 180
curl --fail --silent --show-error https://meet.ineshast.ir/healthCheck
printf '%s\n' "$previous" > releases/current
