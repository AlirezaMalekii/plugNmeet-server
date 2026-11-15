#!/usr/bin/env python3
"""Restore plugNmeet backends while preserving the shared host's other routes."""
import re
import sys
from pathlib import Path


BACKENDS = re.compile(
    r"^backend (turn|plugnmeet|livekit|nats)\n.*?(?=^(?:backend|frontend|global|defaults)\b|\Z)",
    re.MULTILINE | re.DOTALL,
)


def restore(snapshot: str, current: str) -> str:
    saved = {match.group(1): match.group(0) for match in BACKENDS.finditer(snapshot)}
    present = [match.group(1) for match in BACKENDS.finditer(current)]
    expected = {"turn", "plugnmeet", "livekit", "nats"}
    if set(saved) != expected or set(present) != expected or len(present) != 4:
        raise ValueError("Both configurations must contain the four unique plugNmeet backends")
    return BACKENDS.sub(lambda match: saved[match.group(1)], current)


if __name__ == "__main__":
    snapshot, current, output = map(Path, sys.argv[1:])
    output.write_text(restore(snapshot.read_text(), current.read_text()))
