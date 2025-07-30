package models

import (
	"fmt"
	"strings"
)

type Team struct {
	Name      string
	Strategy  Strategy
	Morale    int
	Fitness   int
	Chemistry int
	Players   []Player
	Training  Training
}

type PlayerSearchOptions struct {
	Positions  []PlayerPosition
	Name       string
	Number     PlayerNumber
	Exclusions map[PlayerNumber]string
}

type Strategy struct {
	Tactic             Tactic
	Formation          Formation
	PlayerInstructions map[PlayerNumber]Instruction
	PlayStyle          PlayStyle
}

type Tactic int

const (
	TacticCounter Tactic = iota
	TacticPressing
	TacticDefensive
	TacticHolding
)

type Formation int

const (
	FormationFourThreeThree Formation = iota
	FormationFourFourTwo
	FormationThreeFourTwoOne
)

type Instruction struct {
	Position PositionInstruction
}

type PositionInstruction int

const (
	PositionWing PositionInstruction = iota
	PositionCenter
)

type PlayStyle int

const (
	PlayStyleCreative PlayStyle = iota
	PlayStylePredictable
	PlayStyleDriven
	PlayStyleCrossing
	PlayStyleDefensive
)

type Training struct {
	Focus TrainingFocus
}

type TrainingFocus int

const (
	Passing TrainingFocus = iota
	Defense
	Shooting
	Penalties
	SetPieces
)

type Player struct {
	Name                 string
	Position             PlayerPosition
	Number               PlayerNumber
	Form                 int
	Adaptability         int
	Composure            int
	Technical            TechnicalSkill
	TacticalIntelligence TacticalIntelligence
	Stamina              Stamina
	Fitness              Fitness
}

func (p Player) Initials() string {
	words := strings.Fields(p.Name)
	initials := ""
	if len(words) > 2 {
		words = []string{words[0], words[len(words)-1]}
	}
	for _, word := range words {
		if len(word) > 0 {
			initials += strings.ToUpper(string(word[0]))
		}
	}
	if len(initials) < 2 {
		initials = fmt.Sprintf("%s ", initials)
	}
	return initials
}

type PlayerNumber int

//go:generate stringer -type=PlayerPosition -output player_positions_string.go
type PlayerPosition int

const (
	Goalkeeper PlayerPosition = iota
	RightBack
	RightWingBack
	LeftCentreBack
	RightCentreBack
	LeftBack
	LeftWingBack
	LeftMidfielder
	LeftWinger
	CentralMidfielder
	CentralDefensiveMidfielder
	CentralAttackingMidfielder
	CentreForward
	RightMidfielder
	RightWinger
	Striker
)

var Forwards = []PlayerPosition{
	Striker,
	LeftWinger,
	RightWinger,
	CentreForward,
}

var Midfielders = []PlayerPosition{
	CentralMidfielder,
	CentralDefensiveMidfielder,
	CentralAttackingMidfielder,
	RightMidfielder,
	LeftMidfielder,
}

var Defenders = []PlayerPosition{
	RightCentreBack,
	LeftCentreBack,
	LeftBack,
	RightBack,
	CentralDefensiveMidfielder,
	Goalkeeper,
}

var SimilarPositions = map[PlayerPosition][]PlayerPosition{
	Goalkeeper: {},

	RightBack:       {RightWingBack, RightCentreBack},
	RightWingBack:   {RightBack, RightMidfielder, RightWinger},
	RightMidfielder: {RightWingBack, RightWinger},
	RightWinger:     {RightMidfielder, RightWingBack},

	LeftCentreBack:  {CentralDefensiveMidfielder, RightCentreBack},
	RightCentreBack: {CentralDefensiveMidfielder, LeftCentreBack},

	LeftBack:       {LeftWingBack, LeftCentreBack},
	LeftWingBack:   {LeftBack, LeftMidfielder, LeftWinger},
	LeftMidfielder: {LeftWingBack, LeftWinger},
	LeftWinger:     {LeftMidfielder, LeftWingBack},

	CentralMidfielder:          {CentralDefensiveMidfielder, CentralAttackingMidfielder},
	CentralDefensiveMidfielder: {CentralMidfielder, LeftCentreBack, RightCentreBack},
	CentralAttackingMidfielder: {CentralMidfielder, CentreForward},

	CentreForward: {Striker, CentralAttackingMidfielder},
	Striker:       {CentreForward, LeftWinger, RightWinger},
}

var TeammateAdjacents = map[PlayerPosition][]PlayerPosition{
	Goalkeeper: {LeftCentreBack, RightCentreBack, LeftBack, RightBack},

	RightBack:       {RightWingBack, RightCentreBack},
	RightWingBack:   {RightBack, RightMidfielder, RightWinger},
	RightMidfielder: {RightWingBack, RightWinger},
	RightWinger:     {RightMidfielder, RightWingBack},

	LeftCentreBack:  {CentralDefensiveMidfielder, RightCentreBack},
	RightCentreBack: {CentralDefensiveMidfielder, LeftCentreBack},

	LeftBack:       {LeftWingBack, LeftCentreBack},
	LeftWingBack:   {LeftBack, LeftMidfielder, LeftWinger},
	LeftMidfielder: {LeftWingBack, LeftWinger},
	LeftWinger:     {LeftMidfielder, LeftWingBack},

	CentralMidfielder:          {CentralDefensiveMidfielder, CentralAttackingMidfielder, CentralMidfielder},
	CentralDefensiveMidfielder: {CentralMidfielder, LeftCentreBack, RightCentreBack, CentralAttackingMidfielder},
	CentralAttackingMidfielder: {CentralMidfielder, CentreForward, Striker, LeftWinger, RightWinger},

	CentreForward: {Striker, CentralAttackingMidfielder},
	Striker:       {CentreForward, LeftWinger, RightWinger},
}

var OpponentAdjacents = map[PlayerPosition][]PlayerPosition{
	Goalkeeper: {Striker, LeftWinger, RightWinger, CentreForward},

	RightBack:       {LeftWinger, LeftMidfielder},
	RightWingBack:   {LeftWinger, LeftMidfielder},
	RightMidfielder: {LeftMidfielder, CentralMidfielder},
	RightWinger:     {LeftWinger, LeftMidfielder},

	LeftCentreBack:  {CentreForward, Striker, RightWinger},
	RightCentreBack: {CentreForward, Striker, LeftWinger},

	LeftBack:       {RightWinger, RightMidfielder},
	LeftWingBack:   {RightWinger, RightMidfielder},
	LeftMidfielder: {RightMidfielder, CentralMidfielder},
	LeftWinger:     {RightWinger, RightMidfielder},

	CentralMidfielder:          {CentralDefensiveMidfielder, RightCentreBack, LeftCentreBack},
	CentralDefensiveMidfielder: {CentralAttackingMidfielder, CentreForward},
	CentralAttackingMidfielder: {CentralDefensiveMidfielder, CentreForward},

	CentreForward: {LeftCentreBack, RightCentreBack, CentralDefensiveMidfielder},
	Striker:       {LeftCentreBack, RightCentreBack, CentralDefensiveMidfielder},
}

type TechnicalSkill struct {
	Speed       SpeedSkill
	Passing     PassingSkill
	Shooting    ShootingSkill
	Defending   DefendingSkill
	Dribbling   DribblingSkill
	Goalkeeping GoalkeepingSkill
	FreeKicks   int
	Penalties   int
}

type SpeedSkill struct {
	Speed        int
	Acceleration int
}

type PassingSkill struct {
	ShortPass   int
	LongPass    int
	Cross       int
	Lob         int
	ThroughBall int
	Chip        int
}

type ShootingSkill struct {
	Power     int
	Curve     int
	Finishing int
	Spin      int
}

type DefendingSkill struct {
	Jumping       int
	Interceptions int
	Heading       HeadingSkill
	Blocking      int
}

type HeadingSkill struct {
	Accuracy int
	Power    int
}

type DribblingSkill struct {
	SkillMoves int
	Agility    int
	Dribbling  int
}

type GoalkeepingSkill struct {
	Reflexes    int
	Positioning int
	Reactions   int
}

type TacticalIntelligence struct {
	Positioning int
	Vision      TacticalVision
}

type TacticalVision struct {
	Passing  int // i.e. when to pass and who to pass to
	Shooting int // i.e. when to shoot
	Defence  int // i.e. when to drop back
}

type Stamina struct {
	Stamina int
}

type Fitness struct {
	Strength         int
	Agility          int
	InjuryTolerance  int
	InjuryResistance int
}
