package simulation

import "github.com/notoriousbfg/football-game/models"

//go:generate stringer -type=Interval -output interval_string.go
type Interval int

const (
	Second Interval = iota
	Minute
	Hour
)

type Event struct {
	Type   EventType
	Team   models.Team
	Player *models.Player
}

func (e *Event) Log(s *SimulationState) {
	// fmt.Printf("%s: %s\n", s.Timestamp(), e.Type)
}

type EventTrigger func(e Event, s *SimulationState)

//go:generate stringer -type=EventType -output event_type_string.go
type EventType int

const (
	EventTypeNone EventType = iota
	HalfTimeExtraTimeAnnouncement
	FullTimeExtraTimeAnnouncement
	HalfTime
	FullTime
	Substitution
	Penalty
	FreeKickOnGoal
	FreeKickDefensiveHalf
	Foul
	Advantage
	YellowCard
	RedCard
	Pass
	GoalScoringChance
)

type Outcome struct {
	HomeScore int
	AwayScore int
}

type WeightedEventSet map[EventType]float64

var WeightedGeneralEvents = WeightedEventSet{
	Substitution:          0.02,
	Foul:                  0.10,
	YellowCard:            0.06,
	RedCard:               0.01,
	FreeKickDefensiveHalf: 0.06,
}

var WeightedAttackingEvents = WeightedEventSet{
	Penalty:        0.01,
	FreeKickOnGoal: 0.04,
	Advantage:      0.05,
}

var AllWeightedEvents = mergeWeights(
	WeightedGeneralEvents,
	WeightedAttackingEvents,
)

func RandomWeightedEvent(weights WeightedEventSet, randFloat func() float64) EventType {
	if len(weights) == 0 {
		return EventTypeNone
	}

	// compute the CDF (Cumulative Distribution Function)
	cdf := make(WeightedEventSet, len(weights))
	total := 0.0
	for i, w := range weights {
		total += w
		cdf[i] = total
	}

	r := randFloat() * total
	for evtType := range cdf {
		if cdf[evtType] >= r {
			return evtType
		}
	}

	return EventTypeNone
}

func mergeWeights(sets ...WeightedEventSet) WeightedEventSet {
	newSet := WeightedEventSet{}
	for _, set := range sets {
		for evtType, weight := range set {
			newSet[evtType] = weight
		}
	}
	return newSet
}
