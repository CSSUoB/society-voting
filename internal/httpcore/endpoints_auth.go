package httpcore

import (
	"crypto/sha512"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/guildScraper"
	"github.com/CSSUoB/society-voting/internal/httpcore/htmlutil"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

func (endpoints) authLoginPage(ctx *fiber.Ctx) error {
	titleLine := config.Get().Platform.SocietyName + " voting"
	page := htmlutil.SkeletonPage(titleLine,
		html.H1(g.Text(titleLine)),
		html.P(g.Text("Please enter your details. Please note you must be a member to proceed.")),
		//g.If(requestProblem != "",
		//	html.P(g.Textf("Error: %s", requestProblem), g.Attr("style", "color: red")),
		//),
		html.Script(g.Attr("src", "https://unpkg.com/htmx.org@1.9.10"), g.Attr("defer")),
		html.Div(
			g.Attr("hx-get", loginActionEndpoint),
			g.Attr("hx-trigger", "load"),
			g.Attr("hx-swap", "outerHTML"),
			//g.Text("Loading..."),
		),
		html.P(g.Attr("class", "htmx-indicator"), g.Attr("id", "indicator"), g.Text("Working...")),
	)

	return htmlutil.SendPage(ctx, page)
}

func (endpoints) authLogin(ctx *fiber.Ctx) error {
	var requestProblem string

	var requestData = struct {
		StudentID            string `json:"studentid,omitempty"`
		Password             string `json:"password,omitempty"`
		PasswordConfirmation string `json:"passwordconf,omitempty"`
		FirstName            string `json:"fname,omitempty"`
		LastName             string `json:"lname,omitempty"`
		AuthCode             string `json:"auth,omitempty"`
	}{
		StudentID:            ctx.FormValue("studentid"),
		Password:             ctx.FormValue("password"),
		PasswordConfirmation: ctx.FormValue("passwordconf"),
		FirstName:            ctx.FormValue("fname"),
		LastName:             ctx.FormValue("lname"),
		AuthCode:             ctx.FormValue("auth"),
	}

	requestDataJSON, err := json.Marshal(&requestData)
	if err != nil {
		return fmt.Errorf("authLogin marshal request data to JSON: %w", err)
	}

	if ctx.Method() == fiber.MethodGet {
		goto askStudentID
	}

	{
		// Has SID?
		if requestData.StudentID == "" {
			// No: show form
			goto askStudentID
		}

		var passwordHash [64]byte
		if requestData.Password != "" {
			passwordHash = sha512.Sum512([]byte(requestData.Password))
		}

		// SID already registered?
		user, err := database.GetUser(requestData.StudentID)
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return fmt.Errorf("authLogin get user: %w", err)
		}
		if user != nil {
			// Yes: has password?
			if requestData.Password != "" {
				// Yes: Is password valid?
				if subtle.ConstantTimeCompare(passwordHash[:], user.PasswordHash) == 1 {
					// Yes: issue token, redirect
					goto success
				}

				// No: return to start with error
				requestProblem = "Invalid password - that password did not match the password that is stored for that student ID."
				goto reset
			} else {
				// No: show form
				goto registeredAskPassword
			}
		}

		guildMember, err := guildScraper.GetMember(requestData.StudentID)
		if err != nil {
			return fmt.Errorf("authLogin get guild member: %w", err)
		}
		if guildMember == nil {
			requestProblem = "Invalid student ID - it doesn't look like that student ID corresponds to a " + config.Get().Platform.SocietyName + " member."
			goto reset
		}

		// Has names?
		if requestData.FirstName == "" && requestData.LastName == "" {
			// No: show names form
			goto askNames
		}

		fmt.Println(requestData.FirstName, guildMember.FirstName)
		fmt.Println(requestData.LastName, guildMember.LastName)

		if subtle.ConstantTimeCompare([]byte(requestData.FirstName), []byte(guildMember.FirstName)) == 0 ||
			subtle.ConstantTimeCompare([]byte(requestData.LastName), []byte(guildMember.LastName)) == 0 {
			requestProblem = "Invalid name - please ensure that you are using the name the Guild of Students knows you as."
			goto reset
		}

		// Will be an admin?
		if guildMember.IsCommitteeMember {
			// Yes: has auth code?
			if requestData.AuthCode == "" {
				// No: show form
				goto askAuthCode
			}

			if subtle.ConstantTimeCompare([]byte(requestData.AuthCode), []byte(config.Get().Platform.AdminToken)) == 0 {
				requestProblem = "Invalid admin token."
				goto reset
			}
		}

		fmt.Printf("%#v\n", requestData)

		// Has password and password confirmation?
		if requestData.PasswordConfirmation == "" || requestData.Password == "" {
			// No: show form
			goto unregisteredAskPassword
		}

		if requestData.PasswordConfirmation != requestData.Password {
			requestProblem = "Passwords do not match."
			goto unregisteredAskPassword
		}

		user = &database.User{
			StudentID:    requestData.StudentID,
			Name:         guildMember.FirstName + " " + guildMember.LastName,
			PasswordHash: passwordHash[:],
			IsAdmin:      guildMember.IsCommitteeMember,
		}

		if err := user.Insert(); err != nil {
			return fmt.Errorf("authLogin insert new user: %w", err)
		}
	}

success:
	ctx.Cookie(newSessionTokenCookie(signData(requestData.StudentID)))
	ctx.Set("HX-Redirect", "/")
	ctx.Status(fiber.StatusNoContent)
	return nil

reset:
askStudentID:
	return htmlutil.SendFragment(ctx, html.FormEl(
		g.Attr("hx-indicator", "#indicator"),
		g.Attr("hx-post", loginActionEndpoint),
		g.Attr("hx-swap", "outerHTML"),
		g.If(requestProblem != "", html.P(
			g.Attr("style", "color: red;"),
			html.Em(g.Text(requestProblem+" If you're having trouble logging in, please speak to a member of committee.")),
		)),
		html.H2(g.Text("What's your student ID?")),
		htmlutil.FormInput("text", "studentid", "Your student ID", "Student ID"),
		html.Br(),
		html.Input(
			g.Attr("type", "submit"),
			g.Attr("value", "Next"),
		),
	))

registeredAskPassword:
	return htmlutil.SendFragment(ctx, html.FormEl(
		g.Attr("hx-indicator", "#indicator"),
		g.Attr("hx-post", loginActionEndpoint),
		g.Attr("hx-swap", "outerHTML"),
		g.Attr("hx-vals", string(requestDataJSON)),
		html.H2(g.Text("Welcome back!")),
		html.P(g.Text("Please enter your password. If you've forgotten it, please speak to a member of committee.")),
		htmlutil.FormInput("password", "password", "", "Password"),
		html.Br(),
		html.Input(
			g.Attr("type", "submit"),
			g.Attr("value", "Next"),
		),
	))

askNames:
	return htmlutil.SendFragment(ctx, html.FormEl(
		g.Attr("hx-indicator", "#indicator"),
		g.Attr("hx-post", loginActionEndpoint),
		g.Attr("hx-swap", "outerHTML"),
		g.Attr("hx-vals", string(requestDataJSON)),
		html.H2(g.Text("Welcome!")),
		html.P(g.Text("In order to check that the student ID you entered is yours, please enter your name as the Guild of Students knows it.")),
		htmlutil.FormInput("text", "fname", "Your first name", "First name"),
		html.Br(),
		htmlutil.FormInput("text", "lname", "Your last name", "Last name"),
		html.Br(),
		html.Input(
			g.Attr("type", "submit"),
			g.Attr("value", "Next"),
		),
	))

askAuthCode:
	return htmlutil.SendFragment(ctx, html.FormEl(
		g.Attr("hx-indicator", "#indicator"),
		g.Attr("hx-post", loginActionEndpoint),
		g.Attr("hx-swap", "outerHTML"),
		g.Attr("hx-vals", string(requestDataJSON)),
		html.H2(g.Text("Welcome!")),
		html.P(g.Text("Please enter the authorisation code.")),
		htmlutil.FormInput("text", "auth", "", "Authorisation code"),
		html.Br(),
		html.Input(
			g.Attr("type", "submit"),
			g.Attr("value", "Next"),
		),
	))

unregisteredAskPassword:
	return htmlutil.SendFragment(ctx, html.FormEl(
		g.Attr("hx-indicator", "#indicator"),
		g.Attr("hx-post", loginActionEndpoint),
		g.Attr("hx-swap", "outerHTML"),
		g.Attr("hx-vals", string(requestDataJSON)),
		g.If(requestProblem != "", html.P(
			g.Attr("style", "color: red;"),
			html.Em(g.Text(requestProblem)),
		)),
		html.H2(g.Text("Welcome!")),
		html.P(g.Text("Please choose a password.")),
		htmlutil.FormInput("password", "password", "", "Password"),
		html.Br(),
		htmlutil.FormInput("password", "passwordconf", "", "Confirm password"),
		html.Br(),
		html.Input(
			g.Attr("type", "submit"),
			g.Attr("value", "Next"),
		),
	))

}

func (endpoints) authLogout(ctx *fiber.Ctx) error {
	ctx.Cookie(newSessionTokenDeletionCookie())
	titleLine := config.Get().Platform.SocietyName + " voting"

	return htmlutil.SendPage(ctx, htmlutil.SkeletonPage(
		titleLine,
		html.H1(g.Text("You're all signed out!")),
		html.A(g.Attr("href", "/auth/login"), g.Text("Click here to login again")),
	))
}
