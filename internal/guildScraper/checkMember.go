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
)

func GetMember(studentID string) (*GuildMember, error) {
	cachedMembershipListLock.RLock()

	for _, x := range cachedMembershipList {
		if x.ID == studentID {
			cachedMembershipListLock.RUnlock()
			return x, nil
		}
	}

	cachedMembershipListLock.RUnlock()

	cachedMembershipListLock.Lock()
	defer cachedMembershipListLock.Unlock()

	for _, x := range cachedMembershipList {
		if x.ID == studentID {
			return x, nil
		}
	}

	members, err := GetMembersList()
	if err == nil {
		cachedMembershipList = members
	} else {
		slog.Warn("failed to refresh cached membership list", "error", err)
		return nil, err
	}

	for _, x := range cachedMembershipList {
		if x.ID == studentID {
			return x, nil
		}
	}

	return nil, nil
}
