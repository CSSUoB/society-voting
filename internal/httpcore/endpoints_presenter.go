package httpcore

import (
	_ "embed"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

func (endpoints) presenterQRCode(ctx *fiber.Ctx) error {
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
	url := "https://" + ctx.Get("Host")
	ctx.Type("html")
	return ctx.SendString(
		strings.ReplaceAll(strings.ReplaceAll(presenterPageHTML, "{{url}}", url), "{{votecode}}", voteCode),
	)
}
