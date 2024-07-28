//go:build noscrape

package guildScraper

import (
	"github.com/gofiber/fiber/v2/utils"
	"log/slog"
)

func init() {
	slog.Warn("Built with noscrape tag (fictional user information will be used)")
}

var adminID = ""

func GetMember(studentID string) (*GuildMember, error) {
	if adminID == "" {
		adminID = utils.CopyString(studentID) // This magic copy is because we're retaining Fiber variables that get changed to prevent allocations
	}
	isAdmin := adminID == studentID

	return &GuildMember{
		ID:                studentID,
		FirstName:         "Martin",
		LastName:          "Martinson",
		IsCommitteeMember: isAdmin,
	}, nil
}
