#!/usr/bin/env bash
set -euo pipefail
umask 077
key="${PLUGNMEET_SSH_KEY:-$HOME/.ssh/plugnmeet_ed25519}"
destination="${PLUGNMEET_BACKUP_DIR:-$HOME/Library/Application Support/plugNmeet/backups}"
mkdir -p "$destination"
rsync -av --include='*.age' --exclude='*' -e "ssh -i $key -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new" root@94.184.47.206:/opt/plugnmeet/ops/backups/ "$destination/"
