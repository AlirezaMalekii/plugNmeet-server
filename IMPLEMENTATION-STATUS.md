# plugNmeet implementation checkpoint

## Requested outcome and acceptance
Docker deployment on a shared server, source synchronized with the user's fork, complete Persian documentation, a second meeting provider in iNeshast, recording/video/audio, timestamped speech intervals, and plugNmeet as default after acceptance. Required pilot is two rooms, ten participants, two recordings, independent networks, forced TURN, recording playback, speech alignment within two seconds, webhook recovery, BBB compatibility and rollback.

## Verified baseline and completed work
- Parent iNeshast started clean at 03f03ef55ec0b3adac667ea8eff3988c4b4b9171 on main.
- Fork cloned locally. origin is AlirezaMalekii/plugNmeet-server, upstream is the official repository. Client and recorder are pinned submodules.
- Work branch codex/ineshast-deployment, draft PR #1. Parent local exclude protects the nested repository.
- All 47 official pages and 25 API pages archived with checksums at docs/upstream. PLUGNMEET-REFERENCE.fa.md includes Persian guidance plus original technical contracts.
- Dedicated Mac SSH key installed and independently tested. Password login disabled, root password rotated and stored in Mac Keychain. Pinned host-key SSH alias is plugnmeet-production.
- Host is Ubuntu26.04.1, two vCPU, 3.8GiB RAM, 55GiB initial free disk. A two-GiB swapfile was added.
- Docker/Compose, UFW, fail2ban, host HAProxy, certificates and renewal installed. meet.ineshast.ir/healthCheck returns Healthy.
- Docker Hub and mirror.gcr.io returned403 from the host. Fork GitHub Actions now build source images and mirror dependencies to GHCR, which the host can pull.
- Stable runtime image source 017f6706cd90a6afcbfa32491b2443b2b46faf51 was deployed with LiveKit1.13.11. Images are pinned by digest. Seven services include the speech observer.
- Native API create/status/join/end smoke passed. TLS renewal dry-run passed. Encrypted backup copied to Mac and database restored into an isolated test container. Actual media recovery remains unverified.

## Active work and important findings
- Core and Panel integration changes are currently uncommitted in the parent checkout. They must be reviewed and tested before integration or production deployment.
- Panel selector tests and production build passed in the prior slice. Core targeted tests passed, but a live PHP provider create request returned400 while the smaller native API smoke passed. The Core owner is diagnosing the exact request contract.
- Initial browser evidence covers one moderator, synthetic camera/microphone, Excellent connection and active recording. The old nine-user batch failed to enable devices. Ten-user acceptance has NOT passed.
- Old task-owned pilot rooms were left running across an account limit interruption. Runtime owner is ending only obsolete pilot rooms and preserving their recording evidence before new tests.
- A real observer correlation bug was found. LiveKit room SID is a physical RM_ ID; plugNmeet recording room_sid is a separate logical UUID. Fix017f6706cd90a6afcbfa32491b2443b2b46faf51 resolves the logical ID through authenticated getActiveRoomInfo and validates room generation. Its tests/build and deployment are pending.
- Core processes the observer's signed ineshast_speaking_events payload and must merge multiple segments under the logical session ID.
- Server resources are constrained. A single recorder used about oneGiB and over one core during the first live test. Recorder cap was raised to2300MiB with swap allowance. Two-recording performance still requires measurement.
- Default provider remains BBB until full acceptance. Core and Panel production rollout is pending.

## Resume order
1. Read this checkpoint, ops/acceptance.json, current Git status and latest CI results.
2. Check account remaining quota. User requests safe wrap-up at20percent remaining, with all pilot sessions stopped and checkpoint saved before10percent.
3. Finish the live PHP provider400 fix, run Core/Panel regressions and review compatibility.
4. Deploy the logical-session observer fix after successful CI, coordinate with browser tests.
5. Complete two-room media/recording/TURN/network/speech tests and encrypted media recovery/rollback.
6. Integrate through PRs, synchronize local main with fork main, deploy Core/Panel with matching private credentials and persistent database queue.
7. Activate plugNmeet as default only after the requested acceptance passes.

TDD mode off, decision skipped. Existing tests and targeted regression tests verify contracts. HTTP health and synthetic browser connection alone do not establish full acceptance. Secrets must stay in runtime config and Mac private storage.


## Safe-wrap checkpoint, 2026-11-14
- User requested wrap at20percent remaining and stop before10percent. Account readings jumped from31percent remaining to11percent during concurrent agent work. Agents were interrupted immediately and task-owned pilot rooms ended through the native API. No further feature work is authorized in this quota slice.
- Parent branch codex/plugnmeet-provider is committed and pushed at5922b79. Draft PR https://github.com/AlirezaMalekii/iNeshast/pull/9 . Core32tests/87assertions and Panel16tests passed; Panel production build passed. Native PHP create/status/join/end passed after omitting an empty logout_url.
- Recording timing now uses the signed start_recording marker correlated by room, logicalSID and recording_info.record_id. Native file creation_time is asset-ready time and must never be used as media start.
- Observer fix017f670 tests and image builds passed. It is deployed and journal logicalUUIDs were verified against the native API in both rooms.
- Rollback017->e39 and redeploy017 completed, with native API smoke after each. Runtime snapshots include matching Compose, configuration and HAProxy. Hourly encrypted backup fetch LaunchAgent is installed on Mac with pinned SSH alias.
- Two rooms each had one synthetic microphone/camera participant connected with Excellent quality and no console errors. Both cloud recorders were confirmed active by browser and two ffmpeg/Xvfb process pairs with growing MP4 files. This is two participants, not the requested ten.
- Recorder used about1.09GiB RAM and129percentCPU at the1.3core cap. Capacity with ten users remains unverified.
- Observer journals had correct logical session IDs but only observer_connected entries during synthetic audio. Speech transition generation, actual audio timing, forced TURN, independent networks, media restoration and production Core/Panel deployment remain pending.
- Core and Panel are NOT deployed. PNM selection remains disabled and default remainsBBB. Do not activate before acceptance.
- Durable private browser fixtures/logs and PHP smoke script are under ~/Library/Application Support/plugNmeet/pilot and private. /tmp and browser sessions were lost across interruption, so do not rely on temporary files when resuming.
- Resume by verifying all task-owned test rooms are ended, stack healthy, Git branches/PRs and CI. Then diagnose missing speech transitions with a known non-silent audio source before expanding to ten users or activating the provider.
