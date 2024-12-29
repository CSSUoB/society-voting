package httpcore

import (
	"context"
	cryptoRand "crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/httpcore/htmlutil"
	"github.com/CSSUoB/society-voting/web"
	goalone "github.com/bwmarrin/go-alone"
	validate "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

type endpoints struct{}

type middleware struct{}

const loginActionEndpoint = "/auth/login/do"

func ListenAndServe(ctx context.Context, addr string) error {
	if signer == nil {
		return errors.New("signer not initialised")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			var (
				code = fiber.StatusInternalServerError
				e    *fiber.Error
				msg  string
			)
			if errors.As(err, &e) {
				code = e.Code
				msg = err.Error()
			} else {
				slog.Error("fiber runtime error", "error", err, "URL", ctx.OriginalURL())
				msg = "Internal Server Error"
			}
			ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return ctx.Status(code).SendString(msg)
		},
		DisableStartupMessage: !config.Get().Debug,
		AppName:               "society-voting",
	})

	e := endpoints{}
	m := middleware{}

	app.Use(func(ctx *fiber.Ctx) error {
		authToken := ctx.Cookies(sessionTokenCookieName)
		authSet := authToken != ""

		if authSet {
			decodedToken, isValid := validateData(authToken)
			if !isValid {
				authSet = false
			} else {
				ctx.Locals("token", decodedToken)
			}
		}
		ctx.Locals("isauthset", authSet)

		return ctx.Next()
	})

	app.Get("/auth/login", e.authLoginPage)
	app.Post("/auth/login", e.authLoginPage)
	app.Get(loginActionEndpoint, e.authLogin)
	app.Post(loginActionEndpoint, e.authLogin)

	app.Get("/auth/logout", e.authLogout)

	apiGroup := app.Group("/api", m.requireAuthenticated)
	apiGroup.Get("/me", e.apiMe)
	apiGroup.Put("/me/name", e.apiSetOwnName)

	//todo fix
	//apiGroup.Get("/election/sse", e.apiElectionsSSE)
	apiGroup.Get("/poll", e.apiListPolls)
	apiGroup.Get("/poll/current", e.apiGetActivePollInformation)
	apiGroup.Get("/poll/results", e.apiGetPollOutcome)
	apiGroup.Post("/election/stand", e.apiStandForElection)
	apiGroup.Delete("/election/stand", e.apiWithdrawFromElection)
	apiGroup.Post("/election/vote", e.apiVoteInElection)
	apiGroup.Post("/referendum/vote", e.apiVoteInReferendum)

	adminGroup := apiGroup.Group("/admin", m.requireAdmin)
	adminGroup.Post("/election", e.apiAdminCreateElection)
	adminGroup.Post("/election/start", e.apiAdminStartElection)
	adminGroup.Post("/election/stop", e.apiAdminStopElection)
	adminGroup.Post("/referendum", e.apiAdminCreateReferendum)
	adminGroup.Post("/referendum/start", e.apiAdminStartReferendum)
	adminGroup.Post("/referendum/stop", e.apiAdminStopReferendum)
	adminGroup.Delete("/poll", e.apiAdminDeletePoll)
	adminGroup.Post("/poll/publish", e.apiAdminPublishPollOutcome)
	//todo fix
	//adminGroup.Get("/election/sse", e.apiAdminRunningElectionSSE)
	adminGroup.Get("/user", e.apiAdminListUsers)
	adminGroup.Post("/user/restrict", e.apiAdminToggleRestrictUser)
	adminGroup.Delete("/user/delete", e.apiAdminDeleteUser)

	app.Get("/presenter", m.requireAdmin, e.presenterPage)
	app.Get("/presenter/qr", m.requireAdmin, e.presenterQRCode)

	app.Use("/", func(ctx *fiber.Ctx) error {
		if p := ctx.Path(); p != "/" && !urlFileRegexp.MatchString(p) {
			// override path to serve everything as if it were requested to /
			ctx.Path("/")
			return ctx.RestartRouting()
		}

		return ctx.Next()
	}, web.GetHandler())

	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusNotFound)
		return htmlutil.SendPage(
			ctx,
			htmlutil.SkeletonPage("404 Not Found",
				html.H1(g.Text("404 Not Found")),
				html.P(
					g.Text("Oh, that's weird. We can't find what you're looking for, but chances are "),
					html.A(g.Text("you're probably looking for this instead!"), g.Attr("href", "/")),
				),
			),
		)
	})

	slog.Info("HTTP server alive", "address", addr, "voteCode", voteCode)

	serverExitValChan := make(chan error)
	go func() {
		serverExitValChan <- app.Listen(addr)
	}()

	select {
	case <-ctx.Done():
		slog.Info("Shutting down HTTP server...")
		err := app.ShutdownWithTimeout(time.Second * 5)
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Warn("5 second deadline exceeded, forcibly shutting down")
			return nil
		}
		return err
	case v := <-serverExitValChan:
		return v
	}
}

var urlFileRegexp = regexp.MustCompile(`[\w\-/]+\.[a-zA-Z]+$`)

var (
	signer    *goalone.Sword
	validator = validate.New(validate.WithRequiredStructEnabled())
	voteCode  string
)

func init() {
	for i := 0; i < 4; i += 1 {
		voteCode += string(rune('A' + rand.Intn(26)))
	}
}

func InitialiseSigner(signingKey string) {
	rawKeyBytes := []byte(signingKey)
	keyBytes := make([]byte, 512)

	if len(rawKeyBytes) == 0 {
		goto generate
	}

	if _, err := hex.Decode(keyBytes, rawKeyBytes); err != nil {
		goto generate
	}

	goto commit

generate:
	slog.Info("using randomly generated session signing key")
	keyBytes = make([]byte, 512)
	if !config.Get().Debug { // This is so that the access tokens doesn't change from run-to-run for ease of testing
		if _, err := cryptoRand.Read(keyBytes); err != nil {
			slog.Error("unable to generate random secret for token signing", "error", err)
			os.Exit(1)
		}
	}

commit:
	signer = goalone.New(keyBytes)
}

var sessionTokenCookieName = "vot-tok"

func newSessionTokenCookie(tok string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:        sessionTokenCookieName,
		Value:       tok,
		Secure:      !config.Get().Debug,
		SessionOnly: true,
	}
}

func newSessionTokenDeletionCookie() *fiber.Cookie {
	return &fiber.Cookie{
		Name:        sessionTokenCookieName,
		Secure:      !config.Get().Debug,
		SessionOnly: true,
		Expires:     time.Date(2010, 1, 1, 1, 1, 1, 1, time.UTC),
	}
}

func signData(data string) string {
	return hex.EncodeToString(signer.Sign([]byte(data)))
}

func validateData(token string) (string, bool) {
	dat, err := hex.DecodeString(token)
	if err != nil {
		return "", false
	}
	originalData, err := signer.Unsign(dat)
	if err != nil {
		return "", false
	}
	return string(originalData), true
}

type authType uint

const (
	authNotAuthed   authType = 0
	authRegularUser authType = 1 << iota
	authAdminUser
)

func getSessionAuth(ctx *fiber.Ctx) (string, authType) {
	authSet := ctx.Locals("isauthset").(bool)
	if !authSet {
		return "", authNotAuthed
	}

	decodedToken := ctx.Locals("token").(string)

	var isAdmin bool
	if err := database.Get().NewSelect().Table("users").Column("is_admin").Where("id = ?", decodedToken).Scan(context.Background(), &isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", authNotAuthed
		}
		// if we can't handle sessions properly we may as well not run
		panic(err)
	}

	if isAdmin {
		return decodedToken, authRegularUser | authAdminUser
	}
	return decodedToken, authRegularUser
}

func parseAndValidateRequestBody(ctx *fiber.Ctx, x any) error {
	if err := json.Unmarshal(ctx.Body(), x); err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request JSON (" + err.Error() + ")",
		}
	}

	if err := validator.Struct(x); err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return nil
}
