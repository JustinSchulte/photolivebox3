// Package frontend serves the nuxt frontend
package frontend

import (
	"embed"
	"io/fs"

	"github.com/labstack/echo/v4"
)

//go:embed all:app/.output/public
var files embed.FS

func StaticContent() fs.FS {
	filesystem := echo.MustSubFS(files, "app/.output/public")
	return filesystem
}
