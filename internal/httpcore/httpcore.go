package httpcore

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"git.tdpain.net/codemicro/society-voting/internal/config"
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

	app.Get("/api/elections/current", e.apiGetActiveElectionInformation)
	app.Post("/api/elections/current/vote", e.apiVote)

	app.Post("/api/admin/election", e.apiAdminCreateElection)
	app.Delete("/api/admin/election", e.apiAdminDeleteElection)
	app.Get("/api/admin/election/sse", e.apiAdminRunningElectionSSE)
	app.Post("/api/admin/election/start", e.apiAdminStartElection)
	app.Post("/api/admin/election/stop", e.apiAdminStopElection)

	app.Delete("/api/admin/user/delete", e.apiAdminDeleteUser)

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

const (
	authRegularUser uint = 1 << iota
	authAdminUser
)

func getSessionAuth(ctx *fiber.Ctx, authType uint) (string, bool) {
	cookieVal := ctx.Cookies(sessionTokenCookieName)

	if cookieVal == "" {
		return "", false
	}

	if authType&authRegularUser != 0 {
		decodedToken, tokenOkay := validateData("token", cookieVal)

		if tokenOkay {
			return decodedToken, true
		}
	}

	if authType&authAdminUser != 0 {
		decodedAdmin, adminOkay := validateData("admin", cookieVal)

		if adminOkay {
			return decodedAdmin, true
		}
	}

	return "", false
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
