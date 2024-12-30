package httpcore

import (
	"errors"
	"fmt"
	"sort"

	"github.com/CSSUoB/society-voting/internal/database"
	"github.com/CSSUoB/society-voting/internal/events"
	"github.com/gofiber/fiber/v2"
)

func (endpoints) apiAdminDeleteUser(ctx *fiber.Ctx) error {
	var request = struct {
		UserID string `json:"userID" validate:"ne=0"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	if err := database.DeleteUser(request.UserID); err != nil {
		return fmt.Errorf("apiAdminDeleteUser delete user: %w", err)
	}

	actor := ctx.Locals("userID").(string)
	events.SendEvent(events.TopicUserDeleted, &events.UserDeletedData{
		ActingUserID: actor,
		UserID:       request.UserID,
	})

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (endpoints) apiAdminListUsers(ctx *fiber.Ctx) error {
	users, err := database.GetAllUsers()
	if err != nil {
		return fmt.Errorf("apiAdminListUsers get all users: %w", err)
	}

	// place restricted users at the top of the list
	sort.Slice(users, func(i, j int) bool {
		if users[i].IsRestricted != users[j].IsRestricted {
			return users[i].IsRestricted
		}
		return users[i].Name < users[j].Name
	})

	return ctx.JSON(users)
}

func (endpoints) apiAdminToggleRestrictUser(ctx *fiber.Ctx) error {
	var request = struct {
		UserID string `json:"userID" validate:"ne=0"`
	}{}

	if err := parseAndValidateRequestBody(ctx, &request); err != nil {
		return err
	}

	tx, err := database.GetTx()
	if err != nil {
		return fmt.Errorf("apiAdminToggleRestrictUser start tx: %w", err)
	}
	defer tx.Rollback()

	user, err := database.GetUser(request.UserID, tx)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return &fiber.Error{
				Code:    fiber.StatusNotFound,
				Message: "User with that ID not found",
			}
		}
		return fmt.Errorf("apiAdminToggleRestrictUser get user: %w", err)
	}

	if user.IsAdmin {
		return &fiber.Error{
			Code:    fiber.StatusForbidden,
			Message: "This user is an administrator - administrators cannot be restricted.",
		}
	}

	user.IsRestricted = !user.IsRestricted

	if err := user.Update(tx); err != nil {
		return fmt.Errorf("apiAdminToggleRestrictUser update user: %w", err)
	}

	if user.IsRestricted {
		if err := database.DeleteAllCandidatesForUser(user.StudentID, tx); err != nil {
			return fmt.Errorf("apiAdminToggleRestrictUser delete user candidacies: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("apiAdminToggleRestrictUser commit tx: %w", err)
	}

	actor := ctx.Locals("userID").(string)
	events.SendEvent(events.TopicUserRestricted, &events.UserRestrictedData{
		ActingUserID: actor,
		User:         user,
	})

	return ctx.JSON(map[string]bool{"isRestricted": user.IsRestricted})
}
