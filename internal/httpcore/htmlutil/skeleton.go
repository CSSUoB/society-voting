package htmlutil

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"math/rand"
	"strconv"
)

func SkeletonPage(title string, content ...g.Node) g.Node {
	content = append(content, P(g.Attr("style", "color: grey; font-style: italic;"),
		g.Text("Society Voting is open source software licensed under the BSD 2-Clause License. You can contribute "),
		A(g.Attr("href", "https://github.com/CSSUoB/society-voting"), g.Text("here")),
		g.Text("."),
	))

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
	buf.WriteString("<!DOCTYPE html>")
	_ = page.Render(buf)
	ctx.Type("html")
	return ctx.SendStream(buf)
}

func SendFragment(ctx *fiber.Ctx, frag g.Node) error {
	buf := new(bytes.Buffer)
	_ = frag.Render(buf)
	ctx.Type("html")
	return ctx.SendStream(buf)
}

func FormInput(inputType, name, placeholder, labelText string) g.Node {
	r := strconv.Itoa(rand.Intn(1000000))
	return g.Group([]g.Node{
		Label(
			g.Text(labelText),
			g.Attr("for", r),
			g.Attr("style", "margin-right: 6px"),
		),
		Input(
			g.Attr("type", inputType),
			g.Attr("placeholder", placeholder),
			g.Attr("name", name),
			g.Attr("id", r),
		),
	})
}
