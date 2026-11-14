#!/usr/bin/env python3
"""Generate private runtime configuration once. Requires python3-yaml and OpenSSL."""
import argparse
import base64
import binascii
import json
import os
from pathlib import Path
import secrets
import subprocess
import yaml

parser = argparse.ArgumentParser()
parser.add_argument('--host', default='meet.ineshast.ir')
parser.add_argument('--turn', default='turn.ineshast.ir')
parser.add_argument('--ip', default='94.184.47.206')
parser.add_argument('--image-tag', required=True)
args = parser.parse_args()
root = Path(__file__).resolve().parent
source = root.parent
runtime = root / 'runtime'
os.umask(0o077)
runtime.mkdir(exist_ok=True)


def encode(data):
    checksum = binascii.crc_hqx(data, 0).to_bytes(2, 'little')
    return base64.b32encode(data + checksum).decode().rstrip('=')


def keypair(prefix, algorithm='ED25519'):
    private = subprocess.check_output(['openssl', 'genpkey', '-algorithm', algorithm, '-outform', 'DER'])
    public = subprocess.check_output(['openssl', 'pkey', '-inform', 'DER', '-pubout', '-outform', 'DER'], input=private)
    seed = encode(bytes([144 | (prefix >> 5), (prefix << 3) & 255]) + private[-32:])
    return seed, encode(bytes([prefix]) + public[-32:])


secret_path = runtime / 'secrets.json'
if secret_path.exists():
    credentials = json.loads(secret_path.read_text())
else:
    user_seed, user_public = keypair(160)
    issuer_seed, issuer_public = keypair(0)
    curve_seed, curve_public = keypair(184, 'X25519')
    credentials = {name: secrets.token_hex(32) for name in ['api_secret', 'livekit_secret', 'redis_password', 'db_password', 'db_root_password']}
    credentials.update(api_key='PNM' + secrets.token_hex(10), livekit_key='LK' + secrets.token_hex(10), user_seed=user_seed, user_public=user_public, issuer_seed=issuer_seed, issuer_public=issuer_public, curve_seed=curve_seed, curve_public=curve_public)
    secret_path.write_text(json.dumps(credentials, indent=2))
c = credentials


def write_yaml(name, data):
    (runtime / name).write_text(yaml.safe_dump(data, sort_keys=False))


server = yaml.safe_load((source / 'config_sample.yaml').read_text())
server['client'].update(debug=False, api_key=c['api_key'], secret=c['api_secret'], path='/app/client/dist', bbb_join_host='https://' + args.host)
server['client']['webhook_conf'].update(enable=True, enable_for_per_meeting=True)
server['log_settings'].update(log_level='info')
server['livekit_info'].update(host='http://host.docker.internal:7880', client_host='wss://' + args.host + '/livekit', api_key=c['livekit_key'], secret=c['livekit_secret'])
server['livekit_sip_info']['enabled'] = False
server['redis_info'].update(host='redis:6379', password=c['redis_password'])
server['database_info'].update(host='db', username='plugnmeet', password=c['db_password'])
server['nats_info'].update(nats_urls=['nats://nats:4222'], nats_ws_urls=['wss://' + args.host + '/ws'], nkey=c['user_seed'], auth_callout_issuer_private=c['issuer_seed'], auth_callout_xkey_private=c['curve_seed'])
server['recorder_info']['recording_files_path'] = '/recording_files'
write_yaml('server.yaml', server)
recorder = yaml.safe_load((source / 'components/recorder/config_sample.yaml').read_text())
recorder['recorder'].update(max_limit=2, width=1280, height=720, post_mp4_convert=False)
recorder['recorder']['copy_to_path']['main_path'] = '/recording_files'
recorder['plugNmeet_info'].update(host='http://server:8080', api_key=c['api_key'], api_secret=c['api_secret'], join_host='https://' + args.host + '/?access_token=')
recorder['nats_info']['nats_urls'] = ['nats://nats:4222']
write_yaml('recorder.yaml', recorder)
write_yaml('livekit.yaml', {
    'port': 7880,
    'rtc': {'tcp_port': 7881, 'port_range_start': 50000, 'port_range_end': 60000, 'use_external_ip': True, 'node_ip': args.ip},
    'redis': {'address': '127.0.0.1:16379', 'password': c['redis_password']},
    'keys': {c['livekit_key']: c['livekit_secret']},
    'webhook': {'api_key': c['livekit_key'], 'urls': ['http://127.0.0.1:18080/webhook']},
    'turn': {'enabled': True, 'udp_port': 443, 'tls_port': 5349, 'external_tls': True, 'domain': args.turn, 'proxy_protocol': True, 'proxy_protocol_trusted_cidrs': ['127.0.0.0/8']},
    'logging': {'level': 'warn'},
    'room': {'enable_remote_unmute': True},
})
nats = (source / 'nats_server_sample.conf').read_text().replace('store_dir: /data/jetstream', 'store_dir: /data/jetstream\n    max_mem: 96MB\n    max_file: 1GB')
for old, new in [('UD4JADBMDSQCA5D475KSMFV43TAINASJ3TIEFQ7LP43XENPMQYFTRFKP', c['user_public']), ('AD7EHM6WBIVERNVY6K7T3QEL62EB2BNMK4TTHXTUHTYSJGS3SZFZ4HX5', c['issuer_public']), ('XAVXJLXWZR7W24SAOPN6YATDOF2B6URA4GMBKYM7SBQIFO4O6OLQKSZB', c['curve_public'])]:
    nats = nats.replace(old, new)
(runtime / 'nats.conf').write_text(nats)
(runtime / 'redis.conf').write_text('appendonly yes\nmaxmemory 100mb\nmaxmemory-policy noeviction\nrequirepass ' + c['redis_password'] + '\n')
(runtime / 'observer.env').write_text('LIVEKIT_URL=ws://host.docker.internal:7880\nLIVEKIT_API_KEY=' + c['livekit_key'] + '\nLIVEKIT_API_SECRET=' + c['livekit_secret'] + '\nSPEAKING_OUTPUT=/data/speaking\nPLUGNMEET_URL=https://' + args.host + '\nSPEAKING_WEBHOOK_URL=https://api.ineshast.ir/api/v1/webhooks/plugnmeet\nPLUGNMEET_API_KEY=' + c['api_key'] + '\nPLUGNMEET_API_SECRET=' + c['api_secret'] + '\n')
# Redis drops privileges before reading configuration. Parent runtime stays 0700.
(runtime / 'redis.conf').chmod(0o644)
(runtime / 'db.env').write_text('MARIADB_DATABASE=plugnmeet\nMARIADB_USER=plugnmeet\nMARIADB_PASSWORD=' + c['db_password'] + '\nMARIADB_ROOT_PASSWORD=' + c['db_root_password'] + '\n')
client = (source / 'components/client/src/assets/config_sample.js').read_text()
# The source sample uses serverUrl as its public API location.
client = client.replace('http://localhost:8080', 'https://' + args.host).replace('https://demo.plugnmeet.com', 'https://' + args.host)
(runtime / 'client-config.js').write_text(client)
(root / '.env').write_text(''.join(name.upper() + '_IMAGE=ghcr.io/alirezamalekii/plugnmeet-' + name + ':' + args.image_tag + '\n' for name in ['server', 'recorder', 'livekit', 'redis', 'nats', 'mariadb']))
print('Runtime configuration generated. Secrets are in runtime/secrets.json with private permissions.')
