# Testing Strategy tanpa Client Fisik

Plan untuk cover skenario hotspot (login event, expiry, transaction, SSE fan-out) tanpa butuh device client fisik untuk login — kombinasi mock TCP RouterOS API, in-memory DB, dan live integration ke router fisik (read+write user dummy).

## Tujuan & Coverage Target

- Cover **fitur hotspot yang sebelumnya butuh client real** (active stream, login event flow, expiry workflow, transaction recording) lewat mock.
- Cover **CRUD + voucher generate + expiry e2e** lewat live integration ke router v6+v7 (allowed modify state).
- Tetap **race-detector clean** (`go test -race ./...`).
- Bukan ditujukan untuk replace integration test yang sudah ada — **extend**, bukan replace.

## Audit Existing Infra

| Komponen | Status | Limitation |
|---|---|---|
| `internal/tcpmock` | Wire-protocol encoding OK | Single-conn, FIFO scripted reply, **no login handshake**, **no streaming emit**, no `.tag` echo |
| `test/integration/` | Connect ke router via `testutil.NewClient` | Tidak ada test expiry workflow / transaction record |
| `internal/testutil` | `RequireIntegration`, `UniqueName`, `NewClient` | Lengkap untuk router-level. Butuh helper baru untuk DB in-memory + devmgr stub |
| `service/expiry` | `executeExpiry` siap di-test (interface jelas) | Belum ada test sama sekali |

## Strategi 4 Pilar

### Pilar 1 — Extend `internal/tcpmock` untuk streaming + login (~250 LOC)

**Tambah ke `server.go`:**

- `Server.AcceptLogin()` — emulate RouterOS modern login (post-6.43): first sentence `/login` direply `!done` tanpa challenge. Opsional support legacy challenge kalau perlu.
- `Server.OnSentence(matcher func([]string) bool, replies ...[]string)` — registered handler yang match sentence + reply. Mengganti scripted FIFO yang fragile.
- `Server.StreamReply(matcher, taggedSentences ...[]string)` — saat sentence cocok dengan listen command, server echo `.tag` ke setiap reply dan kirim `!re` berkali-kali sampai `Cancel()` atau test `Close()`. Mode "follow" yang persistent.
- `Server.EmitToStream(streamID string, words ...[]string)` — push sentence baru ke stream yang aktif (untuk simulate "user login event" beberapa ms setelah subscribe).
- `Server.AssertReceived(t, matcher)` — testify-style assertion untuk command yang masuk.

**File baru `internal/tcpmock/router.go` (~150 LOC):**
Sentence builder helpers untuk test readability:

```go
package tcpmock

// Reply sentences ala RouterOS for tests
func DoneReply(extras ...string) []string                  // !done [+extras]
func ReReply(tag string, kvs ...string) []string           // !re =.tag=X =k=v ...
func TrapReply(tag, msg string) []string                   // !trap =.tag=X =message=...
func ActiveLoggedIn(tag, user, addr, mac string) []string  // typed shortcut
func ActiveLoggedOut(tag, id string) []string              // .dead=true
func UserPrint(tag, id, name, profile, comment string) []string
```

**File baru `internal/tcpmock/login.go` (~50 LOC):**
Modern login handler — match `/login` sentence dengan `=name=` `=password=` → reply `!done`.

### Pilar 2 — Test Infrastructure Helpers (~150 LOC)

**File baru `internal/testutil/mock.go`:**

- `NewMockDevice(t)` — start tcpmock, build `*roslib.Device` yang dial ke mock, return `(dev, mockSrv)`. Setup login auto. `t.Cleanup` shutdown semua.
- `NewSQLiteStore(t)` — sqlite in-memory + auto-migrate semua model (device, transaction, profile_config). Return `(deviceStore, txStore, profileStore)`.
- `NewTestDevMgr(t, dev)` — devmgr.Manager stub yang return ClientSet untuk satu device dummy.
- `Eventually(t, condition, timeout)` — poll wrapper untuk assert eventual state.

### Pilar 3 — Behavioral Tests pakai Mock (~600 LOC)

**File baru `api/sse/broker_streaming_test.go` (~150 LOC):**

- `TestBroker_loginEvent_emitsSSE` — mock emit active row → broker → 1 subscriber dapat `change` event dengan payload user yang benar.
- `TestBroker_logoutEvent_deadFlag` — mock emit `.dead=true` → SSE event `dead: true`.
- `TestBroker_3clients_singleStream` — 3 subscribers, 1 backend stream, semua dapat event sama (assert via channel count).
- `TestBroker_slowClient_othersStillReceive` — 1 subscriber tidak baca, 2 lain baca normal, slow drop counter incremented.

**File baru `service/expiry/service_test.go` (~250 LOC):**

- `TestService_expiredUser_remMode_callsDeleteUser` — sqlite + mock router. Mock reply user-print dengan comment past time. Run 1 tick. Assert: mock receive `/ip/hotspot/user/remove =.id=*1`.
- `TestService_expiredUser_ntfcMode_setsLimitAndKicks` — assert mock receive `user/set =limit-uptime=1s` + `active/remove`.
- `TestService_expiredUser_remcMode_recordsTransaction` — sqlite + mock. Tick. Query DB → 1 row dengan SaleMonth, Username, Price benar.
- `TestService_expiredUser_ntfMode_noTransactionRecorded` — non-`c` mode tidak insert tx.
- `TestService_validComment_notExpired_skip` — comment future date → no action.
- `TestService_invalidComment_skip` — comment bukan format mikhmon → no action.
- `TestService_backoff_onDeviceDisconnect` — devmgr return ErrDeviceNotConnected → state transisi ke backoff → log sekali (count log lines).
- `TestService_recover_onReconnect` — backoff → next tick devmgr return ClientSet → state kembali normal + log "recovered".

**File baru `mikrotik/hotspot/active_stream_test.go` (~100 LOC):**

Pure test untuk sentence handler tanpa router layer:

- `TestActiveStream_registerSendsFollowCommand` — register stream → tcpmock receive `/ip/hotspot/active/print follow =.tag=...`.
- `TestActiveStream_followOnly_sendsFollowOnlyFlag`.
- `TestActiveStream_unregisterSendsCancel` — `dev.UnregisterStream(id)` → mock receive `/cancel =tag=...`.

**File baru `api/handler/hotspot_voucher_test.go` (~100 LOC):**

End-to-end di gin handler level pakai mock router:

- `TestVoucher_generate_50_callsUserAdd50times` — POST request → mock receive 50× `user/add` dengan name unik, semua reply OK → response 200 dengan 50 vouchers.
- `TestVoucher_partialFailure_returns207` — mock reply `!trap` di voucher ke-10 → response 207 Multi-Status dengan 9 voucher OK + error field.
- `TestVoucher_invalidSpec_returns400` — batch_size=0 → 400 ValidationErr.

### Pilar 4 — Live Integration Tests (router fisik, no client) (~250 LOC)

User confirmed OK modify router state. Prefix `mikhmon-it-` + `t.Cleanup`.

**File baru `test/integration/expiry_workflow_test.go`:**

- `TestIntegration_Expiry_remMode_userDeleted` — add user `mikhmon-it-exp-<rand>` dengan comment past time + profile mode=rem. Run expiry checker 1 tick. Assert user `UserByName` return ErrNotFound (already deleted).
- `TestIntegration_Expiry_ntfMode_limitUptimeSet` — add user mode=ntf. Run tick. UserByName → assert `LimitUptime=1s`.
- `TestIntegration_Expiry_remcMode_recordsTransaction` — add user mode=remc + setup profile_config dengan Price. Run tick. Query tx store → 1 row exist.
- Cleanup: t.Cleanup hapus user yg masih ada (kalau test fail).

**File baru `test/integration/voucher_workflow_test.go`:**

- `TestIntegration_VoucherGenerate_10_thenList_thenBulkDelete` — full lifecycle 10 voucher di router fisik via workflows.GenerateVouchers → assert hot.UserList contains 10 → bulk delete → list assert empty.

**Extend `test/integration/hotspot_active_stream_test.go`:**

- `TestIntegration_ActiveStream_addUserDoesNotEmit` — verify add user (without login) TIDAK emit ActiveStream event (event = login event, bukan user creation). Important untuk konfirmasi pemahaman semantik RouterOS.

## File Changes Summary

**New files:**

- `internal/tcpmock/router.go` — sentence builder helpers
- `internal/tcpmock/login.go` — login handler
- `internal/tcpmock/streaming.go` — streaming emit + tag echo
- `internal/tcpmock/streaming_test.go` — self-test
- `internal/testutil/mock.go` — mock factory
- `internal/testutil/sqlite.go` — in-memory DB
- `api/sse/broker_streaming_test.go`
- `service/expiry/service_test.go`
- `mikrotik/hotspot/active_stream_test.go`
- `api/handler/hotspot_voucher_test.go`
- `test/integration/expiry_workflow_test.go`
- `test/integration/voucher_workflow_test.go`

**Modified:**

- `internal/tcpmock/server.go` — add `OnSentence`, `StreamReply`, `EmitToStream`, `AssertReceived` methods
- `internal/tcpmock/doc.go` — update scope description
- `test/integration/hotspot_active_stream_test.go` — extend with semantic test
- `docs/TESTING.md` — dokumentasi 4 layer strategy

## Estimasi LOC + Risk

| Pilar | LOC | Risk |
|---|---|---|
| 1. tcpmock extension | 250 | **Medium** — RouterOS login handshake bisa ada quirk di v6 vs v7 (mitigasi: support modern flow saja, dokumentasikan limitation) |
| 2. Test helpers | 150 | Low |
| 3. Behavioral tests | 600 | Low (deterministic mock) |
| 4. Live integration | 250 | Medium — butuh router online, user manual run `go test -tags=integration` |
| **Total** | **~1250** | |

## Order of Execution

Implementasi sekuensial (dependent):

1. **Phase 1** — tcpmock extension (login + streaming + helpers) + self-tests. Verifikasi: `go test ./internal/tcpmock/...` PASS dengan race.
2. **Phase 2** — testutil mock factories. Verifikasi: contoh trivial test pakai mock device PASS.
3. **Phase 3** — behavioral tests (SSE broker → expiry service → handler). Verifikasi: full test suite PASS, coverage report.
4. **Phase 4** — live integration (perlu env router). Verifikasi manual oleh user di v6 dan v7.
5. **Phase 5** — dokumentasi `docs/TESTING.md` dengan diagram alur + cara run per layer.

## Yang TIDAK di-cover (limitasi)

Sengaja keluar scope karena fundamentally butuh client:

- **Captive portal HTTP flow** — form submit, redirect, cookie issuance.
- **RADIUS Access-Request** — butuh RADIUS server + client.
- **Walled-garden DNS** — butuh DNS client.
- **Bandwidth throttling enforcement** — butuh real traffic.

Untuk fitur ini: smoke test manual oleh operator pre-release (dokumentasikan di `docs/TESTING.md`).

## Out of Scope (di plan lain)

- Auth middleware + JWT (Phase 5 terpisah).
- Password device encryption (Phase 5 terpisah).
- Replay fixture dari capture router (opsi yang tidak dipilih user di question sebelumnya).

## Verifikasi Final

- `go test -race ./...` PASS di mikhmon repo + roslib repo
- `go test -tags=integration -race ./test/integration/...` PASS (user yang run, butuh env `ROSLIB_ROUTER_ADDRESS`)
- Coverage report: target 70%+ di `service/expiry`, `api/sse`, `internal/tcpmock`
- `docs/TESTING.md` updated dengan 4 layer strategy + cara run
