#!/usr/bin/env bash
set -euo pipefail
certificate=/etc/letsencrypt/live/meet.ineshast.ir
install -d -m 700 /etc/haproxy/ssl
umask 077
cat "$certificate/fullchain.pem" "$certificate/privkey.pem" > /etc/haproxy/ssl/meet.pem.new
mv /etc/haproxy/ssl/meet.pem.new /etc/haproxy/ssl/meet.pem
haproxy -c -f /etc/haproxy/haproxy.cfg
if systemctl is-active --quiet haproxy; then
    systemctl reload haproxy
else
    systemctl start haproxy
fi
