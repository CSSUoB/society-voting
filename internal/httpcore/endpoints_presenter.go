package httpcore

import (
	_ "embed"
	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"strings"
)

func (endpoints) presenterQRCode(ctx *fiber.Ctx) error {
	if _, status := getSessionAuth(ctx); status&authAdminUser == 0 {
		return fiber.ErrUnauthorized
	}

	var png []byte
	png, err := qrcode.Encode("https://"+ctx.Get("Host"), qrcode.Medium, 600)
	if err != nil {
		return err
	}

	ctx.Type("png")
	return ctx.Send(png)
}

//go:embed presenter.html
var presenterPageHTML string

func (endpoints) presenterPage(ctx *fiber.Ctx) error {
	if _, status := getSessionAuth(ctx); status&authAdminUser == 0 {
		return fiber.ErrUnauthorized
	}

	url := "https://" + ctx.Get("Host")
	ctx.Type("html")
	return ctx.SendString(
		strings.ReplaceAll(strings.ReplaceAll(presenterPageHTML, "{{url}}", url), "{{votecode}}", voteCode),
	)
}
