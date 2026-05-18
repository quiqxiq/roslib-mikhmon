# Testing Guide

Empat lapis test, dipisah lewat build tag supaya tiap CI run bisa memilih
scope yang tepat.

```
┌──────────────────────────────────────────────────────────────────────────┐
│  Layer 1: Pure unit                  (no tag)                            │
│    └─ domain/, scripts/, internal/* (rosfmt, tcpmock self-test, dst)     │
├──────────────────────────────────────────────────────────────────────────┤
│  Layer 2: Behavioral via mock router (no tag)                            │
│    └─ tcpmock + roslib.Device       (mikrotik/hotspot, api/sse,          │
│                                      api/handler/hotspot_voucher)        │
├──────────────────────────────────────────────────────────────────────────┤
│  Layer 3: Behavioral via mock + DB   (build tag: dbtest)                 │
│    └─ tcpmock + postgres testcontainer  (service/expiry behavioral)      │
├──────────────────────────────────────────────────────────────────────────┤
│  Layer 4: Live e2e ke router fisik   (build tag: integration            │
│                                       [+ dbtest untuk expiry workflow]) │
│    └─ ROSLIB_ROUTER_ADDRESS env       (test/integration/*)              │
└──────────────────────────────────────────────────────────────────────────┘
```

## Cara Run per Layer

### Layer 1+2: tanpa tag (default CI)

Cepat (<5s), tidak butuh Docker, tidak butuh router.

```bash
rtk go test -p 1 ./...
rtk go test -p 1 -race ./...
```

`-p 1` (sequential packages) **dianjurkan** karena beberapa test memakai
`tcpmock + *roslib.Device`. Saat dijalankan paralel antar paket, kontensi
resource OS dapat menyebabkan flake (test hang). Dalam isolasi tiap paket
PASS; gunakan `-p 1` untuk full-suite run.

### Layer 3: dbtest (postgres testcontainer)

Butuh Docker (atau Podman + `TESTCONTAINERS_RYUK_DISABLED=true`). Container
postgres start ~3-5 detik per test yang memanggil `testutil.NewStores`.

```bash
TESTCONTAINERS_RYUK_DISABLED=true \
  rtk go test -tags=dbtest -count=1 -timeout 300s ./service/expiry/...
```

Untuk Docker (bukan Podman), env var `TESTCONTAINERS_RYUK_DISABLED` tidak
perlu di-set; default Ryuk reaper bisa connect ke `/var/run/docker.sock`
tanpa perlu auth ulang.

Coverage target Layer 3:

```bash
TESTCONTAINERS_RYUK_DISABLED=true \
  rtk go test -tags=dbtest -coverprofile=cov.out \
    ./service/expiry/... ./api/sse/... ./internal/tcpmock/...
rtk go tool cover -func=cov.out | grep total
```

Target: ≥70% di `service/expiry`, `api/sse`, `internal/tcpmock`.

### Layer 4: integration ke router fisik

Build tag `integration`. Skip otomatis kalau `ROSLIB_ROUTER_ADDRESS` kosong.
**Test ini memodifikasi state router** (create/delete user, scheduler,
script). Konvensi: semua resource diawali `mikhmon-it-`, dibersihkan via
`t.Cleanup`.

```bash
cp .env.example .env
# isi ROSLIB_ROUTER_ADDRESS, USERNAME, PASSWORD
export $(grep -v '^#' .env | xargs)

# Smoke test minimum (tanpa DB):
rtk go test -tags=integration -count=1 ./test/integration/...

# Expiry workflow e2e (butuh router + Docker):
TESTCONTAINERS_RYUK_DISABLED=true \
  rtk go test -tags='integration dbtest' -count=1 \
    ./test/integration/...
```

`-count=1` mencegah Go cache hasil — wajib untuk integration test yang
state-nya bergantung pada router.

### Layer 4 verifikasi manual (out-of-scope untuk automated test)

Skenario yang tidak ter-cover automated karena fundamentally butuh
captive-portal client:

- Captive portal HTTP flow (form submit, redirect, cookie issuance).
- RADIUS Access-Request (butuh RADIUS server + client).
- Walled-garden DNS (butuh DNS resolver client).
- Bandwidth throttling enforcement (butuh real traffic).

Operator melakukan smoke test manual pre-release untuk skenario di atas.

## Mock Infrastructure

### `internal/tcpmock` — TCP server RouterOS API

Multi-conn TCP server yang berbicara protocol RouterOS API. Tiga mode
dispatch:

| Mode | API | Kegunaan |
|---|---|---|
| Matcher | `OnSentence(matcher, replies...)` | Reply identik tiap match, bagus untuk command idempoten |
| Matcher dinamis | `OnSentenceFunc(matcher, fn)` | Reply per-call berdasarkan counter/state (partial failure) |
| Streaming | `OnStream(matcher, taggedReplies...)` + `EmitToStream(tag, words...)` | Subscribe + push event tambahan |
| FIFO | `Script(replies...)` | Reply group berurutan (legacy / fallback) |

Sentence builder helpers di `builders.go`:
`DoneReply`, `ReReply`, `TrapReply`, `ActiveLoggedIn`, `ActiveLoggedOut`,
`UserPrint`, `UserPrintExpired`.

Assertion helpers (`assertions.go`):
`AssertReceived`, `AssertReceivedAll`, `AssertNotReceived` — polling-based,
race-aware.

### `internal/testutil` — test wiring

| Helper | Tujuan |
|---|---|
| `NewClient(t)` | dial router fisik via env (Layer 4) |
| `NewMockDevice(t)` | `*roslib.Device` ke tcpmock yang sudah AcceptLogin |
| `NewTestClientSet(t)` | wrap NewMockDevice + bangun `*devmgr.ClientSet` |
| `NewTestDevMgr(t)` | `*devmgr.Manager` dengan satu device dummy registered |
| `NewPostgres(t)` (tag `dbtest`) | postgres testcontainer + auto-migrate |
| `NewStores(t)` (tag `dbtest`) | NewPostgres + 3 store siap pakai |
| `Eventually(t, cond, timeout, msg)` | poll wrapper untuk eventual state |
| `RaceEnabled` const | true kalau binary build dengan `-race` |

## Catatan Race Detector

go-routeros v3.0.1 memiliki **race internal** antara `ctxReader.Close()`
dan `ctxReader.Cancel()` di `proto/io_context.go` (channel `close` + `send`
tanpa mutex). Race ini ter-flag walau goroutine sudah selesai.

Test yang exercise streaming atau async `!trap` skip saat `-race` via
`if testutil.RaceEnabled { t.Skip(...) }`. Fungsionalitasnya tetap
ter-cover saat run tanpa `-race`. Affected tests:

- `mikrotik/hotspot/active_stream_test.go` — semua 4 stream tests.
- `api/sse/broker_streaming_test.go` — 2 e2e SSE tests.
- `api/handler/hotspot_voucher_test.go` — `TestVoucher_partialFailure_returns207`.

Test pure (tcpmock self-test, testutil unit, broker fan-out unit) tetap
race-clean di semua kondisi.

## Golden File Update (`scripts/`)

Generator script (`scripts/onlogin`, `scripts/onevent`, `scripts/transaction`,
`scripts/quickprint`) di-test dengan golden file.

```bash
rtk go test ./scripts/... -update
git diff scripts/*/testdata/
git add scripts/*/testdata/ && git commit
```

Test tanpa `-update` FAIL kalau output beda dari golden.

## Coverage Aggregate

```bash
make test-cover
```

Atau manual:

```bash
TESTCONTAINERS_RYUK_DISABLED=true \
  rtk go test -tags=dbtest -coverprofile=cov.out -p 1 ./...
rtk go tool cover -html=cov.out -o cov.html
xdg-open cov.html
```

## Quick Reference

| Tujuan | Command |
|---|---|
| Smoke + unit, fast | `rtk go test -p 1 ./...` |
| + race detector | `rtk go test -p 1 -race ./...` |
| Layer 3 (DB) | `TESTCONTAINERS_RYUK_DISABLED=true rtk go test -tags=dbtest -p 1 ./...` |
| Layer 4 (router) | `rtk go test -tags=integration -count=1 ./test/integration/...` |
| Layer 4 full | `TESTCONTAINERS_RYUK_DISABLED=true rtk go test -tags='integration dbtest' -count=1 ./test/integration/...` |
| Update golden | `rtk go test ./scripts/... -update` |
