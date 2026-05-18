//go:build dbtest

package expiry

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/quiqxiq/roslib-mikhmon/internal/tcpmock"
	"github.com/quiqxiq/roslib-mikhmon/internal/testutil"
	"github.com/quiqxiq/roslib-mikhmon/service/devmgr"
	"github.com/quiqxiq/roslib-mikhmon/store/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test ini exercise jalur checkDevice + executeExpiry end-to-end:
//   - postgres testcontainer untuk store seed
//   - tcpmock untuk router replies
//   - white-box test (package expiry) supaya bisa panggil checkDevice
//     yang unexported.

// setupExpiryEnv menyiapkan service + mock + dev model siap di-test.
func setupExpiryEnv(t *testing.T, mode string) (*Service, *tcpmock.Server, model.MikrotikDevice) {
	t.Helper()

	devStore, profStore, txStore := testutil.NewStores(t)
	mgr, srv, dev := testutil.NewTestDevMgr(t)

	// Seed device row (ID harus match dev.ID = 1 dari NewTestDevMgr).
	ctx := context.Background()
	err := devStore.Create(ctx, &dev)
	require.NoError(t, err)

	// Seed profile config dengan mode yang diuji.
	require.NoError(t, profStore.Upsert(ctx, &model.HotspotProfileConfig{
		DeviceID:    dev.ID,
		ProfileName: "default",
		ExpiryMode:  mode,
		Price:       10000,
		SellPrice:   12000,
		Validity:    "30d",
	}))

	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)

	svc := New(mgr, devStore, profStore, txStore, log)
	svc.rootCtx = ctx
	return svc, srv, dev
}

// (note: tests pakai OnSentence matcher-based dispatch supaya cocok
//  per-command-path. Pakai -48h offset supaya past time tetap past walau
//  ParseExpiry mengembalikan UTC sementara time.Now() bisa local TZ.)

func TestService_remMode_userDeleted(t *testing.T) {
	svc, srv, dev := setupExpiryEnv(t, "rem")

	past := time.Now().Add(-48 * time.Hour)
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/user/print"),
		tcpmock.UserPrintExpired("1", "alice", "default", past),
		tcpmock.DoneReply(),
	)
	srv.OnSentence(tcpmock.MatchCommand("/system/script/print"), tcpmock.DoneReply())
	srv.OnSentence(tcpmock.MatchCommand("/system/scheduler/print"), tcpmock.DoneReply())
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/user/remove"), tcpmock.DoneReply())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, svc.checkDevice(ctx, dev))

	srv.AssertReceived(t, tcpmock.MatchAll(
		tcpmock.MatchCommand("/ip/hotspot/user/remove"),
		tcpmock.MatchHas("numbers", "*1"),
	), "user/remove =numbers=*1")
}

func TestService_remcMode_recordsTransaction(t *testing.T) {
	svc, srv, dev := setupExpiryEnv(t, "remc")

	past := time.Now().Add(-48 * time.Hour)
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/user/print"),
		tcpmock.UserPrintExpired("1", "alice", "default", past),
		tcpmock.DoneReply(),
	)
	srv.OnSentence(tcpmock.MatchCommand("/system/script/print"), tcpmock.DoneReply())
	srv.OnSentence(tcpmock.MatchCommand("/system/scheduler/print"), tcpmock.DoneReply())
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/user/remove"), tcpmock.DoneReply())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, svc.checkDevice(ctx, dev))

	month := strings.ToLower(time.Now().Format("Jan2006"))
	txs, err := svc.txStore.ListByDevice(ctx, dev.ID, month)
	require.NoError(t, err)
	require.Len(t, txs, 1)
	assert.Equal(t, "alice", txs[0].Username)
	assert.Equal(t, 10000, txs[0].Price)
	assert.Equal(t, 12000, txs[0].SellPrice)
	assert.Equal(t, "default", txs[0].Profile)
}

func TestService_ntfMode_setsLimitAndKicks(t *testing.T) {
	svc, srv, dev := setupExpiryEnv(t, "ntf")

	past := time.Now().Add(-48 * time.Hour)
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/user/print"),
		tcpmock.UserPrintExpired("1", "alice", "default", past),
		tcpmock.DoneReply(),
	)
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/user/set"), tcpmock.DoneReply())
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/active/print"),
		tcpmock.ActiveLoggedIn("99", "alice", "10.0.0.5", "aa:bb:cc:dd:ee:ff"),
		tcpmock.DoneReply(),
	)
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/active/remove"), tcpmock.DoneReply())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, svc.checkDevice(ctx, dev))

	srv.AssertReceived(t, tcpmock.MatchAll(
		tcpmock.MatchCommand("/ip/hotspot/user/set"),
		tcpmock.MatchHas("limit-uptime", "1s"),
	), "user/set =limit-uptime=1s")
	srv.AssertReceived(t, tcpmock.MatchCommand("/ip/hotspot/active/remove"), "active/remove session")
}

func TestService_ntfMode_noTransactionRecorded(t *testing.T) {
	svc, srv, dev := setupExpiryEnv(t, "ntf")

	past := time.Now().Add(-48 * time.Hour)
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/user/print"),
		tcpmock.UserPrintExpired("1", "alice", "default", past),
		tcpmock.DoneReply(),
	)
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/user/set"), tcpmock.DoneReply())
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/active/print"), tcpmock.DoneReply())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, svc.checkDevice(ctx, dev))

	month := strings.ToLower(time.Now().Format("Jan2006"))
	txs, err := svc.txStore.ListByDevice(ctx, dev.ID, month)
	require.NoError(t, err)
	assert.Empty(t, txs, "ntf mode should not record transaction")
}

func TestService_invalidComment_skip(t *testing.T) {
	svc, srv, dev := setupExpiryEnv(t, "rem")

	// User dengan comment bukan format expiry — harus di-skip.
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/user/print"),
		tcpmock.UserPrint("1", "alice", "default", "not-a-date"),
		tcpmock.DoneReply(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, svc.checkDevice(ctx, dev))

	srv.AssertNotReceived(t, tcpmock.MatchCommand("/ip/hotspot/user/remove"))
	srv.AssertNotReceived(t, tcpmock.MatchCommand("/ip/hotspot/user/set"))
}

func TestService_validComment_futureExpiry_skip(t *testing.T) {
	svc, srv, dev := setupExpiryEnv(t, "rem")

	future := time.Now().Add(72 * time.Hour)
	srv.OnSentence(tcpmock.MatchCommand("/ip/hotspot/user/print"),
		tcpmock.UserPrintExpired("1", "alice", "default", future),
		tcpmock.DoneReply(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, svc.checkDevice(ctx, dev))

	srv.AssertNotReceived(t, tcpmock.MatchCommand("/ip/hotspot/user/remove"))
}

func TestService_backoff_deviceNotConnected_returnsErr(t *testing.T) {
	// Manager tanpa device terdaftar → Get return ErrDeviceNotConnected.
	devStore, profStore, txStore := testutil.NewStores(t)
	log := logrus.New()
	log.SetLevel(logrus.WarnLevel)

	// New manager, no device registered.
	emptyMgr := devmgr.New(devStore, log)
	svc := New(emptyMgr, devStore, profStore, txStore, log)
	svc.rootCtx = context.Background()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := svc.checkDevice(ctx, model.MikrotikDevice{ID: 99, Slug: "missing"})

	// runChecker akan switch ke backoff atas error ini; di sini kita verifikasi
	// surface ErrDeviceNotConnected langsung dari checkDevice.
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

