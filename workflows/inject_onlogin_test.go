package workflows

import "testing"

// TestBuildWebhookURL verifikasi format URL yang di-construct cocok dengan
// route yang di-mount di api/routes.go: /api/v1/hook/hotspot/login/<id>.
//
// Trailing slash di base URL harus di-strip supaya tidak jadi double-slash
// (mis. "http://x:8080/" + "/api/v1/..." = "http://x:8080//api/v1/...").
func TestBuildWebhookURL(t *testing.T) {
	cases := []struct {
		name      string
		base      string
		deviceID  uint
		want      string
	}{
		{
			name:     "empty base returns empty",
			base:     "",
			deviceID: 1,
			want:     "",
		},
		{
			name:     "no trailing slash",
			base:     "http://192.168.1.10:8080",
			deviceID: 1,
			want:     "http://192.168.1.10:8080/api/v1/hook/hotspot/login/1",
		},
		{
			name:     "single trailing slash stripped",
			base:     "http://192.168.1.10:8080/",
			deviceID: 7,
			want:     "http://192.168.1.10:8080/api/v1/hook/hotspot/login/7",
		},
		{
			name:     "multiple trailing slashes stripped",
			base:     "http://x:8080///",
			deviceID: 42,
			want:     "http://x:8080/api/v1/hook/hotspot/login/42",
		},
		{
			name:     "https scheme",
			base:     "https://api.example.com",
			deviceID: 100,
			want:     "https://api.example.com/api/v1/hook/hotspot/login/100",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := buildWebhookURL(tt.base, tt.deviceID)
			if got != tt.want {
				t.Errorf("buildWebhookURL(%q, %d) = %q, want %q",
					tt.base, tt.deviceID, got, tt.want)
			}
		})
	}
}

// TestErrUnwrap memastikan inner-most error di-extract.
func TestErrUnwrap(t *testing.T) {
	innerMost := errSentinel("inner")
	wrapped := wrapErr(wrapErr(innerMost))
	got := errUnwrap(wrapped)
	if got.Error() != innerMost.Error() {
		t.Errorf("errUnwrap should return inner-most, got %q, want %q",
			got.Error(), innerMost.Error())
	}
}

// errSentinel + wrapErr helpers untuk test.
type errSentinel string

func (e errSentinel) Error() string { return string(e) }

type wrappedErr struct{ inner error }

func (w *wrappedErr) Error() string { return "wrap: " + w.inner.Error() }
func (w *wrappedErr) Unwrap() error { return w.inner }

func wrapErr(err error) error { return &wrappedErr{inner: err} }
