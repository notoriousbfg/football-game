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
		Score  float64
	}

	similarityWeights := make(map[PlayerPosition]float64)
	similarityWeights[options.Position] = 2.0 // direct match

	visited := map[PlayerPosition]bool{options.Position: true}
	queue := []struct {
		Pos   PlayerPosition
		Depth int
	}{{options.Position, 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.Depth >= 3 {
			continue
		}

		for _, similar := range SimilarPositions[current.Pos] {
			if !visited[similar] {
				visited[similar] = true

				weight := map[int]float64{
					0: 2.0,  // direct match
					1: 1.0,  // 1st-degree similar
					2: 0.5,  // 2nd-degree similar
					3: 0.25, // 3rd-degree similar
				}[current.Depth+1]

				similarityWeights[similar] = weight

				queue = append(queue, struct {
					Pos   PlayerPosition
					Depth int
				}{similar, current.Depth + 1})
			}
		}
	}

	var highest *PlayerScore
	for _, player := range t.Players {
		if _, excluded := options.Exclusions[player.Number]; excluded {
			continue
		}

		score := 0.0

		if weight, ok := similarityWeights[player.Position]; ok {
			score += weight
		}

		if player.Name == options.Name {
			score += 2
		}

		if player.Number == options.Number {
			score += 2
		}

		if score == 0 {
			continue
		}

		if highest == nil || score > highest.Score {
			highest = &PlayerScore{Player: player, Score: score}
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
	Striker:       {CentreForward, LeftWinger, RightWinger},
}

type TechnicalSkill struct {
	Speed     SpeedSkill
	Passing   PassingSkill
	Shooting  ShootingSkill
	Defending DefendingSkill
	Dribbling int
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
