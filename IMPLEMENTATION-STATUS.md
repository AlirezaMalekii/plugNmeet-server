# Implementation status

User-approved target: Docker deployment on 94.184.47.206, meet.ineshast.ir and turn.ineshast.ir, independent shared-server stack, BBB plus plugNmeet provider, recording/audio and speaking intervals, default only after acceptance.

## Evidence and baseline
- Parent iNeshast started clean at 03f03ef on main.
- Fork cloned, origin is AlirezaMalekii/plugNmeet-server, upstream is official.
- Work branch codex/ineshast-deployment. Local parent exclude protects nested repository.
- Dedicated Mac SSH public key installed. Independent key-only connection returned uid 0.
- Host Ubuntu 26.04.1, 2 vCPU, 3.8 GiB RAM, 55 GiB free, no swap.
- Password rotation and SSH password disable pending.

## Pending deliverables
- Complete Persian reference with coverage of all official pages.
- Pinned client and recorder submodules.
- Docker stack, TLS, backups, health and rollback.
- Laravel provider selection, authenticated durable webhook ingestion, reconciliation.
- Recorded audio and server-side speaking intervals.
- Two-room ten-user two-recorder pilot, TURN and BBB compatibility.
- PR, merge, local/main sync and image digest deployment evidence.

TDD mode off, decision skipped. Verification uses existing tests and targeted regression tests for provider, signatures, duplicate webhooks and interval alignment. Deployment acceptance requires real media evidence.
