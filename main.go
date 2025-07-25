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
}
