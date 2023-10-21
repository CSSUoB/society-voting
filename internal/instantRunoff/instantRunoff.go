package instantRunoff

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type Vote struct {
	RankedChoices []int
}

func Run(votes []*Vote, nameMap map[int]string) string {
	var ids []int
	for k := range nameMap {
		ids = append(ids, k)
	}

	var log strings.Builder

	logFrequencies := func(freqs map[int]int) {
		var pairs [][2]int
		for k, v := range freqs {
			pairs = append(pairs, [2]int{k, v})
		}

		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i][0] < pairs[j][0]
		})

		var strs []string
		for _, pair := range pairs {
			strs = append(strs, fmt.Sprintf(" - %s: %d votes", nameMap[pair[0]], pair[1]))
		}

		log.WriteString(strings.Join(strs, "\n"))
	}

	round := 1
	for {
		log.WriteString("================================================================\n")
		log.WriteString("ROUND ")
		log.WriteString(strconv.Itoa(round))
		log.WriteRune('\n')

		freqs := getFrequencies(ids, votes, 1)

		log.WriteString("\nCurrent first-choice votes:\n")
		logFrequencies(freqs)
		log.WriteString("\n\n")

		if len(ids) == 1 {
			log.WriteString(nameMap[ids[0]])
			log.WriteString(" has won\n")
			break
		}

		eliminatedID, err := eliminate(ids, votes)
		if err != nil {
			log.WriteString("Error: ")
			log.WriteString(err.Error())
			log.WriteString("\n")
			break
		}

		{
			n := 0
			for _, v := range ids {
				if v != eliminatedID {
					ids[n] = v
					n += 1
				}
			}
			ids = ids[:n]
		}

		log.WriteString("Eliminated ")
		log.WriteString(nameMap[eliminatedID])
		log.WriteRune('\n')

		round += 1
	}

	return log.String()
}

func getFrequencies(ids []int, votes []*Vote, choiceNumber int) map[int]int {
	res := make(map[int]int)
	for _, id := range ids {
		res[id] = 0
	}

	for _, vote := range votes {
		if len(vote.RankedChoices) < choiceNumber {
			continue
		}
		res[vote.RankedChoices[choiceNumber-1]] = res[vote.RankedChoices[choiceNumber-1]] + 1
	}

	return res
}

func eliminate(ids []int, votes []*Vote) (int, error) {
	freqs := getFrequencies(ids, votes, 1)
	minVotes := math.MaxInt
	for _, v := range freqs {
		if v < minVotes {
			minVotes = v
		}
	}

	var lowestScoringCandidates []int
	for k, v := range freqs {
		if v == minVotes {
			lowestScoringCandidates = append(lowestScoringCandidates, k)
		}
	}

	var toEliminate int

	if len(lowestScoringCandidates) == 1 {
		toEliminate = lowestScoringCandidates[0]
	} else {
		toEliminate = lowestScoringCandidates[rand.Intn(len(lowestScoringCandidates))]
	}

	for _, vote := range votes {
		n := 0
		for _, val := range vote.RankedChoices {
			if val != toEliminate {
				vote.RankedChoices[n] = val
				n += 1
			}
		}
		vote.RankedChoices = vote.RankedChoices[:n]
	}

	return toEliminate, nil
}
