//go:build noscrape

package guildScraper

import "log/slog"

func init() {
	slog.Warn("Built with noscrape tag (fictional user information will be used)")
}

func GetMember(studentID string) (*GuildMember, error) {
	return &GuildMember{
		ID:   studentID,
		Name: "Martin Martinson",
	}, nil
}
