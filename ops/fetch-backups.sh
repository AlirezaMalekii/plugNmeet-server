#!/usr/bin/env bash
set -euo pipefail
umask 077
host="${PLUGNMEET_SSH_HOST:-plugnmeet-production}"
destination="${PLUGNMEET_BACKUP_DIR:-$HOME/Library/Application Support/plugNmeet/backups}"
mkdir -p "$destination"
rsync -av --include='*.age' --exclude='*' -e "ssh -o BatchMode=yes" "$host:/opt/plugnmeet/ops/backups/" "$destination/"
