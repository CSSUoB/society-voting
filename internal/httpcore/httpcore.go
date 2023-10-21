package httpcore

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"git.tdpain.net/codemicro/society-voting/internal/config"
	"git.tdpain.net/codemicro/society-voting/internal/database"
	"github.com/bwmarrin/go-alone"
	validate "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"os"
	"strings"
	"time"
)

type endpoints struct{}

func ListenAndServe(addr string) error {
	if config.Get().Debug {
		slog.Debug("debug token", "token", signData("token", "1111111"))
	}

	app := fiber.New()

	e := endpoints{}

	app.Get("/", e.authCheck)

	app.Get("/auth/login", e.authLogin)
	app.Post("/auth/login", e.authLogin)

	app.Get("/api/me", e.apiMe)
	app.Put("/api/me/name", e.apiSetOwnName)
	app.Post("/api/me/stand", e.apiStandForElection)
	app.Delete("/api/me/stand", e.apiWithdrawFromElection)

	app.Get("/api/elections", e.apiListElections)
	app.Get("/api/elections/sse", e.apiElectionsSSE)

	app.Post("/api/admin/election", e.apiAdminCreateElection)
	app.Delete("/api/admin/election", e.apiAdminDeleteElection)

	slog.Info("HTTP server alive", "address", addr)
	return app.Listen(addr)
}

var (
	signer    *goalone.Sword
	validator = validate.New(validate.WithRequiredStructEnabled())
)

func init() {
	secret := make([]byte, 512)
	if !config.Get().Debug {
		// This is so that the access tokens doesn't change from run-to-run for ease of testing
		if _, err := rand.Read(secret); err != nil {
			slog.Error("unable to generate random secret for token signing", "error", err)
			os.Exit(1)
		}
	}

	signer = goalone.New(secret)
}

var sessionTokenCookieName = "vot-tok"

func newSessionTokenCookie(tok string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:        sessionTokenCookieName,
		Value:       tok,
		Secure:      true,
		SessionOnly: true,
	}
}

func newSessionTokenDeletionCookie() *fiber.Cookie {
	return &fiber.Cookie{
		Name:        sessionTokenCookieName,
		Secure:      true,
		SessionOnly: true,
		Expires:     time.Date(2010, 1, 1, 1, 1, 1, 1, time.UTC),
	}
}

func signData(datatype string, data string) string {
	return hex.EncodeToString(signer.Sign([]byte(datatype + "." + data)))
}

func validateData(datatype string, token string) (string, bool) {
	dat, err := hex.DecodeString(token)
	if err != nil {
		return "", false
	}
	originalData, err := signer.Unsign(dat)
	if err != nil {
		return "", false
	}
	if !strings.HasPrefix(string(originalData), datatype+".") {
		return "", false
	}
	return string(originalData[len(datatype)+1:]), true
}

func getSessionAuth(ctx *fiber.Ctx) (*database.User, bool, error) {
	cookieVal := ctx.Cookies(sessionTokenCookieName)
	decodedVal, tokenOkay := validateData("token", cookieVal)

	if cookieVal == "" || !tokenOkay {
		return nil, false, nil
	}

	user, err := database.GetUser(decodedVal)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			// If the user has been deleted
			ctx.Cookie(newSessionTokenDeletionCookie())
			return nil, false, nil
		}
		return nil, false, err
	}

	return user, true, nil
}

func isAdminSession(ctx *fiber.Ctx) bool {
	cookieVal := ctx.Cookies(sessionTokenCookieName)
	_, tokenOkay := validateData("admin", cookieVal)

	if cookieVal == "" || !tokenOkay {
		return false
	}

	return true
}

func parseAndValidateRequestBody(ctx *fiber.Ctx, x any) error {
	if err := json.Unmarshal(ctx.Body(), x); err != nil {
		return fiber.ErrBadRequest
	}

	if err := validator.Struct(x); err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return nil
}
