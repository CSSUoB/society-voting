package guildScraper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
	"strings"
	"time"
)

type GuildMember struct {
	ID                string
	FirstName         string
	LastName          string
	IsCommitteeMember bool
}

func GetMembersList() ([]*GuildMember, error) {
	pageData, err := fetchMembershipPage()
	if err != nil {
		return nil, fmt.Errorf("get guild membership list: %w", err)
	}

	res, err := parseGuildMemberPage(pageData)
	if err != nil {
		return nil, fmt.Errorf("parse guild membership list: %w", err)
	}

	return res, nil
}

func fetchMembershipPage() (string, error) {
	conf := config.Get().Guild

	var pageData string
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := requests.URL("https://www.guildofstudents.com").
		Pathf("/organisation/memberlist/%s", conf.SocietyID).
		Param("sort", "groups").
		Headers(map[string][]string{
			"Cache-Control": {"no-cache"},
			"Pragma":        {"no-cache"},
			"Expires":       {"0"},
			"User-Agent":    {"Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/119.0"},
		}).
		Cookie(".ASPXAUTH", conf.SessionToken).
		ToString(&pageData).
		Fetch(ctx)

	if err != nil {
		return "", err
	}

	// MSL does not return proper error codes so we have to resort to this instead
	if strings.Contains(pageData, "Student Login") || strings.Contains(pageData, "You do not have permission to see this organisation's members.") {
		return "", errors.New("access token expired or invalid")
	}

	return pageData, nil
}

func parseGuildMemberPage(pageData string) ([]*GuildMember, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(pageData)))

	if err != nil {
		return nil, fmt.Errorf("create document reader: %w", err)
	}

	var (
		// both maps are in the format <id>:<name>
		members    map[string]string
		cmtMembers map[string]string
	)

	doc.Find("div.member_list_group").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		header := selection.Find("h3")
		if header == nil {
			err = fmt.Errorf("cannot find member list group title for item %d", i)
			return false
		}

		normalisedHeaderText := strings.ToLower(strings.TrimSpace(header.Text()))

		const (
			standardMembers = "standard membership"
			allCmtMembers   = "all committee members"
		)

		if normalisedHeaderText == standardMembers || normalisedHeaderText == allCmtMembers {
			table := selection.Find("table.msl_table")
			if table == nil {
				err = fmt.Errorf("no table found in group %d", i)
				return false
			}

			var parsedMembers map[string]string
			parsedMembers, err = extractMembersFromTable(table)
			if err != nil {
				return false
			}

			if normalisedHeaderText == standardMembers {
				members = parsedMembers
			} else {
				cmtMembers = parsedMembers
			}

		}

		return members == nil || cmtMembers == nil
	})
	if err != nil {
		return nil, fmt.Errorf("scrape data from guild HTML: %w", err)
	}

	var res []*GuildMember

	for id, name := range members {
		// Last name comes first
		ns := strings.Split(name, ", ")

		_, isCmt := cmtMembers[id]
		res = append(res, &GuildMember{
			ID:                strings.TrimSpace(id),
			FirstName:         strings.TrimSpace(ns[1]),
			LastName:          strings.TrimSpace(ns[0]),
			IsCommitteeMember: isCmt,
		})
	}

	return res, nil
}

func extractMembersFromTable(table *goquery.Selection) (map[string]string, error) {
	var (
		err error
		res = make(map[string]string)
	)

	table.Find("tr.msl_row,tr.msl_altrow").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		err = nil

		cols := selection.Find("td").Nodes

		if len(cols) != 4 {
			// unexpected number of columns here, what's going on??
			err = fmt.Errorf("unexpected number of columns: expected 4, got %d", len(cols))
			return false
		}

		id := strings.TrimSpace(goquery.NewDocumentFromNode(cols[1]).Text())
		name := strings.TrimSpace(goquery.NewDocumentFromNode(cols[0]).Text())
		res[id] = name

		return true
	})
	if err != nil {
		return nil, fmt.Errorf("scrape guild HTML table: %w", err)
	}
	return res, nil
}
