//go:build !noscrape

package guildScraper

import (
	"fmt"
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

		if err != nil {
			cachedMembershipListLock.Unlock()
			return nil, fmt.Errorf("refresh cached membership list: %w", err)
		}

		cachedMembershipList = members

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
