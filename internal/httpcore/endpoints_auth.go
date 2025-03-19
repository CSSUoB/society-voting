package httpcore

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/guildScraper"
	"github.com/CSSUoB/society-voting/internal/httpcore/htmlutil"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
	"strings"
)

func (endpoints) authLoginPage(ctx *fiber.Ctx) error {
	page := htmlutil.SkeletonPage(config.Get().Platform.SocietyName+" voting",
		//g.If(requestProblem != "",
		//	html.P(g.Textf("Error: %s", requestProblem), g.Attr("style", "color: red")),
		//),
		html.Script(g.Attr("src", "https://unpkg.com/htmx.org@1.9.10"), g.Attr("defer")),
		html.H1(g.Attr("class", "h3 mb-3 fw-normal"), g.Text("Welcome!")),
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
		StudentID:            strings.TrimSpace(ctx.FormValue("studentid")),
		Password:             ctx.FormValue("password"),
		PasswordConfirmation: ctx.FormValue("passwordconf"),
		FirstName:            strings.TrimSpace(ctx.FormValue("fname")),
		LastName:             strings.TrimSpace(ctx.FormValue("lname")),
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

		// SID already registered?
		user, err := database.GetUser(requestData.StudentID)
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return fmt.Errorf("authLogin get user: %w", err)
		}
		if user != nil {
			// Yes: has password?
			if requestData.Password != "" {
				// Yes: Is password valid?
				match, err := argon2id.ComparePasswordAndHash(requestData.Password, user.PasswordHash)
				if err != nil {
					return fmt.Errorf("authLogin compare password hash: %w", err)
				}

				if match {
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

		if subtle.ConstantTimeCompare([]byte(strings.ToLower(requestData.FirstName)), []byte(strings.ToLower(guildMember.FirstName))) == 0 ||
			subtle.ConstantTimeCompare([]byte(strings.ToLower(requestData.LastName)), []byte(strings.ToLower(guildMember.LastName))) == 0 {
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

		// Has password and password confirmation?
		if requestData.PasswordConfirmation == "" || requestData.Password == "" {
			// No: show form
			goto unregisteredAskPassword
		}

		if requestData.PasswordConfirmation != requestData.Password {
			requestProblem = "Passwords do not match."

			// Since unregisteredAskPassword still includes previous request data, we need to remove the old passwords to prevent them from overriding the new passwords the user will input.
			// If this were not done, a user that entered an non-matching password pair would never be able to set their password.

			requestData.Password = ""
			requestData.PasswordConfirmation = ""

			requestDataJSON, err = json.Marshal(&requestData)
			if err != nil {
				return fmt.Errorf("authLogin marshal request data to JSON after removing passwords: %w", err)
			}

			goto unregisteredAskPassword
		}
		var passwordHash string
		if requestData.Password != "" {
			var err error
			passwordHash, err = argon2id.CreateHash(requestData.Password, argon2id.DefaultParams)
			if err != nil {
				return fmt.Errorf("authLogin hash password: %w", err)
			}
		}

		user = &database.User{
			StudentID:    requestData.StudentID,
			Name:         guildMember.FirstName + " " + guildMember.LastName,
			PasswordHash: passwordHash,
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
		htmlutil.SmallTitle("What's your student ID?"),
		htmlutil.FormInput("text", "studentid", "Your student ID", "Student ID"),
		htmlutil.FormSubmitButton(),
	))

registeredAskPassword:
	return htmlutil.SendFragment(ctx, html.FormEl(
		g.Attr("hx-indicator", "#indicator"),
		g.Attr("hx-post", loginActionEndpoint),
		g.Attr("hx-swap", "outerHTML"),
		g.Attr("hx-vals", string(requestDataJSON)),
		html.P(g.Text("Please enter your password. If you've forgotten it, please speak to a member of committee.")),
		htmlutil.FormInput("password", "password", "", "Password"),
		htmlutil.FormSubmitButton(),
	))

askNames:
	return htmlutil.SendFragment(ctx, html.FormEl(
		g.Attr("hx-indicator", "#indicator"),
		g.Attr("hx-post", loginActionEndpoint),
		g.Attr("hx-swap", "outerHTML"),
		g.Attr("hx-vals", string(requestDataJSON)),
		html.P(g.Text("In order to check that the student ID you entered is yours, please enter your name as the Guild of Students knows it.")),
		htmlutil.FormInput("text", "fname", "Your first name", "First name"),
		htmlutil.FormInput("text", "lname", "Your last name", "Last name"),
		htmlutil.FormSubmitButton(),
	))

askAuthCode:
	return htmlutil.SendFragment(ctx, html.FormEl(
		g.Attr("hx-indicator", "#indicator"),
		g.Attr("hx-post", loginActionEndpoint),
		g.Attr("hx-swap", "outerHTML"),
		g.Attr("hx-vals", string(requestDataJSON)),
		html.P(g.Text("Please enter the authorisation code.")),
		htmlutil.FormInput("password", "auth", "", "Authorisation code"),
		htmlutil.FormSubmitButton(),
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
		html.P(g.Text("Please choose a password.")),
		htmlutil.FormInput("password", "password", "", "Password"),
		htmlutil.FormInput("password", "passwordconf", "", "Confirm password"),
		htmlutil.FormSubmitButton(),
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
