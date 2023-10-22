package web

import (
	"embed"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"net/http"
)

//go:generate npm install
//go:generate npm run build

//go:embed dist/*
var fs embed.FS

func GetHandler() fiber.Handler {
	return filesystem.New(filesystem.Config{
		Root:       http.FS(fs),
		PathPrefix: "dist",
	})
}
