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
	Type            EventType
	Team            models.Team
	StartingPlayer  *models.Player
	FinishingPlayer *models.Player
	EventMeta       EventMeta
}

type EventMeta map[string]interface{}

type EventTrigger func(e Event, s *SimulationState)

//go:generate stringer -type=EventType -output event_type_string.go
type EventType int

const (
	ETNone EventType = iota
	ETHalfTimeExtraTimeAnnouncement
	ETFullTimeExtraTimeAnnouncement
	ETHalfTime
	ETFullTime
	ETSubstitution
	ETPenalty
	ETFreeKickOnGoal
	ETFreeKickDefensiveHalf
	ETFoul
	ETAdvantage
	ETYellowCard
	ETRedCard
	ETPass
	ETGoalScoringChance
	ETInterception
	ETDribble
	ETPossession
	ETSave
	ETGoal
	ETMiss
	ETCross
	ETEndOfFirstHalf
	ETEndOfFirstHalfExtraTime
	ETEndOfSecondHalf
	ETEndOfSecondHalfExtraTime
	ETReset
)

//go:generate stringer -type=Decision -output decision_string.go
type Decision int

const (
	NoDecision Decision = iota
	DecisionLongPass
	DecisionShortPass
	DecisionCross
	DecisionDribble
	DecisionShoot
)

type Outcome struct {
	HomeScore int
	AwayScore int
}

type WeightedEventSet map[EventType]float64

var WeightedGeneralEvents = WeightedEventSet{
	ETSubstitution:          0.02,
	ETFoul:                  0.10,
	ETYellowCard:            0.06,
	ETRedCard:               0.01,
	ETFreeKickDefensiveHalf: 0.06,
}

var WeightedAttackingEvents = WeightedEventSet{
	ETPenalty:        0.01,
	ETFreeKickOnGoal: 0.04,
	ETAdvantage:      0.05,
}

var AllWeightedEvents = mergeWeights(
	WeightedGeneralEvents,
	WeightedAttackingEvents,
)

func RandomWeightedEvent(weights WeightedEventSet, randFloat func() float64) EventType {
	if len(weights) == 0 {
		return ETNone
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

	return ETNone
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
