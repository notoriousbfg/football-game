package main

import (
	"fmt"

	"github.com/notoriousbfg/football-game/scenarios"
	"github.com/notoriousbfg/football-game/simulation"
)

func main() {
	sim := simulation.CreateSimulation(
		scenarios.HomeTeam(),
		scenarios.AwayTeam(),
	)

	sim.Run()

	fmt.Print(sim.State.Outcome)

	// exclusions := make(map[string]map[models.PlayerNumber]string)
	// team := scenarios.HomeTeam()
	// positions := []models.PlayerPosition{
	// 	models.LeftBack,
	// 	models.LeftCentreBack,
	// 	models.RightCentreBack,
	// 	models.RightBack,
	// }
	// for _, p := range positions {
	// 	player := team.SearchPlayers(models.PlayerSearchOptions{
	// 		Position:   p,
	// 		Exclusions: exclusions[team.Name],
	// 	})
	// 	if exclusions[team.Name] == nil {
	// 		exclusions[team.Name] = make(map[models.PlayerNumber]string)
	// 	}
	// 	exclusions[team.Name][player.Number] = player.Initials()
	// }
	// fmt.Print(exclusions)
}
