# Implementation checkpoint, 2026-11-15

## Goal and account boundary
Shared-server Docker deployment, source synchronized with AlirezaMalekii/plugNmeet-server, complete documentation, plugNmeet alongside BBB in iNeshast, recordings/audio and speech intervals, then default activation after acceptance. User requests safe wrap at20percent account quota remaining and stop before10percent. Current slice reached22percent remaining. All browser pilots ended and both browsers closed. Only CI observation and safe wrap remain in this slice.

## Verified current state
- iNeshast branch codex/plugnmeet-provider, draft PR9, most recent pushed commit93b2178. Original clean baseline03f03ef. Core and Panel are not deployed and provider default remainsBBB.
- plugNmeet branch codex/ineshast-deployment, PR1. origin is user's fork, upstream official. Components are pinned submodules. All47 official pages including25API pages archived with checksums and Persian reference.
- SSH alias plugnmeet-production uses a dedicated Mac key and pinned host key. Root password was rotated into Mac Keychain; password SSH is disabled. Do not rotate again.
- Server94.184.47.206 now also runs AqayHerz shop/consultation/landing services. Preserve their shared HAProxy routes, data and networks. The plugNmeet rollback helper restores only owned backends and was validated against live shared configuration.
- Current plugNmeet release7461dcfb771cb9ceb5cb7ee81684cc5cd9067493. Seven services healthy, LiveKit1.13.11, image digests pinned. Certificates renew successfully. Native and actual PHP create/status/join/end smoke passed.
- Physical LiveKit RM_ SID differs from logical plugNmeet UUID. Observer resolves the logical ID through signed native API and verifies room generation.
- Missing speech transitions were caused by LiveKit scoping speaker updates to subscribed publishers. Observer now subscribes only to audio, drains RTP without decoding/storing media, and uses the internal pnm_agent prefix.
- Controlled15.007second synthetic WAV generated real RTC audio. Observer produced three start/end pairs. First transition975ms after source start, final transition1593ms after stop.
- Native recorded MP4 was decoded. Using the actual native recording-start second, first/last speech boundary errors were0.686458s and1.101375s. This proves one synthetic participant, not physical microphone or overlap accuracy.
- Forced TURN UDP passed with candidate type relay and actual188500inbound audio RTP bytes in750packets. Both participants were on the same Mac. TCP/TLS and independent networks remain unverified.
- Two rooms/two simultaneous recordings were previously verified with two synthetic participants. Ten-participant capacity remains unverified.
- Encrypted backup20261115T023617Z downloaded completely to Mac. A25456862byte MP4 was restored, SHA256 matched live source and restored video/audio decoded successfully. Daily remote backup and hourly Mac fetch timers are installed.
- Application rollback017->e39->017 passed previously. Current shared-proxy-safe rollback merger has two passing regression tests and live HAProxy validation preserving AqayHerz routes.

## Core implementation and current CI issue
- Core32tests/87assertions and Panel16tests passed in prior slice; Panel production build passed. Later Core33tests/90assertions and six webhook tests/21assertions passed.
- Native PHP400 was an explicitly empty optional URI logout_url. It is now omitted when empty and native PHP smoke passed.
- File creation_time is asset-ready time, not media start. Core uses signed start_recording createdAt correlated by room, logicalSID and recording_info.record_id. Missing start marker defers speech processing.
- Late recording discovery replays that session's journals. Inbox retries are ordered by updated_at then id, rotating failed attempts so101old invalid rows cannot starve a new valid event.
- PR9 full functional suite and MySQL checks passed. Twelve final PHPStan errors were corrected in93b2178 with four PHPDoc lines declaring relation generics and the speaking-state shape. Targeted PHPStan on ten related files, Pint, and fourteen regressions with37assertions passed. Updated full CI is pending. No baseline or ignore changes were made.
- Recording reconciliation currently scans a three-day updated_at window and fetches at most100artifacts. Very late or larger histories need manual sync or a later scoped improvement.

## Next steps
1. Reconcile this file with current Git status, PRs and quota. Private evidence is under ~/Library/Application Support/plugNmeet/pilot and private/recovery. Temporary files and browser sessions may disappear across interruption.
2. Observe the new full CI for93b2178 and require all PR9 checks before merge/deploy. The type fixes are already committed and pushed.
3. Finish source PR1 integration and synchronize local main with user's fork main after passing checks.
4. Deploy Core/Panel with matching private API credentials, database queue and inbox migration. Keep selection/default gated while testing signed webhook and recording ingestion end to end.
5. Complete ten-user/two-room/two-recorder capacity, independent-network, screen share, TURN TCP/TLS and overlapping-speaker tests. Activate plugNmeet default only after full acceptance.
6. During every quota wrap, end only task-owned pilot rooms and stop recording/browser processes. Keep the shared production stacks healthy.

Secrets belong only in private runtime files and Mac private storage. Never commit fixtures containing join tokens or raw secret values. See ops/acceptance.json for machine-readable evidence.
