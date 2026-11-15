#!/usr/bin/env python3
"""Authenticated native lifecycle smoke. Creates and ends only a unique test room."""
import hashlib
import hmac
import json
from pathlib import Path
import time
import urllib.request

root = Path(__file__).resolve().parent
credentials = json.loads((root / 'runtime/secrets.json').read_text())
room = 'ineshast-smoke-' + str(int(time.time()))


def call(endpoint, payload):
    body = json.dumps(payload, separators=(',', ':')).encode()
    signature = hmac.new(credentials['api_secret'].encode(), body, hashlib.sha256).hexdigest()
    request = urllib.request.Request('https://meet.ineshast.ir/auth' + endpoint, data=body, headers={'Content-Type': 'application/json', 'API-KEY': credentials['api_key'], 'HASH-SIGNATURE': signature})
    with urllib.request.urlopen(request, timeout=30) as response:
        result = json.load(response)
    print(endpoint, json.dumps({k: v for k, v in result.items() if k in ['status', 'status_code', 'msg', 'is_active']}))
    assert result.get('status'), 'API operation failed'
    return result


created = False
try:
    call('/room/create', {'room_id': room, 'metadata': {'room_title': 'iNeshast lifecycle smoke', 'room_features': {'allow_webcams': True, 'allow_screen_share': True, 'auto_gen_user_id': False, 'enable_analytics': True, 'recording_features': {'is_allow': True, 'is_allow_cloud': True, 'enable_auto_cloud_recording': False}}}})
    created = True
    assert call('/room/isRoomActive', {'room_id': room})['is_active']
    result = call('/room/getJoinToken', {'room_id': room, 'user_info': {'name': 'Smoke moderator', 'user_id': 'smoke-host', 'is_admin': True, 'is_hidden': False}})
    assert result.get('token'), 'Join token missing'
    print('Join token issued and redacted')
finally:
    if created:
        call('/room/endRoom', {'room_id': room})
