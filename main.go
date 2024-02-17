package main

import (
	"github.com/JustinSchulte/photolivebox3/backend"
	"github.com/JustinSchulte/photolivebox3/frontend"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	// serve frontend
	e.StaticFS("/", frontend.StaticContent())
	// serve backend endpoints
	e.GET("/api/v1/getImage/:filename", backend.GetImage)
	e.GET("/api/v1/listImages", backend.ListImages)
	e.POST("/api/v1/upload", backend.UploadImage)
	// start echo server at the port
	e.Start(":80")
}
