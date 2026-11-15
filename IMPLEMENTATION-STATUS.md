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


## Continuation checkpoint, 2026-11-15 evening

- Both main branches were synchronized with their GitHub forks. plugNmeet PR1 merged at f76729efafce2648c0be677ebe2c5844214ab6ed. iNeshast PR9 merged at 5e8a9b0af49f4cb53ffb4274f068d394382e79e3. Full PR9 CI passed after keeping relation types scoped to plugNmeet code. Local full PHPStan passed, and27 targeted tests with77 assertions passed.
- A real two-network test passed using Mac and GitHub-hosted Chrome, with bidirectional synthetic audio and video. Evidence is in ~/Library/Application Support/plugNmeet/pilot/external-network-proof.json and GitHub run43730494740.
- Distributed ten-participant test passed admission and tracks simultaneously in two rooms, five participants each. Four independent GitHub-hosted runners each joined both rooms with audio/video sent and received. Run43730920828 succeeded. Both simultaneous recordings were created and ffprobe found audio and video, durations260.516 and240.300 seconds. Two Mac hosts and eight GitHub participants used synthetic media.
- Screen sharing used a simulated canvas through the client UI. Server reported SCREEN_SHARE, and a separate listener received live1280x720 video. Physical OS screen selection was not exercised.
- During the load sample, recorder CPU131.27percent with1.32GiB memory, LiveKit30.75percent and host load18.84 on two vCPU. Six samples of meet and AqayHerz public endpoints returned200. This sample does not establish capacity for heavier future co-located workloads.
- Test rooms and browsers closed. Temporary GitHub secrets and temporary workflow branches were deleted. Private proofs are under ~/Library/Application Support/plugNmeet/pilot, including distributed-proof.json, distributed-samples.jsonl, distributed-recordings-proof.json and distributed-health.log.
- iNeshast production deployment run43730183313 failed twice during private ghcr.io/alirezamalekii/ineshast-core:5e8a9b0af49f4cb53ffb4274f068d394382e79e3 image layer download with connection reset to registry CDN. Other three application images and dependency images were downloaded. The deployed release remains03f03ef and API/panel returned200.
- The deployment wrote new images.env and .env before failure. Both were restored atomically from /opt/ineshast/shared/backups/before-5e8a9b0af49f4cb53ffb4274f068d394382e79e3 after confirming current still points to03f03ef. Existing containers were retained. GitHub production CORE_ENV secret was updated with current credentials, PLUGNMEET_ENABLED=false and MEETING_DEFAULT_PROVIDER=bbb.
- Next step is complete the authorized deployment without changing current running workloads prematurely. Registry image is private and Mac gh token lacks read:packages. A task-scoped GitHub Actions job in the private iNeshast repo can pull the exact core image with GITHUB_TOKEN and upload a compressed docker-save artifact. Download the artifact to Mac and stream it via the existing SSH alias iran-vps into docker load. Then rerun failed deploy43730183313 and monitor to completion, verifying health, migration, exact running SHA, environment flags and actual signed webhook delivery. Do not make private images public or grant non-Mac SSH access to the plugNmeet host.
- Core and Panel production integration remains unverified and disabled until successful deploy. Default remains BBB. After deploy, verify native signed recording-start, recording-ready, audio extraction, journal ingestion, late recording replay and legacy BBB. Exercise actual overlapping-speaker timing. TURN UDP passed, TCP/TLS remains untested.
- User requested start safe wrap at20percent remaining, before10percent. At this checkpoint remaining was21percent. Any in-progress GitHub deployment must be checked explicitly before another attempt; no rerun was started after failure.
