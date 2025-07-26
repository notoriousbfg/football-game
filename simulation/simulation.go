package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/notoriousbfg/football-game/models"
)

type Simulation struct {
	State             *SimulationState
	Match             Match
	RandomFloat       func() float64
	SynergyMultiplier int
	TacticalCounters  map[int]TacticalCounter
	EventTriggers     map[EventType]EventTrigger
	Pitch             *Pitch
}

func (sim *Simulation) Run() {
	sim.Pitch = NewPitch(&sim.Match)
	sim.Pitch.Draw()

	halfSeconds := 45 * 60

	// first half
	sim.runPeriod(0, halfSeconds)
	sim.runPeriod(halfSeconds, sim.State.FirstHalfExtraTime)

	// second half
	sim.runPeriod(halfSeconds, halfSeconds)
	sim.runPeriod(halfSeconds, sim.State.SecondHalfExtraTime)

	sim.State.Outcome = &Outcome{
		HomeScore: sim.State.HomeScore,
		AwayScore: sim.State.AwayScore,
	}
}

func (sim *Simulation) runPeriod(start, seconds int) {
	for i := start; i <= start+seconds; i++ {
		sim.State.Time = sim.State.Time.Add(time.Second)
		// if evtType := RandomWeightedEvent(AllWeightedEvents, sim.RandomFloat); evtType != EventTypeNone {
		// 	if trigger, exists := sim.EventTriggers[evtType]; exists {
		// 		trigger(
		// 			Event{Type: evtType, Source: sim.eventSource(evtType)},
		// 			sim.State,
		// 		)
		// 	}
		// }

	}
}

type SimulationState struct {
	Start               time.Time
	Time                time.Time
	HomeScore           int
	AwayScore           int
	HomeYellowCards     int
	AwayYellowCards     int
	HomeRedCards        int
	AwayRedCards        int
	HomeMomentum        float64
	AwayMomentum        float64
	HomeTeamAttacking   bool
	AwayTeamAttacking   bool
	Stalemate           bool
	FirstHalfExtraTime  int // seconds
	SecondHalfExtraTime int // seconds
	Events              []Event
	Outcome             *Outcome
}

func (s *SimulationState) CaptureEvent(e Event) {
	s.Events = append(s.Events, e)

	switch e.Type {
	case YellowCard:
		// TODO is home or away?
		s.HomeYellowCards++
	}
}

func (s *SimulationState) Timestamp() string {
	duration := s.Time.Sub(s.Start)
	return fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
}

type Match struct {
	H models.Team
	A models.Team
}

type TacticalCounter struct{}

func goalChanceEvent() EventType {
	return FreeKickOnGoal
}

func CreateSimulation(home, away models.Team) *Simulation {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFloat := func() float64 { return randGen.Float64() }

	startingPlayer := home.SearchPlayers(models.PlayerSearchOptions{
		Position: models.Striker | models.CentralMidfielder,
	})

	state := &SimulationState{
		HomeScore:         0,
		AwayScore:         0,
		HomeYellowCards:   0,
		AwayYellowCards:   0,
		HomeRedCards:      0,
		AwayRedCards:      0,
		HomeMomentum:      1.0,
		AwayMomentum:      1.0,
		HomeTeamAttacking: false,
		AwayTeamAttacking: false,
		Stalemate:         false,
		Events: []Event{
			{
				Type:   Pass,
				Team:   home,
				Player: &startingPlayer,
			},
		},
	}

	triggers := make(map[EventType]EventTrigger)

	triggers[YellowCard] = func(e Event, s *SimulationState) {}

	// triggers[HomeTeamAttacking] = func(e Event, s *SimulationState) {
	// 	if s.evaluateAttack(e, randomFloat) {
	// 		s.HomeScore++
	// 		s.HomeMomentum += 0.1
	// 		s.AwayMomentum -= 0.1
	// 	}
	// 	e.Log(s)
	// }

	// triggers[AwayTeamAttacking] = func(e Event, s *SimulationState) {
	// 	if s.evaluateAttack(e, randomFloat) {
	// 		s.AwayScore++
	// 		s.AwayMomentum += 0.1
	// 		s.HomeMomentum -= 0.1
	// 	}
	// 	e.Log(s)
	// }

	return &Simulation{
		State:             state,
		Match:             Match{H: home, A: away},
		RandomFloat:       randomFloat,
		SynergyMultiplier: 1,
		TacticalCounters:  make(map[int]TacticalCounter),
		EventTriggers:     triggers,
	}
}

func (s *SimulationState) evaluateAttack(event Event, randFloat func() float64) bool {
	// get last event

	// get team/player from event source

	// pass, shoot, dribble?

	// evaluate likelihood of chance

	// evaluate likelihood of goal

	// weighted dice roll

	return false
}
