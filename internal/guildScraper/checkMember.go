//go:build !noscrape

package guildScraper

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

var (
	cachedMembershipList              []*GuildMember
	cachedMembershipListLock          = new(sync.RWMutex)
	cachedMembershipListLastRefreshed time.Time
)

func GetMember(studentID string) (*GuildMember, error) {
	cachedMembershipListLock.RLock()

	if time.Now().Sub(cachedMembershipListLastRefreshed) > time.Minute*5 {
		cachedMembershipListLock.RUnlock()
		cachedMembershipListLock.Lock()

		members, err := GetMembersList()

		if err == nil {
			cachedMembershipList = members
			cachedMembershipListLastRefreshed = time.Now()
		} else if cachedMembershipListLastRefreshed.IsZero() {
			cachedMembershipListLock.Unlock()
			return nil, fmt.Errorf("initially load membership list: %w", err)
		} else {
			slog.Warn("failed to refresh cached membership list", "error", err)
		}

		cachedMembershipListLock.Unlock()
		cachedMembershipListLock.RLock()
	}

	var target *GuildMember

	for _, x := range cachedMembershipList {
		if x.ID == studentID {
			target = x
			break
		}
	}

	cachedMembershipListLock.RUnlock()

	return target, nil
}
