package httpcore

import (
	"errors"
	"fmt"

	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/gofiber/fiber/v2"
)

func (middleware) requireAuthenticated(ctx *fiber.Ctx) error {
	userID, authStatus := getSessionAuth(ctx)
	if authStatus == authNotAuthed {
		return fiber.ErrUnauthorized
	}

	ctx.Locals("userID", userID)
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

func (middleware) validateUserExists(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	user, err := database.GetUser(userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			// User has been deleted
			ctx.Cookie(newSessionTokenDeletionCookie())
			return fiber.ErrUnauthorized
		}
		return fmt.Errorf("validateUserExists: %w", err)
	}

	ctx.Locals("user", user)

	return ctx.Next()
}
