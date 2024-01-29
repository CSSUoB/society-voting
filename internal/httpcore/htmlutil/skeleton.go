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
	return HTML(
		Head(
			Meta(g.Attr("charset", "utf-8")),
			Meta(g.Attr("name", "viewport"), g.Attr("content", "width=device-width, initial-scale=1")),
			TitleEl(g.Text(title)),
			Link(g.Attr("rel", "stylesheet"), g.Attr("type", "text/css"), g.Attr("href", "https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css")),
			Link(g.Attr("rel", "icon"), g.Attr("href", "/favicon.svg"), g.Attr("type", "image/svg")),
			StyleEl(g.Raw(`body,html{height:100%}.form-signin{max-width:330px;padding:1rem}.form-signin .form-floating:focus-within{z-index:2}`)),
		),
		Body(
			g.Attr("class", "d-flex align-items-center py-4 bg-body-tertiary"),
			Main(
				g.Attr("class", "form-signin w-100 m-auto"),
				Img(g.Attr("src", "https://cssuob.github.io/resources/dinosaur/raster/tex_ballot.png"), g.Attr("style", "max-height: 72px; margin-bottom: 15px;")),
				g.Group(content),
				P(
					g.Attr("class", "mt-3 mb-3 text-body-secondary"),
					g.Text("Society Voting is open source software licensed under the BSD 2-Clause License. You can contribute "),
					A(g.Attr("href", "https://github.com/CSSUoB/society-voting"), g.Text("here")),
					g.Text("."),
				),
			),
		),
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
	return Div(
		g.Attr("class", "form-floating"),
		Input(
			g.Attr("type", inputType),
			g.Attr("placeholder", placeholder),
			g.Attr("name", name),
			g.Attr("id", r),
			g.Attr("class", "form-control"),
		),
		Label(
			g.Text(labelText),
			g.Attr("for", r),
		),
	)
}

func SmallTitle(content string) g.Node {
	return H2(
		g.Attr("class", "h5 mb-3 fw-normal"),
		g.Text(content),
	)
}

func FormSubmitButton() g.Node {
	return Button(
		g.Attr("class", "btn btn-primary w-100 py-2 mt-2"),
		g.Attr("type", "submit"),
		g.Text("Next"),
	)
}