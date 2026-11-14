# Deployment operations

The host has one HAProxy listener on TCP80/443. TURN UDP443 and media TCP7881/UDP50000–60000 belong to LiveKit. NATS18222, server18080 and Redis16379 publish on loopback only. Add future websites using host routing in websites_http with independent Compose stacks.

## Deploy

Install Docker Compose, HAProxy, Certbot, python3-yaml, OpenSSL, age. Place this checkout at /opt/plugnmeet. Runtime config is generated only on the host and excluded from Git. configure.py generates API/LiveKit/DB/Redis/NATS credentials and reuses them on upgrades. Preserve runtime/secrets.json.

The source-pinned client is built into the server image. Server and recorder share recordings. The speaking observer has a dedicated volume. Pilot captures720p, allows two recordings, skips extra post-conversion to preserve live-media CPU. Recording acceptance remains required.

Run ineshast-images GitHub Actions at the exact source SHA. All mirrors and app images share the full-SHA tag. Once jobs pass, copy ops and source templates to the host and run `cd /opt/plugnmeet/ops; ./deploy.sh FULL_COMMIT_SHA; ./health.sh`.

Deployment preserves prior image environment and records image/container identity after health succeeds. `rollback.sh PREVIOUS_HEALTHY_SHA` restores a recorded healthy release. Check database schema compatibility before upgrades and take backup first.

## Certificate

First issuance uses standalone80. Renewal uses standalone9080 through HAProxy ACME route. Host /etc/letsencrypt/cli.ini sets http-01-port9080. Install renew-certificate.sh in /etc/letsencrypt/renewal-hooks/deploy/plugnmeet-haproxy. Verify `certbot renew --dry-run --http-01-port 9080`. User supplied no ACME contact email, so supported no-email registration was selected.

## Backup and recovery

Run backup.sh AGE_PUBLIC_RECIPIENT after meetings end and transfer encrypted artifact off-host. Backup includes SQL, private config, recordings, uploads, artifacts and speaking events. Encrypted backups are copied to the Mac at ~/Library/Application Support/plugNmeet/backups using fetch-backups.sh. The age private identity is ~/.ssh/plugnmeet-backup.age with mode0600, outside Git. Server timer runs daily02:30UTC and retains seven snapshots. The local fetch command can be rerun after deployments.

Recover into a clean separate project by decrypting with the corresponding age private identity, creating volumes, importing database.sql using root database credential, restoring each volume archive preserving paths, restoring runtime and image environment, then starting services. Verify real recording playback, speaking file and authenticated API before switching traffic. An initial empty-data rehearsal decrypted the archive and imported all four database tables in an isolated no-network MariaDB container. Recovery of real media still requires the pilot recording.

## Acceptance

Default-provider activation requires two-room/ten-user/two-recorder tests, forced TURN, audio/video/screenshare, speech interval alignment, webhook reconciliation, BBB compatibility and rollback. HTTP health proves API availability only.
