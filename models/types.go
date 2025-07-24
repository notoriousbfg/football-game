package models

type Team struct {
	Name      string
	Strategy  Strategy
	Morale    int
	Fitness   int
	Chemistry int
	Players   []Player
	Training  Training
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
	Number               PlayerNumber
	Form                 int
	Adaptability         int
	Composure            int
	Technical            TechnicalSkill
	TacticalIntelligence TacticalIntelligence
	Stamina              Stamina
	Fitness              Fitness
}

type PlayerNumber int

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
