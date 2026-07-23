// go-snappymail — SnappyMail UX webmail in Go (IMAP/SMTP).
package main

import (
	"embed"

	"go-snappymail/cmd"
)

//go:embed all:web/dist
//go:embed all:web/files
var embeddedFiles embed.FS

func main() {
	cmd.Execute(embeddedFiles)
}
