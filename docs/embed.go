// Package docs embeds the bundled OpenAPI spec so the server binary is
// self-contained. Source files live under openapi/ and are bundled via
// `make openapi-bundle` (Redocly CLI). The bundle is committed and CI
// verifies it stays in sync.
package docs

import _ "embed"

//go:embed openapi/openapi.bundle.yaml
var OpenAPIBundle []byte

//go:embed scalar/index.html
var ScalarIndexHTML []byte
