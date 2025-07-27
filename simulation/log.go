package simulation

import (
	"fmt"
)

type Log struct{}

func (l *Log) logPass(s *SimulationState, e Event) {
	fmt.Printf("(%s) %s passes to %s\n", s.Timestamp(), e.StartingPlayer.Name, e.FinishingPlayer.Name)
}

func (l *Log) logGoal(s *SimulationState, e Event) {
	fmt.Printf("(%s) %s shoots and scores!\n", s.Timestamp(), e.FinishingPlayer.Name)
}

func (l *Log) logCross(s *SimulationState, e Event) {
	fmt.Printf("(%s) %s crosses to %s\n", s.Timestamp(), e.StartingPlayer.Name, e.FinishingPlayer.Name)
}

func (l *Log) logDribble(s *SimulationState, e Event) {
	fmt.Printf("(%s) %s is dribbling with the ball\n", s.Timestamp(), e.StartingPlayer.Name)
}

func (l *Log) logInterception(s *SimulationState, e Event) {
	fmt.Printf("(%s) %s loses the ball to %s\n", s.Timestamp(), e.StartingPlayer.Name, e.FinishingPlayer.Name)
}

func (l *Log) logPossession(s *SimulationState, e Event) {
	fmt.Printf("(%s) %s has the ball\n", s.Timestamp(), e.FinishingPlayer.Name)
}

func (l *Log) logYellowCard(s *SimulationState, e Event) {
	fmt.Printf("(%s) %s is given a yellow card for a foul\n", s.Timestamp(), e.FinishingPlayer.Name)
}

func (l *Log) logRestart(s *SimulationState) {
	fmt.Printf("(%s) The game restarts.\n", s.Timestamp())
}

func (l *Log) logMissingTrigger(s *SimulationState, e Event) {
	fmt.Printf("(%s) No trigger for Event type '%s'.\n", s.Timestamp(), e.Type)
}

func (l *Log) logSave(s *SimulationState, e Event) {
	fmt.Printf("(%s) %s took a shot but it was saved by %s\n", s.Timestamp(), e.StartingPlayer.Name, e.FinishingPlayer.Name)
}
