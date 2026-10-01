// Package semgreprules embeds the semgrep rules run by the go-semgrep hook.
package semgreprules

import "embed"

//go:embed *.yml
var FS embed.FS
