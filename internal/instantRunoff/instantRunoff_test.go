package instantRunoff

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		votes          []*Vote
		nameMap        map[int]string
		expectedWinner string
		expectedRounds int
	}{
		{
			votes: []*Vote{
				{RankedChoices: []int{1}},
				{RankedChoices: []int{1}},
				{RankedChoices: []int{2}},
				{RankedChoices: []int{2}},
				{RankedChoices: []int{2}},
			},
			nameMap: map[int]string{
				1: "Alice",
				2: "Bob",
				3: "Charlie",
			},
			expectedWinner: "Bob",
			expectedRounds: 3,
		},
		{
			votes: []*Vote{
				{RankedChoices: []int{1, 2, 3}},
				{RankedChoices: []int{1, 2, 3}},
				{RankedChoices: []int{2, 1, 3}},
				{RankedChoices: []int{2, 3, 1}},
				{RankedChoices: []int{3, 2, 1}},
			},
			nameMap: map[int]string{
				1: "Alice",
				2: "Bob",
				3: "Charlie",
			},
			expectedWinner: "Bob",
			expectedRounds: 3,
		},
		{
			votes: []*Vote{
				{RankedChoices: []int{1}},
				{RankedChoices: []int{1, 3}},
				{RankedChoices: []int{1, 2}},
				{RankedChoices: []int{2, 1, 3}},
				{RankedChoices: []int{3, 2, 1}},
				{RankedChoices: []int{3, 2}},
			},
			nameMap: map[int]string{
				1: "Alice",
				2: "Bob",
				3: "Charlie",
			},
			expectedWinner: "Alice",
			expectedRounds: 3,
		},
	}

	for dataset, tc := range testCases {
		t.Run(fmt.Sprintf("dataset #%d", dataset), func(t *testing.T) {
			ir, err := Run(tc.votes, tc.nameMap)

			require.NoError(t, err)
			require.Equal(t, tc.expectedWinner, ir.Winner)
			require.Equal(t, tc.expectedRounds, ir.Rounds)
		})
	}
}

func TestResultsAsString(t *testing.T) {
	tallies := []*Tally{
		{ID: 1, Name: "Alice", Round: 1, Count: 2},
		{ID: 2, Name: "Bob", Round: 1, Count: 2},
		{ID: 3, Name: "Charlie", Round: 1, Count: 1, Eliminated: true},
		{ID: 1, Name: "Alice", Round: 2, Count: 2, Eliminated: true},
		{ID: 2, Name: "Bob", Round: 2, Count: 3},
		{ID: 2, Name: "Bob", Round: 3, Count: 5},
	}

	ir := &InstantRunoff{
		Rounds:  3,
		Tallies: tallies,
		Winner:  "Bob",
	}

	results := ir.ResultsAsString()

	fmt.Printf("%s\n", ir.ResultsAsString())

	expectedResults := "Round 1\n\nCurrent first-choice votes:\n - Alice: 2 votes\n - Bob: 2 votes\n - Charlie: 1 votes\n\nEliminated Charlie\n\n---\n\nRound 2\n\nCurrent first-choice votes:\n - Alice: 2 votes\n - Bob: 3 votes\n\nEliminated Alice\n\n---\n\nRound 3\n\nCurrent first-choice votes:\n - Bob: 5 votes\n\nThe winner of this election is Bob"

	require.Equal(t, expectedResults, results)
}
