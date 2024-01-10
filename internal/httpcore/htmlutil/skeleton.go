package htmlutil

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func SkeletonPage(title string, content ...g.Node) g.Node {
	return HTML(
		Head(
			Meta(g.Attr("charset", "utf-8")),
			Meta(g.Attr("name", "viewport"), g.Attr("content", "width=device-width, initial-scale=1.0")),
			TitleEl(g.Text(title)),
			StyleEl(g.Text(`
body {
	font-family: sans-serif;
	font-size: 1.1rem;
	padding: 1em;
}
`)),
		),
		Body(content...),
	)
}

func SendPage(ctx *fiber.Ctx, page g.Node) error {
	buf := new(bytes.Buffer)
	_ = page.Render(buf)
	ctx.Type("html")
	return ctx.SendStream(buf)
}
