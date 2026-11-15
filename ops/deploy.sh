#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
if [[ $# != 1 || ! "$1" =~ ^[a-f0-9]{40}$ ]]; then
    echo 'Usage: deploy.sh <full source commit SHA>' >&2
    exit 2
fi
umask 077
mkdir -p releases
release="$1"
if [[ -f .env ]]; then cp .env "releases/before-$release.env"; fi
python3 configure.py --image-tag "$release"
docker compose -f compose.yml config --quiet
docker compose -f compose.yml pull
python3 pin-images.py
cp .env "releases/$release.env"
docker compose -f compose.yml create
subnet="$(docker network inspect plugnmeet_default --format '{{(index .IPAM.Config 0).Subnet}}')"
ufw allow from "$subnet" to any port 7880 proto tcp
docker compose -f compose.yml up -d --wait --wait-timeout 180
for attempt in {1..30}; do
    if curl --fail --silent http://127.0.0.1:18080/healthCheck >/dev/null; then
        docker compose -f compose.yml images --format json > "releases/$release.images.json"
        docker inspect $(docker compose -f compose.yml ps -q) --format '{{.Name}} {{.Image}} {{json .Config.Labels}}' > "releases/$release.containers.txt"
        mkdir -p "releases/$release.runtime"
        cp -a runtime/. "releases/$release.runtime/"
        cp compose.yml "releases/$release.compose.yml"
        cp /etc/haproxy/haproxy.cfg "releases/$release.haproxy.cfg"
        printf '%s\n' "$release" > releases/current
        exit 0
    fi
    sleep 2
done
echo 'Health check failed. Inspect logs and roll back with releases/before-<SHA>.env.' >&2
exit 1
