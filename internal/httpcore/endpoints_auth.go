package httpcore

import (
	"crypto/sha512"
	"crypto/subtle"
	"errors"
	"fmt"
	"git.tdpain.net/codemicro/society-voting/internal/config"
	"git.tdpain.net/codemicro/society-voting/internal/database"
	"git.tdpain.net/codemicro/society-voting/internal/guildScraper"
	"git.tdpain.net/codemicro/society-voting/internal/httpcore/htmlutil"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
	"strconv"
)

func (endpoints) authCheck(ctx *fiber.Ctx) error {
	_, isAuthed, err := getSessionAuth(ctx)
	if err != nil {
		return err
	}

	var nextURL string

	if isAuthed {
		// redirect to app
		nextURL = "/app"
	} else {
		// redirect to login page
		nextURL = "/auth/login"
	}

	return ctx.Redirect(nextURL)
}

func (endpoints) authLogin(ctx *fiber.Ctx) error {
	var requestProblem string

	switch ctx.Method() {
	case fiber.MethodGet:
		goto staticPage
	case fiber.MethodPost:
		break
	default:
		return fiber.ErrMethodNotAllowed
	}

	if method := ctx.Method(); method == fiber.MethodGet {
		goto staticPage
	} else if method == fiber.MethodPost {
		// Process form with <studentid> field
		studentID := ctx.FormValue("studentid")
		if studentID == "" {
			requestProblem = "missing student ID"
			goto staticPage
		}

		if _, err := strconv.Atoi(studentID); err != nil && studentID != "admin" {
			requestProblem = "invalid student ID"
			goto staticPage
		}

		passwordPlaintext := ctx.FormValue("password")
		if passwordPlaintext == "" {
			requestProblem = "missing password"
			goto staticPage
		}

		if studentID == "admin" {
			if subtle.ConstantTimeCompare([]byte(config.Get().Platform.AdminToken), []byte(passwordPlaintext)) == 0 {
				goto incorrectPassword
			} else {
				ctx.Cookie(newSessionTokenCookie(signData("admin", "admin")))
				return ctx.Redirect("/app")
			}
		}

		// Provision user if needed

		user, err := database.GetUser(studentID)
		if err != nil {
			if !errors.Is(err, database.ErrNotFound) {
				return fmt.Errorf("authLogin call getUser: %w", err)
			}
		}

		passwordHash := sha512.Sum512([]byte(passwordPlaintext))

		if user == nil {
			// provision user

			guildMember, err := guildScraper.GetMember(studentID)
			if err != nil {
				return fmt.Errorf("authLogin membership check: %w", err)
			}

			if guildMember == nil {
				requestProblem = "not a member of " + config.Get().Platform.SocietyName
				goto staticPage
			}

			user = &database.User{
				StudentID:    studentID,
				Name:         guildMember.Name,
				PasswordHash: passwordHash[:],
			}

			if err := user.Insert(); err != nil {
				return fmt.Errorf("authLogin insert new user: %w", err)
			}
		} else {
			if subtle.ConstantTimeCompare(passwordHash[:], user.PasswordHash) == 0 {
				goto incorrectPassword
			}
		}

		// Issue token
		ctx.Cookie(newSessionTokenCookie(signData("token", studentID)))

		// Redirect to app
		return ctx.Redirect("/app")
	} else {
		return fiber.ErrMethodNotAllowed
	}

incorrectPassword:
	requestProblem = "incorrect password (speak to a member of committee if you've forgotten your password)"

staticPage:
	// GET request: serve page

	titleLine := config.Get().Platform.SocietyName + " voting"
	page := htmlutil.SkeletonPage(titleLine,
		html.H1(g.Text(titleLine)),
		html.P(g.Text("Please enter your details. Please note you must be a member to proceed.")),
		html.P(g.Text("If you have not logged in before, please choose a password and enter it here, otherwise use your existing password.")),
		g.If(requestProblem != "",
			html.P(g.Textf("Error: %s", requestProblem), g.Attr("style", "color: red")),
		),
		html.FormEl(
			g.Attr("method", "POST"),
			html.Label(
				g.Text("Student ID"),
				g.Attr("for", "student-id-input"),
				g.Attr("style", "margin-right: 6px"),
			),
			html.Input(
				g.Attr("type", "text"),
				g.Attr("id", "student-id-input"),
				g.Attr("placeholder", "Student ID"),
				g.Attr("name", "studentid"),
			),
			html.Br(),
			html.Label(
				g.Text("Password"),
				g.Attr("for", "password-input"),
				g.Attr("style", "margin-right: 6px"),
			),
			html.Input(
				g.Attr("type", "password"),
				g.Attr("placeholder", "Password"),
				g.Attr("name", "password"),
			),
			html.Br(),
			html.Input(
				g.Attr("type", "submit"),
				g.Attr("value", "Submit"),
			),
		),
	)

	return htmlutil.SendPage(ctx, page)
}