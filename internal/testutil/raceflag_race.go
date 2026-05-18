//go:build race

package testutil

// RaceEnabled true kalau binary di-compile dengan -race. Test yang trigger
// race di third-party lib (mis. go-routeros v3.0.1 race antara
// ctxReader.Close() dan ctxReader.Cancel() di proto/io_context.go) dapat
// skip dirinya kalau RaceEnabled.
const RaceEnabled = true
