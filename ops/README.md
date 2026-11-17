# Deployment operations

The host has one HAProxy listener on TCP80/443. TURN UDP443 and media TCP7881/UDP50000–60000 belong to LiveKit. NATS18222, server18080 and Redis16379 publish on loopback only. Add future websites using host routing in websites_http with independent Compose stacks.

## Deploy

Install Docker Compose, HAProxy, Certbot, python3-yaml, OpenSSL, age. Place this checkout at /opt/plugnmeet. Runtime config is generated only on the host and excluded from Git. configure.py generates API/LiveKit/DB/Redis/NATS credentials and reuses them on upgrades. Preserve runtime/secrets.json.

The source-pinned client is built into the server image. Server and recorder share recordings. The speaking observer has a dedicated volume. Pilot captures720p, allows two recordings, skips extra post-conversion to preserve live-media CPU. Recording acceptance remains required.

Run ineshast-images GitHub Actions at the exact source SHA. All mirrors and app images share the full-SHA tag. Once jobs pass, copy ops and source templates to the host and run `cd /opt/plugnmeet/ops; ./deploy.sh FULL_COMMIT_SHA; ./health.sh`.

Deployment preserves prior image environment and records image/container identity after health succeeds. A successful release also snapshots its private runtime configuration, Compose file and installed HAProxy configuration. Rollback restores that matched set before starting the recorded immutable images. `rollback.sh PREVIOUS_HEALTHY_SHA` restores a recorded healthy release. Check database schema compatibility before upgrades and take backup first.

## رابط فارسی جلسه

ساخت `Dockerfile.ineshast` تغییرات رابط را روی commit ثابت زیرمخزن کلاینت اعمال می‌کند. فایل‌های ایران‌یکان از نشانی دارایی‌های آی‌نشست دریافت و با SHA-256 بررسی می‌شوند. در ساخت محلی، از مسیر `components/client` ابتدا `git apply ../../docker-build/ineshast-client.patch` را اجرا کنید. سپس فایل‌های `Shared/fonts/iranyekan/iranyekan{regular,bold}.woff` را به ترتیب به `components/client/src/assets/ineshast/{regular,bold}.woff` کپی کنید. سپس از مسیر کلاینت `pnpm install --frozen-lockfile` و `pnpm build` را اجرا کنید.

در انتشار مستقیم نسخهٔ رابط، لایهٔ `client/dist` روی تصویر سرور نصب‌شده قرار می‌گیرد و فقط کانتینر وب از نو ساخته می‌شود. رسید در مسیر `/opt/plugnmeet/ops/releases/client-<SHA>` ثبت می‌شود. رسید میزبان شناسهٔ تصاویر قبل و بعد و commit رابط را نگه می‌دارد. نسخهٔ پشتیبان `.env` در همان پوشهٔ رسید قرار می‌گیرد. برای بازگشت رابط، در زمان خالی بودن اتاق‌ها آن را به `.env` بازگردانید و `docker compose up -d --no-deps --pull never server` و سپس `./health.sh` را اجرا کنید.

### آزمون‌های تکمیلی رابط

- یک جلسه با تعداد دوربین بیش از ظرفیت یک صفحه برگزار کنید. عبور از مرز صفحه، اسکرول، دوربین سنجاق‌شده و قطع و وصل دوربین را بررسی کنید تا هر تصویر در مسیر اسکرول و صفحه‌بندی قابل دسترسی بماند.
- در سافاری آیفون و یک مرورگر اندرویدی در دو جهت نمایش، پیش‌نمایش، مجوز دستگاه، جابه‌جایی پنل و فهرست دوربین‌ها را امتحان کنید. سیستم مجوزها و ابعاد نمایش در این مرورگرها جداگانه عمل می‌کنند.
- متن منوهای پیشرفتهٔ تختهٔ سفید، رأی‌گیری، اتاق گروهی، خطاهای دستگاه و صادر کردن گزارش را در نقش میزبان و مهمان بررسی کنید. این مسیرها متن‌هایی دارند که فقط هنگام اجرای همان قابلیت نمایش داده می‌شوند.
- یک جلسهٔ دارای اشتراک‌گذاری صفحه و ضبط را پایان دهید و خروجی را بازبینی کنید. تصویر کلاس زنده و چیدمان ضبط از مسیرهای جداگانه تولید می‌شوند.
- با حرکت کاهش‌یافته در تنظیمات سیستم، ناوبری با صفحه‌کلید و نمایشگر کوچک، تمرکز روی دکمه‌ها و خاموش بودن حرکت پس‌زمینه را جداگانه بررسی کنید.

## Certificate

First issuance uses standalone80. Renewal uses standalone9080 through HAProxy ACME route. Host /etc/letsencrypt/cli.ini sets http-01-port9080. Install renew-certificate.sh in /etc/letsencrypt/renewal-hooks/deploy/plugnmeet-haproxy. Verify `certbot renew --dry-run --http-01-port 9080`. User supplied no ACME contact email, so supported no-email registration was selected.

## Backup and recovery

Run backup.sh AGE_PUBLIC_RECIPIENT after meetings end and transfer encrypted artifact off-host. Backup includes SQL, private config, recordings, uploads, artifacts and speaking events. Encrypted backups are copied to the Mac at ~/Library/Application Support/plugNmeet/backups using fetch-backups.sh. The age private identity is ~/.ssh/plugnmeet-backup.age with mode0600, outside Git. Server timer runs daily02:30UTC and retains seven snapshots. The Mac LaunchAgent ir.ineshast.plugnmeet-backups runs the encrypted-only fetch hourly using the pinned SSH alias plugnmeet-production. Its plist is in ~/Library/LaunchAgents and private logs in ~/Library/Logs/plugNmeet. The local fetch command can also be rerun after deployments.

Recover into a clean separate project by decrypting with the corresponding age private identity, creating volumes, importing database.sql using root database credential, restoring each volume archive preserving paths, restoring runtime and image environment, then starting services. Verify real recording playback, speaking file and authenticated API before switching traffic. An initial empty-data rehearsal decrypted the archive and imported all four database tables in an isolated no-network MariaDB container. Recovery of real media still requires the pilot recording.

## Acceptance

Default-provider activation requires two-room/ten-user/two-recorder tests, forced TURN, audio/video/screenshare, speech interval alignment, webhook reconciliation, BBB compatibility and rollback. HTTP health proves API availability only.

## LiveKit version

Final deployment targets LiveKit v1.13.11, whose TURN configuration supports proxy_protocol and trusted proxy CIDRs. HAProxy sends PROXYv2 to preserve the real TURN client address. The initial API smoke used v1.9.11 with these newer options disabled. Upgrade and repeat forced-TURN/media acceptance before final use. Upstream source evidence is https://github.com/livekit/livekit/blob/v1.13.11/pkg/config/config.go#L325-L337.
