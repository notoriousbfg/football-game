package main

import (
	"fmt"

	"github.com/notoriousbfg/football-game/simulation"
)

func main() {
	sim := simulation.CreateSimulation(
		simulation.HomeTeam(),
		simulation.AwayTeam(),
	)

	sim.Run()

	fmt.Print(sim.State.Outcome)
}
