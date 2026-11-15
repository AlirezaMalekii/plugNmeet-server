#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
if [[ $# != 1 ]]; then echo 'Usage: backup.sh <age public recipient>' >&2; exit 2; fi
command -v age >/dev/null
umask 077
mkdir -p backups
stamp="$(date -u +%Y%m%dT%H%M%SZ)"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT
docker compose -f compose.yml exec -T db sh -c 'exec mariadb-dump -uroot -p"$MARIADB_ROOT_PASSWORD" --single-transaction --routines --events plugnmeet' > "$work/database.sql"
cp -a runtime .env "$work/"
# Uploaded files can change during active meetings. Run after sessions finish for a consistent media snapshot.
for volume in recordings upload artifacts speaking; do
    mount="$(docker volume inspect "plugnmeet_$volume" --format '{{.Mountpoint}}')"
    tar -C "$mount" -czf "$work/$volume.tar.gz" .
done
tar -C "$work" -czf - . | age -r "$1" > "backups/$stamp.tar.gz.age"
printf '%s\n' "backups/$stamp.tar.gz.age"

python3 - <<'PYRETENTION'
from pathlib import Path
for path in sorted(Path('backups').glob('*.tar.gz.age'))[:-7]:
    path.unlink()
PYRETENTION
