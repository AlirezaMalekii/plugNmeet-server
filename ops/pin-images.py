#!/usr/bin/env python3
"""Replace release tags with pulled immutable repository digests."""
import json
from pathlib import Path
import subprocess

path = Path(__file__).resolve().parent / '.env'
lines = []
for line in path.read_text().splitlines():
    name, reference = line.split('=', 1)
    repository = reference.split('@')[0].rsplit(':', 1)[0]
    image = json.loads(subprocess.check_output(['docker', 'image', 'inspect', reference]))[0]
    digest = next(value for value in image['RepoDigests'] if value.startswith(repository + '@'))
    lines.append(name + '=' + digest)
path.write_text('\n'.join(lines) + '\n')
