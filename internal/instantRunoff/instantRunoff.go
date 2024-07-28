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

type Tally struct {
	ID         int
	Name       string
	Round      int
	Count      int
	Eliminated bool
	Winner     bool
}

type InstantRunoff struct {
	Rounds  int
	Tallies []*Tally
	Winner  string
}

func (ir *InstantRunoff) ResultsAsString() string {
	talliesByRound := make(map[int][]*Tally)

	for _, tally := range ir.Tallies {
		if v, ok := talliesByRound[tally.Round]; ok {
			talliesByRound[tally.Round] = append(v, tally)
			continue
		}
		talliesByRound[tally.Round] = []*Tally{tally}
	}

	var str strings.Builder

	keys := make([]int, 0, len(talliesByRound))
	for k := range talliesByRound {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, round := range keys {
		tallies := talliesByRound[round]

		str.WriteString("Round ")
		str.WriteString(strconv.Itoa(round))
		str.WriteString("\n\nCurrent first-choice votes:\n")

		var eliminated []string
		for _, tally := range tallies {
			str.WriteString(fmt.Sprintf(" - %s: %d votes\n", tally.Name, tally.Count))

			if tally.Eliminated {
				eliminated = append(eliminated, tally.Name)
			}
		}
		if len(eliminated) > 0 {
			str.WriteString(fmt.Sprintf("\nEliminated %s\n\n---\n\n", strings.Join(eliminated, ", ")))
		}
	}

	str.WriteString(fmt.Sprintf("\nThe winner of this election is %s", ir.Winner))

	return str.String()
}

func Run(votes []*Vote, nameMap map[int]string) (*InstantRunoff, error) {
	var ids []int
	for k := range nameMap {
		ids = append(ids, k)
	}

	logFrequencies := func(freqs map[int]int) []*Tally {
		var pairs [][2]int
		for k, v := range freqs {
			pairs = append(pairs, [2]int{k, v})
		}

		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i][0] < pairs[j][0]
		})

		var tallies []*Tally
		for _, pair := range pairs {
			name := nameMap[pair[0]]
			count := pair[1]
			tallies = append(tallies, &Tally{ID: pair[0], Name: name, Count: count})
		}

		return tallies
	}

	round := 1
	var allTallies []*Tally
	var winner string
	for {
		freqs := getFrequencies(ids, votes, 1)

		tallies := logFrequencies(freqs)

		for _, tally := range tallies {
			tally.Round = round
		}
		allTallies = append(allTallies, tallies...)

		if len(ids) == 1 {
			tallies[0].Winner = true
			winner = tallies[0].Name
			break
		}

		eliminatedID, err := eliminate(ids, votes)
		if err != nil {
			return nil, err
		}

		{
			n := 0
			for _, v := range ids {
				if v != eliminatedID {
					ids[n] = v
					n += 1
				}
			}
			for _, tally := range tallies {
				if tally.ID == eliminatedID {
					tally.Eliminated = true
				}
			}
			ids = ids[:n]
		}

		round += 1
	}

	return &InstantRunoff{Rounds: round, Tallies: allTallies, Winner: winner}, nil
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
