package httpcore

import "github.com/gofiber/fiber/v2"

func (middleware) requireAuthenticated(ctx *fiber.Ctx) error {
	if _, status := getSessionAuth(ctx); status == authNotAuthed {
		return fiber.ErrUnauthorized
	}

	return ctx.Next()
}

func (middleware) requireAdmin(ctx *fiber.Ctx) error {
	if _, status := getSessionAuth(ctx); status&authAdminUser == 0 {
		return fiber.ErrForbidden
	}

	return ctx.Next()
}
