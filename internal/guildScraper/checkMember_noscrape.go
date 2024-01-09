//go:build noscrape

package guildScraper

import "log/slog"

func init() {
	slog.Warn("Built with noscrape tag (fictional user information will be used)")
}

var hasAdmined = false

func GetMember(studentID string) (*GuildMember, error) {
	isAdmin := !hasAdmined
	hasAdmined = true
	return &GuildMember{
		ID:                studentID,
		Name:              "Martin Martinson",
		IsCommitteeMember: isAdmin,
	}, nil
}
