package httpcore

import (
	"github.com/gofiber/fiber/v2"
)

func (middleware) requireAuthenticated(ctx *fiber.Ctx) error {
	authToken := ctx.Cookies(sessionTokenCookieName)
	user, authStatus := getAuthStatus(authToken)

	if authStatus == authNotAuthed {
		return fiber.ErrUnauthorized
	} else if authStatus&authInvalid != 0 {
		// Token signature invalid or user has been deleted
		ctx.Cookie(newSessionTokenDeletionCookie())
		return fiber.ErrUnauthorized
	}

	ctx.Locals("user", user)
	ctx.Locals("userID", user.StudentID)
	ctx.Locals("authStatus", authStatus)

	return ctx.Next()
}

func (middleware) requireAdmin(ctx *fiber.Ctx) error {
	authStatus := ctx.Locals("authStatus").(authType)
	if authStatus&authAdminUser == 0 {
		return fiber.ErrForbidden
	}

	return ctx.Next()
}

func (middleware) requireNotRestricted(ctx *fiber.Ctx) error {
	authStatus := ctx.Locals("authStatus").(authType)
	if authStatus&authRestricted == authRestricted {
		return &fiber.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can't do that because you're restricted - please speak to a member of committee.",
		}
	}

	return ctx.Next()
}
