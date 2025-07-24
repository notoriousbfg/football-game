package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/notoriousbfg/football-game/models"
)

type Simulation struct {
	State             *SimulationState
	Comparison        Comparison
	RandomFloat       func() float64
	SynergyMultiplier int
	TacticalCounters  map[int]TacticalCounter
	EventTriggers     map[EventType]EventTrigger
}

func (sim *Simulation) Run() {
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
		if evtType, ok := randomMatchEvent(sim.RandomFloat); ok {
			if trigger, exists := sim.EventTriggers[evtType]; exists {
				trigger(Event{Type: evtType, Source: sim.Comparison}, sim.State)
			}
		}
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

func (s *SimulationState) Timestamp() string {
	duration := s.Time.Sub(s.Start)
	return fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
}

type Comparison struct {
	H models.Team
	A models.Team
}

type TacticalCounter struct{}

//go:generate stringer -type=Interval -output interval_string.go
type Interval int

const (
	Second Interval = iota
	Minute
	Hour
)

type Event struct {
	Count  int
	Type   EventType
	Source Comparison
}

func (e *Event) Log(s *SimulationState) {
	fmt.Printf("%s: %s\n", s.Timestamp(), e.Type)
}

type EventTrigger func(e Event, s *SimulationState)

//go:generate stringer -type=EventType -output event_type_string.go
type EventType int

const (
	HalfTimeExtraTimeAnnouncement EventType = iota
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
	HomeTeamAttacking
	AwayTeamAttacking
)

type Outcome struct {
	HomeScore int
	AwayScore int
}

func randomMatchEvent(randFloat func() float64) (EventType, bool) {
	eventProbabilities := map[EventType]float64{
		Substitution:          0.02,
		Penalty:               0.01,
		FreeKickOnGoal:        0.04,
		FreeKickDefensiveHalf: 0.06,
		Foul:                  0.10,
		Advantage:             0.05,
		YellowCard:            0.06,
		RedCard:               0.01,
		HomeTeamAttacking:     0.05,
		AwayTeamAttacking:     0.05,
	}

	totalWeight := 0.0
	for _, weight := range eventProbabilities {
		totalWeight += weight
	}

	if randFloat() > totalWeight {
		return 0, false
	}

	r := randFloat() * totalWeight
	acc := 0.0
	for event, weight := range eventProbabilities {
		acc += weight
		if r <= acc {
			return event, true
		}
	}

	return 0, false
}

func goalChanceEvent() EventType {
	return FreeKickOnGoal
}

func CreateSimulation(home, away models.Team) *Simulation {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFloat := func() float64 { return randGen.Float64() }

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
	}

	triggers := make(map[EventType]EventTrigger)

	triggers[HomeTeamAttacking] = func(e Event, s *SimulationState) {
		comparison := e.Source
		if evaluateAttack(comparison.H, comparison.A, randomFloat) {
			s.HomeScore++
			s.HomeMomentum += 0.1
			s.AwayMomentum -= 0.1
		}
		e.Log(s)
	}

	triggers[AwayTeamAttacking] = func(e Event, s *SimulationState) {
		comparison := e.Source
		if evaluateAttack(comparison.A, comparison.H, randomFloat) {
			s.AwayScore++
			s.AwayMomentum += 0.1
			s.HomeMomentum -= 0.1
		}
		e.Log(s)
	}

	return &Simulation{
		State:             state,
		Comparison:        Comparison{H: home, A: away},
		RandomFloat:       randomFloat,
		SynergyMultiplier: 1,
		TacticalCounters:  make(map[int]TacticalCounter),
		EventTriggers:     triggers,
	}
}

func evaluateAttack(attacking models.Team, defending models.Team, randFloat func() float64) bool {
	// rewrite me not using AI
	return true
}

func HomeTeam() models.Team {
	return models.Team{
		Name:      "Egg Fried Reus",
		Morale:    88,
		Fitness:   82,
		Chemistry: 91,
		Strategy: models.Strategy{
			Tactic:    models.TacticPressing,
			Formation: models.FormationFourThreeThree,
			PlayerInstructions: map[models.PlayerNumber]models.Instruction{
				1:  {Position: models.PositionCenter}, // Goalkeeper
				2:  {Position: models.PositionWing},   // Right Back
				3:  {Position: models.PositionCenter}, // Center Back
				4:  {Position: models.PositionCenter}, // Center Back
				5:  {Position: models.PositionWing},   // Left Back
				6:  {Position: models.PositionCenter}, // Defensive Mid
				7:  {Position: models.PositionWing},   // Right Wing
				8:  {Position: models.PositionCenter}, // Center Mid
				9:  {Position: models.PositionCenter}, // Striker
				10: {Position: models.PositionCenter}, // Attacking Mid
				11: {Position: models.PositionWing},   // Left Wing
			},
			PlayStyle: models.PlayStyleCreative,
		},
		Training: models.Training{
			Focus: models.Shooting,
		},
		Players: []models.Player{
			{
				Number: 1, // Goalkeeper
				Form:   80, Adaptability: 75, Composure: 85,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 50, Acceleration: 55},
					Passing:  models.PassingSkill{ShortPass: 60, LongPass: 65, Cross: 40, Lob: 50, ThroughBall: 45, Chip: 40},
					Shooting: models.ShootingSkill{Power: 30, Curve: 20, Finish: 20, Spin: 25},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 90,
					Vision:      models.TacticalVision{Passing: 65, Shooting: 30, Defence: 85},
				},
				Stamina: models.Stamina{Stamina: 65},
				Fitness: models.Fitness{Strength: 70, Agility: 80, InjuryTolerance: 85, InjuryResistance: 78},
			},
			{
				Number: 2, // Right Back
				Form:   75, Adaptability: 70, Composure: 68,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 82, Acceleration: 85},
					Passing:  models.PassingSkill{ShortPass: 70, LongPass: 68, Cross: 75, Lob: 60, ThroughBall: 65, Chip: 55},
					Shooting: models.ShootingSkill{Power: 50, Curve: 48, Finish: 45, Spin: 52},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 72,
					Vision:      models.TacticalVision{Passing: 70, Shooting: 40, Defence: 80},
				},
				Stamina: models.Stamina{Stamina: 78},
				Fitness: models.Fitness{Strength: 70, Agility: 82, InjuryTolerance: 75, InjuryResistance: 80},
			},
			{
				Number: 3, // Center Back
				Form:   78, Adaptability: 65, Composure: 80,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 70, Acceleration: 65},
					Passing:  models.PassingSkill{ShortPass: 68, LongPass: 70, Cross: 40, Lob: 65, ThroughBall: 50, Chip: 45},
					Shooting: models.ShootingSkill{Power: 55, Curve: 40, Finish: 35, Spin: 45},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 85,
					Vision:      models.TacticalVision{Passing: 68, Shooting: 30, Defence: 90},
				},
				Stamina: models.Stamina{Stamina: 72},
				Fitness: models.Fitness{Strength: 85, Agility: 60, InjuryTolerance: 80, InjuryResistance: 85},
			},
			{
				Number: 4, // Center Back
				Form:   82, Adaptability: 67, Composure: 78,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 72, Acceleration: 68},
					Passing:  models.PassingSkill{ShortPass: 65, LongPass: 66, Cross: 42, Lob: 60, ThroughBall: 55, Chip: 50},
					Shooting: models.ShootingSkill{Power: 52, Curve: 38, Finish: 30, Spin: 40},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 83,
					Vision:      models.TacticalVision{Passing: 70, Shooting: 35, Defence: 88},
				},
				Stamina: models.Stamina{Stamina: 74},
				Fitness: models.Fitness{Strength: 82, Agility: 62, InjuryTolerance: 80, InjuryResistance: 82},
			},
			{
				Number: 5, // Left Back
				Form:   74, Adaptability: 72, Composure: 69,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 84, Acceleration: 87},
					Passing:  models.PassingSkill{ShortPass: 70, LongPass: 65, Cross: 78, Lob: 62, ThroughBall: 67, Chip: 58},
					Shooting: models.ShootingSkill{Power: 50, Curve: 52, Finish: 48, Spin: 55},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 76,
					Vision:      models.TacticalVision{Passing: 72, Shooting: 45, Defence: 78},
				},
				Stamina: models.Stamina{Stamina: 80},
				Fitness: models.Fitness{Strength: 72, Agility: 84, InjuryTolerance: 75, InjuryResistance: 77},
			},
			{
				Number: 6, // Defensive Midfielder
				Form:   77, Adaptability: 70, Composure: 75,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 75, Acceleration: 78},
					Passing:  models.PassingSkill{ShortPass: 78, LongPass: 72, Cross: 65, Lob: 70, ThroughBall: 75, Chip: 60},
					Shooting: models.ShootingSkill{Power: 62, Curve: 58, Finish: 50, Spin: 55},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 80,
					Vision:      models.TacticalVision{Passing: 80, Shooting: 55, Defence: 82},
				},
				Stamina: models.Stamina{Stamina: 85},
				Fitness: models.Fitness{Strength: 78, Agility: 75, InjuryTolerance: 80, InjuryResistance: 80},
			},
			{
				Number: 7, // Right Wing
				Form:   85, Adaptability: 80, Composure: 82,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 90, Acceleration: 95},
					Passing:  models.PassingSkill{ShortPass: 80, LongPass: 70, Cross: 85, Lob: 68, ThroughBall: 75, Chip: 70},
					Shooting: models.ShootingSkill{Power: 75, Curve: 80, Finish: 78, Spin: 82},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 85,
					Vision:      models.TacticalVision{Passing: 85, Shooting: 80, Defence: 60},
				},
				Stamina: models.Stamina{Stamina: 88},
				Fitness: models.Fitness{Strength: 70, Agility: 90, InjuryTolerance: 75, InjuryResistance: 78},
			},
			{
				Number: 8, // Center Midfielder
				Form:   80, Adaptability: 76, Composure: 77,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 78, Acceleration: 80},
					Passing:  models.PassingSkill{ShortPass: 85, LongPass: 80, Cross: 70, Lob: 75, ThroughBall: 82, Chip: 65},
					Shooting: models.ShootingSkill{Power: 68, Curve: 70, Finish: 72, Spin: 65},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 80,
					Vision:      models.TacticalVision{Passing: 88, Shooting: 78, Defence: 70},
				},
				Stamina: models.Stamina{Stamina: 82},
				Fitness: models.Fitness{Strength: 75, Agility: 80, InjuryTolerance: 80, InjuryResistance: 85},
			},
			{
				Number: 9, // Striker
				Form:   90, Adaptability: 85, Composure: 88,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 88, Acceleration: 90},
					Passing:  models.PassingSkill{ShortPass: 75, LongPass: 65, Cross: 72, Lob: 68, ThroughBall: 70, Chip: 72},
					Shooting: models.ShootingSkill{Power: 92, Curve: 88, Finish: 95, Spin: 90},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 92,
					Vision:      models.TacticalVision{Passing: 80, Shooting: 95, Defence: 55},
				},
				Stamina: models.Stamina{Stamina: 84},
				Fitness: models.Fitness{Strength: 80, Agility: 88, InjuryTolerance: 78, InjuryResistance: 82},
			},
			{
				Number: 10, // Attacking Midfielder
				Form:   88, Adaptability: 82, Composure: 85,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 85, Acceleration: 88},
					Passing:  models.PassingSkill{ShortPass: 90, LongPass: 82, Cross: 75, Lob: 78, ThroughBall: 88, Chip: 70},
					Shooting: models.ShootingSkill{Power: 80, Curve: 82, Finish: 86, Spin: 80},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 88,
					Vision:      models.TacticalVision{Passing: 92, Shooting: 88, Defence: 65},
				},
				Stamina: models.Stamina{Stamina: 85},
				Fitness: models.Fitness{Strength: 76, Agility: 87, InjuryTolerance: 77, InjuryResistance: 83},
			},
			{
				Number: 11, // Left Wing
				Form:   87, Adaptability: 80, Composure: 80,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 92, Acceleration: 94},
					Passing:  models.PassingSkill{ShortPass: 78, LongPass: 72, Cross: 88, Lob: 70, ThroughBall: 75, Chip: 68},
					Shooting: models.ShootingSkill{Power: 85, Curve: 83, Finish: 87, Spin: 86},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 86,
					Vision:      models.TacticalVision{Passing: 85, Shooting: 88, Defence: 60},
				},
				Stamina: models.Stamina{Stamina: 90},
				Fitness: models.Fitness{Strength: 74, Agility: 91, InjuryTolerance: 80, InjuryResistance: 82},
			},
		},
	}
}

func AwayTeam() models.Team {
	return models.Team{
		Name:      "Iron Titans",
		Morale:    82,
		Fitness:   79,
		Chemistry: 87,
		Strategy: models.Strategy{
			Tactic:    models.TacticHolding,
			Formation: models.FormationFourFourTwo,
			PlayerInstructions: map[models.PlayerNumber]models.Instruction{
				1:  {Position: models.PositionCenter}, // Goalkeeper
				2:  {Position: models.PositionWing},   // Right Back
				3:  {Position: models.PositionCenter}, // Center Back
				4:  {Position: models.PositionCenter}, // Center Back
				5:  {Position: models.PositionWing},   // Left Back
				6:  {Position: models.PositionCenter}, // Central Midfielder
				7:  {Position: models.PositionWing},   // Right Midfielder
				8:  {Position: models.PositionCenter}, // Central Midfielder
				9:  {Position: models.PositionCenter}, // Striker
				10: {Position: models.PositionCenter}, // Second Striker
				11: {Position: models.PositionWing},   // Left Midfielder
			},
			PlayStyle: models.PlayStyleDriven,
		},
		Training: models.Training{
			Focus: models.Defense,
		},
		Players: []models.Player{
			{
				Number: 1, // Goalkeeper
				Form:   78, Adaptability: 68, Composure: 80,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 45, Acceleration: 50},
					Passing:  models.PassingSkill{ShortPass: 55, LongPass: 60, Cross: 35, Lob: 45, ThroughBall: 40, Chip: 38},
					Shooting: models.ShootingSkill{Power: 25, Curve: 20, Finish: 15, Spin: 30},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 92,
					Vision:      models.TacticalVision{Passing: 55, Shooting: 25, Defence: 88},
				},
				Stamina: models.Stamina{Stamina: 60},
				Fitness: models.Fitness{Strength: 72, Agility: 82, InjuryTolerance: 90, InjuryResistance: 85},
			},
			{
				Number: 2, // Right Back
				Form:   70, Adaptability: 65, Composure: 72,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 78, Acceleration: 80},
					Passing:  models.PassingSkill{ShortPass: 68, LongPass: 60, Cross: 75, Lob: 58, ThroughBall: 60, Chip: 55},
					Shooting: models.ShootingSkill{Power: 48, Curve: 50, Finish: 42, Spin: 46},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 75,
					Vision:      models.TacticalVision{Passing: 68, Shooting: 38, Defence: 78},
				},
				Stamina: models.Stamina{Stamina: 76},
				Fitness: models.Fitness{Strength: 68, Agility: 78, InjuryTolerance: 80, InjuryResistance: 80},
			},
			{
				Number: 3, // Center Back
				Form:   75, Adaptability: 60, Composure: 78,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 65, Acceleration: 60},
					Passing:  models.PassingSkill{ShortPass: 60, LongPass: 65, Cross: 38, Lob: 50, ThroughBall: 48, Chip: 45},
					Shooting: models.ShootingSkill{Power: 55, Curve: 40, Finish: 30, Spin: 35},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 84,
					Vision:      models.TacticalVision{Passing: 60, Shooting: 30, Defence: 88},
				},
				Stamina: models.Stamina{Stamina: 70},
				Fitness: models.Fitness{Strength: 82, Agility: 65, InjuryTolerance: 78, InjuryResistance: 75},
			},
			{
				Number: 4, // Center Back
				Form:   77, Adaptability: 62, Composure: 74,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 68, Acceleration: 63},
					Passing:  models.PassingSkill{ShortPass: 65, LongPass: 60, Cross: 40, Lob: 52, ThroughBall: 48, Chip: 46},
					Shooting: models.ShootingSkill{Power: 52, Curve: 42, Finish: 32, Spin: 38},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 80,
					Vision:      models.TacticalVision{Passing: 64, Shooting: 35, Defence: 85},
				},
				Stamina: models.Stamina{Stamina: 72},
				Fitness: models.Fitness{Strength: 80, Agility: 66, InjuryTolerance: 76, InjuryResistance: 78},
			},
			{
				Number: 5, // Left Back
				Form:   73, Adaptability: 68, Composure: 70,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 82, Acceleration: 85},
					Passing:  models.PassingSkill{ShortPass: 70, LongPass: 65, Cross: 78, Lob: 60, ThroughBall: 68, Chip: 62},
					Shooting: models.ShootingSkill{Power: 52, Curve: 54, Finish: 45, Spin: 50},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 78,
					Vision:      models.TacticalVision{Passing: 70, Shooting: 48, Defence: 80},
				},
				Stamina: models.Stamina{Stamina: 82},
				Fitness: models.Fitness{Strength: 72, Agility: 85, InjuryTolerance: 80, InjuryResistance: 82},
			},
			{
				Number: 6, // Central Midfielder
				Form:   80, Adaptability: 75, Composure: 78,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 75, Acceleration: 78},
					Passing:  models.PassingSkill{ShortPass: 82, LongPass: 78, Cross: 65, Lob: 70, ThroughBall: 80, Chip: 68},
					Shooting: models.ShootingSkill{Power: 68, Curve: 65, Finish: 60, Spin: 62},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 82,
					Vision:      models.TacticalVision{Passing: 85, Shooting: 72, Defence: 70},
				},
				Stamina: models.Stamina{Stamina: 85},
				Fitness: models.Fitness{Strength: 75, Agility: 78, InjuryTolerance: 78, InjuryResistance: 80},
			},
			{
				Number: 7, // Right Midfielder
				Form:   79, Adaptability: 74, Composure: 76,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 85, Acceleration: 90},
					Passing:  models.PassingSkill{ShortPass: 75, LongPass: 70, Cross: 80, Lob: 68, ThroughBall: 72, Chip: 70},
					Shooting: models.ShootingSkill{Power: 70, Curve: 75, Finish: 72, Spin: 78},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 80,
					Vision:      models.TacticalVision{Passing: 78, Shooting: 70, Defence: 65},
				},
				Stamina: models.Stamina{Stamina: 88},
				Fitness: models.Fitness{Strength: 70, Agility: 88, InjuryTolerance: 77, InjuryResistance: 80},
			},
			{
				Number: 8, // Central Midfielder
				Form:   78, Adaptability: 72, Composure: 80,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 77, Acceleration: 79},
					Passing:  models.PassingSkill{ShortPass: 80, LongPass: 75, Cross: 68, Lob: 72, ThroughBall: 78, Chip: 65},
					Shooting: models.ShootingSkill{Power: 72, Curve: 70, Finish: 65, Spin: 66},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 84,
					Vision:      models.TacticalVision{Passing: 83, Shooting: 75, Defence: 72},
				},
				Stamina: models.Stamina{Stamina: 84},
				Fitness: models.Fitness{Strength: 73, Agility: 80, InjuryTolerance: 80, InjuryResistance: 82},
			},
			{
				Number: 9, // Striker
				Form:   85, Adaptability: 82, Composure: 88,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 88, Acceleration: 90},
					Passing:  models.PassingSkill{ShortPass: 70, LongPass: 65, Cross: 68, Lob: 65, ThroughBall: 70, Chip: 72},
					Shooting: models.ShootingSkill{Power: 90, Curve: 85, Finish: 94, Spin: 88},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 90,
					Vision:      models.TacticalVision{Passing: 75, Shooting: 92, Defence: 58},
				},
				Stamina: models.Stamina{Stamina: 80},
				Fitness: models.Fitness{Strength: 80, Agility: 88, InjuryTolerance: 78, InjuryResistance: 84},
			},
			{
				Number: 10, // Second Striker
				Form:   83, Adaptability: 80, Composure: 85,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 84, Acceleration: 86},
					Passing:  models.PassingSkill{ShortPass: 78, LongPass: 70, Cross: 68, Lob: 66, ThroughBall: 82, Chip: 70},
					Shooting: models.ShootingSkill{Power: 82, Curve: 80, Finish: 88, Spin: 85},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 88,
					Vision:      models.TacticalVision{Passing: 80, Shooting: 90, Defence: 60},
				},
				Stamina: models.Stamina{Stamina: 82},
				Fitness: models.Fitness{Strength: 78, Agility: 85, InjuryTolerance: 75, InjuryResistance: 80},
			},
			{
				Number: 11, // Left Midfielder
				Form:   81, Adaptability: 76, Composure: 78,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 89, Acceleration: 91},
					Passing:  models.PassingSkill{ShortPass: 72, LongPass: 70, Cross: 85, Lob: 70, ThroughBall: 74, Chip: 68},
					Shooting: models.ShootingSkill{Power: 78, Curve: 82, Finish: 80, Spin: 80},
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 83,
					Vision:      models.TacticalVision{Passing: 78, Shooting: 82, Defence: 65},
				},
				Stamina: models.Stamina{Stamina: 87},
				Fitness: models.Fitness{Strength: 72, Agility: 89, InjuryTolerance: 80, InjuryResistance: 82},
			},
		},
	}
}
