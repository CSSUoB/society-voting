package guildScraper

import (
	"bytes"
	"context"
	"fmt"
	"github.com/CSSUoB/society-voting/internal/config"
	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
	"strings"
	"time"
)

type GuildMember struct {
	ID   string
	Name string
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

	// Rewrite names from `Bloggs, Joe` to `Joe Bloggs`
	for _, x := range res {
		ns := strings.Split(x.Name, ", ")
		x.Name = ns[1] + " " + ns[0]
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

	return pageData, nil
}

func parseGuildMemberPage(pageData string) ([]*GuildMember, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(pageData)))

	if err != nil {
		return nil, fmt.Errorf("create document reader: %w", err)
	}

	var res []*GuildMember

	doc.Find("div.member_list_group").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		header := selection.Find("h3")
		if header == nil {
			err = fmt.Errorf("cannot find member list group title for item %d", i)
			return false
		}

		normalisedHeaderText := strings.ToLower(strings.TrimSpace(header.Text()))

		if normalisedHeaderText == "all members" {
			table := selection.Find("table.msl_table")
			if table == nil {
				err = fmt.Errorf("no table found in group %d", i)
				return false
			}

			res, err = extractMembersFromTable(table)
			if err != nil {
				return false
			}
		}

		return res == nil
	})
	if err != nil {
		return nil, fmt.Errorf("scrape data from guild HTML: %w", err)
	}

	return res, nil
}

func extractMembersFromTable(table *goquery.Selection) ([]*GuildMember, error) {
	var (
		err error
		res []*GuildMember
	)

	table.Find("tr.msl_row,tr.msl_altrow").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		err = nil

		cols := selection.Find("td").Nodes

		if len(cols) != 4 {
			// unexpected number of columns here, what's going on??
			err = fmt.Errorf("unexpected number of columns: expected 4, got %d", len(cols))
			return false
		}

		member := &GuildMember{
			ID:   strings.TrimSpace(goquery.NewDocumentFromNode(cols[1]).Text()),
			Name: strings.TrimSpace(goquery.NewDocumentFromNode(cols[0]).Text()),
		}
		res = append(res, member)

		return true
	})
	if err != nil {
		return nil, fmt.Errorf("scrape guild HTML table: %w", err)
	}
	return res, nil
}
