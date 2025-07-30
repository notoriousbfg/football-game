package models

import (
	"fmt"
	"slices"
)

func (t *Team) SearchPlayers(options PlayerSearchOptions) Player {
	type PlayerScore struct {
		Player Player
		Score  float64
	}

	similarityWeights := make(map[PlayerPosition]float64)
	visited := make(map[PlayerPosition]bool)
	queue := []struct {
		Pos   PlayerPosition
		Depth int
	}{}

	for _, pos := range options.Positions {
		similarityWeights[pos] = 2.0 // direct match weight
		visited[pos] = true
		queue = append(queue, struct {
			Pos   PlayerPosition
			Depth int
		}{pos, 0})
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Depth >= 3 {
			continue
		}

		for _, similar := range SimilarPositions[current.Pos] {
			if !visited[similar] {
				visited[similar] = true

				weight := map[int]float64{
					1: 1.0,  // 1st-degree similar
					2: 0.5,  // 2nd-degree similar
					3: 0.25, // 3rd-degree similar
				}[current.Depth+1]

				similarityWeights[similar] = weight

				queue = append(queue, struct {
					Pos   PlayerPosition
					Depth int
				}{similar, current.Depth + 1})
			}
		}
	}

	var highest *PlayerScore
	for _, player := range t.Players {
		if _, excluded := options.Exclusions[player.Number]; excluded {
			continue
		}

		score := 0.0

		if weight, ok := similarityWeights[player.Position]; ok {
			score += weight
		}

		if player.Name == options.Name {
			score += 2
		}

		if t.Name == options.Name {
			score += 2
		}

		if player.Number == options.Number {
			score += 2
		}

		if score == 0 {
			continue
		}

		if highest == nil || score > highest.Score {
			highest = &PlayerScore{Player: player, Score: score}
		}
	}

	if highest == nil {
		panic(fmt.Errorf("no player found with options (positions: %+v, exclusions: %+v)", options.Positions, options.Exclusions))
	}

	return highest.Player
}

func (t *Team) RandomPlayerInGroup(team *Team, positions []PlayerPosition, randomFloat func() float64) *Player {
	players := make([]Player, 0)
	for _, player := range team.Players {
		if slices.Contains(positions, player.Position) {
			players = append(players, player)
		}
	}
	randomIndex := int(randomFloat() * float64(len(players)))
	if randomIndex > 0 {
		randomIndex -= 1
	}
	return &players[randomIndex]
}

func (t *Team) ChooseReceiver(passingPlayer Player, underPressure, isLongPass bool, randomFloat func() float64) *Player {
	var receiver *Player
	if underPressure {
		if slices.Contains(Forwards, passingPlayer.Position) {
			receiver = t.RandomPlayerInGroup(t, Midfielders, randomFloat)
		} else {
			receiver = t.RandomPlayerInGroup(t, Defenders, randomFloat)
		}
	} else {
		if slices.Contains(Defenders, passingPlayer.Position) {
			if isLongPass {
				receiver = t.RandomPlayerInGroup(t, Forwards, randomFloat)
			} else {
				receiver = t.RandomPlayerInGroup(t, Midfielders, randomFloat)
			}
		} else {
			receiver = t.RandomPlayerInGroup(t, Forwards, randomFloat)
		}
	}
	return receiver
}
