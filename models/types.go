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
	Position   PlayerPosition
	Name       string
	Number     PlayerNumber
	Exclusions map[PlayerNumber]bool
}

func (t *Team) SearchPlayers(options PlayerSearchOptions) Player {
	type PlayerScore struct {
		Player Player
		Score  int
	}
	scores := make(map[PlayerNumber]PlayerScore, 0)
	for _, player := range t.Players {
		score := PlayerScore{
			Player: player,
			Score:  0,
		}
		if player.Position == options.Position {
			score.Score += 2
		}
		if similar, ok := SimilarPositions[player.Position]; ok {
			for _, s := range similar {
				if options.Position == s {
					score.Score += 1
				}
			}
		}
		if player.Name == options.Name {
			score.Score += 2
		}
		if player.Number == options.Number {
			score.Score += 2
		}
		if _, found := options.Exclusions[player.Number]; found {
			score.Score = 0
		}
		scores[player.Number] = score
	}
	var highest *PlayerScore
	for _, score := range scores {
		if score.Score == 0 {
			continue
		}
		if highest == nil {
			highest = &score
		}
		if score.Score > highest.Score {
			highest = &score
		}
	}
	if highest == nil {
		panic(fmt.Errorf("no player found with options (position: %s, exclusions: %+v)", options.Position.String(), options.Exclusions))
	}
	return highest.Player
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

var SimilarPositions = map[PlayerPosition][]PlayerPosition{
	Goalkeeper: {},

	RightBack:       {RightWingBack},
	RightWingBack:   {RightBack, RightMidfielder, RightWinger},
	RightMidfielder: {RightWingBack, RightWinger},
	RightWinger:     {RightMidfielder, RightWingBack},

	LeftCentreBack:  {CentralDefensiveMidfielder, RightCentreBack},
	RightCentreBack: {CentralDefensiveMidfielder, LeftCentreBack},

	LeftBack:       {LeftWingBack},
	LeftWingBack:   {LeftBack, LeftMidfielder, LeftWinger},
	LeftMidfielder: {LeftWingBack, LeftWinger},
	LeftWinger:     {LeftMidfielder, LeftWingBack},

	CentralMidfielder:          {CentralDefensiveMidfielder, CentralAttackingMidfielder},
	CentralDefensiveMidfielder: {CentralMidfielder, LeftCentreBack, RightCentreBack},
	CentralAttackingMidfielder: {CentralMidfielder, CentreForward},

	CentreForward: {Striker, CentralAttackingMidfielder},
	Striker:       {CentreForward},
}

type TechnicalSkill struct {
	Speed     SpeedSkill
	Passing   PassingSkill
	Shooting  ShootingSkill
	Defending DefendingSkill
	FreeKicks int
	Penalties int
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
	Power  int
	Curve  int
	Finish int
	Spin   int
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
