package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/notoriousbfg/football-game/helpers"
	"github.com/notoriousbfg/football-game/models"
)

type Simulation struct {
	Match             Match
	KickoffTeam       models.Team
	State             *SimulationState
	RandomFloat       func() float64
	SynergyMultiplier int
	TacticalCounters  map[int]TacticalCounter
	EventTriggers     map[EventType]EventTrigger
	Pitch             *Pitch
}

func (sim *Simulation) Run() {
	sim.Pitch = NewPitch(&sim.Match)
	sim.Pitch.Draw()

	// halfSeconds := 45 * 60

	// first half
	sim.State.CaptureEvent(
		sim.State.startingEvent(sim.KickoffTeam),
	)

	// sim.runPeriod(0, halfSeconds)
	// sim.runPeriod(halfSeconds, sim.State.FirstHalfExtraTime)

	// // second half
	// sim.runPeriod(halfSeconds, halfSeconds)
	// sim.runPeriod(halfSeconds, sim.State.SecondHalfExtraTime)

	sim.State.Outcome = &Outcome{
		HomeScore: sim.State.HomeScore,
		AwayScore: sim.State.AwayScore,
	}
}

func (sim *Simulation) runPeriod(start, seconds int) {
	for i := start; i <= start+seconds; i++ {
		sim.State.Time = sim.State.Time.Add(time.Second)
	}
}

func (sim *Simulation) opposingTeam(team models.Team) models.Team {
	if sim.Match.A.Name == team.Name {
		return sim.Match.H
	} else {
		return sim.Match.A
	}
}

// things that change
type SimulationState struct {
	Simulation          *Simulation
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
	RandomFloat         func() float64
	Events              []Event
	Outcome             *Outcome
}

func (s *SimulationState) CaptureEvent(e Event) {
	s.Events = append(s.Events, e)

	switch e.Type {
	case ETPass:
		s.Time = s.Time.Add(time.Second * 3)
	case ETDribble:
		s.Time = s.Time.Add(time.Second * 5)
	case ETInterception:
		s.Time = s.Time.Add(time.Second * 2)
	case ETYellowCard:
		switch e.Team.Name {
		case s.Simulation.Match.H.Name:
			s.HomeYellowCards++
		case s.Simulation.Match.A.Name:
			s.AwayYellowCards++
		}
		s.Time = s.Time.Add(time.Second * 30)
	case ETRedCard:
		switch e.Team.Name {
		case s.Simulation.Match.H.Name:
			s.HomeRedCards++
		case s.Simulation.Match.A.Name:
			s.HomeRedCards++
		}
		s.Time = s.Time.Add(time.Second * 30)
	}

	if trigger, exists := s.Simulation.EventTriggers[e.Type]; exists {
		trigger(e, s)
	}
}

func (s *SimulationState) LastEvent() Event {
	return s.Events[len(s.Events)-1]
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
	return ETFreeKickOnGoal
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
		Events:            make([]Event, 0),
	}

	sim := &Simulation{
		Match:             Match{H: home, A: away},
		State:             state,
		SynergyMultiplier: 1,
		TacticalCounters:  make(map[int]TacticalCounter),
		RandomFloat:       randomFloat,
	}

	state.Simulation = sim

	coinFlip := randomFloat()
	if coinFlip < 0.5 {
		sim.KickoffTeam = home
	} else {
		sim.KickoffTeam = away
	}

	sim.EventTriggers = make(map[EventType]EventTrigger)

	sim.EventTriggers[ETPass] = func(e Event, s *SimulationState) {
		s.logPass(e)
		s.CaptureEvent(
			s.evaluatePass(&e.Team),
		)
	}

	sim.EventTriggers[ETDribble] = func(e Event, s *SimulationState) {
		s.logDribble(e)
		s.CaptureEvent(e)
	}

	sim.EventTriggers[ETInterception] = func(e Event, s *SimulationState) {
		s.logInterception(e)
		s.CaptureEvent(e)
	}

	sim.EventTriggers[ETYellowCard] = func(e Event, s *SimulationState) {
		s.CaptureEvent(e)
	}

	sim.EventTriggers[ETRedCard] = func(e Event, s *SimulationState) {
		s.CaptureEvent(e)
	}

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

	return sim
}

func (s *SimulationState) startingEvent(team models.Team) Event {
	startingPlayer := team.SearchPlayers(models.PlayerSearchOptions{
		Position: models.Striker,
	})

	receivingPlayer := team.SearchPlayers(models.PlayerSearchOptions{
		Position: models.CentralMidfielder,
	})

	return Event{
		Type:            ETPass,
		Team:            team,
		StartingPlayer:  &startingPlayer,
		FinishingPlayer: &receivingPlayer,
		EventMeta: EventMeta{
			"quality": 100,
		},
	}
}

func (s *SimulationState) evaluateAttack(randFloat func() float64) bool {
	// get last event
	// lastEvent := s.LastEvent()

	// get team/player from event source

	// pass, shoot, dribble?

	// evaluate likelihood of chance

	// evaluate likelihood of goal

	// weighted dice roll

	return false
}

func (s *SimulationState) evaluatePass(team *models.Team) Event {
	lastEvent := s.LastEvent()
	player := lastEvent.FinishingPlayer

	// decide what to do next
	// pass, dribble, shoot
	decision := s.makeDecision()

	return s.evaluateDecision(*team, player, decision)
}

func (s *SimulationState) evaluateDecision(team models.Team, player *models.Player, decision Decision) Event {
	switch decision {
	case DecisionLongPass:
		success := s.evaluateLongPass(*player)
		receivingPlayer := team.SearchPlayers(models.PlayerSearchOptions{
			Position:   models.Striker,
			Exclusions: map[models.PlayerNumber]string{player.Number: player.Initials()},
		})
		if success {
			return Event{
				Type:            ETPass,
				Team:            team,
				StartingPlayer:  player,
				FinishingPlayer: &receivingPlayer,
			}
		} else {
			opposingTeam := s.Simulation.opposingTeam(team)
			interceptor := s.nearestTo(receivingPlayer.Position, opposingTeam.Players)
			return Event{
				Type:            ETInterception,
				Team:            opposingTeam,
				StartingPlayer:  player,
				FinishingPlayer: interceptor,
			}
		}
	case DecisionShortPass:
		success := s.evaluateShortPass(*player)
		receivingPlayer := team.SearchPlayers(models.PlayerSearchOptions{
			Position:   models.Striker,
			Exclusions: map[models.PlayerNumber]string{player.Number: player.Initials()},
		})
		if success {
			return Event{
				Type:            ETPass,
				Team:            team,
				StartingPlayer:  player,
				FinishingPlayer: &receivingPlayer,
			}
		} else {
			opposingTeam := s.Simulation.opposingTeam(team)
			interceptor := s.nearestTo(receivingPlayer.Position, opposingTeam.Players)
			return Event{
				Type:            ETInterception,
				Team:            opposingTeam,
				StartingPlayer:  player,
				FinishingPlayer: interceptor,
			}
		}
	case DecisionDribble:
		opposingTeam := s.Simulation.opposingTeam(team)
		success := s.evaluateDribble(*player, opposingTeam)
		if success {
			return Event{
				Type:            ETDribble,
				Team:            team,
				StartingPlayer:  player,
				FinishingPlayer: player,
			}
		} else {
			interceptor := s.nearestTo(player.Position, opposingTeam.Players)
			return Event{
				Type:            ETInterception,
				Team:            s.Simulation.opposingTeam(team),
				StartingPlayer:  player,
				FinishingPlayer: interceptor,
			}
		}
	case NoDecision:
		success := s.evaluateHold(*player)
		if success {
			return Event{}
		} else {
			return Event{
				Type: ETInterception,
			}
		}
	default:
		panic(fmt.Errorf("decision '%s' not handled", decision.String()))
	}
}

func (s *SimulationState) evaluateLongPass(player models.Player) bool {
	skill := player.Technical.Passing.LongPass
	vision := player.TacticalIntelligence.Vision.Passing
	agility := player.Fitness.Agility

	visionWeight := 0.4
	executionWeight := 0.4
	agilityWeight := 0.2

	successChance := float64(vision)*visionWeight +
		float64(skill)*executionWeight +
		float64(agility)*agilityWeight

	successChance /= 100.0

	return s.Simulation.RandomFloat() < successChance
}

func (s *SimulationState) evaluateShortPass(player models.Player) bool {
	skill := player.Technical.Passing.ShortPass
	vision := player.TacticalIntelligence.Vision.Passing
	agility := player.Fitness.Agility

	visionWeight := 0.2
	executionWeight := 0.6
	agilityWeight := 0.2

	successChance := float64(vision)*visionWeight +
		float64(skill)*executionWeight +
		float64(agility)*agilityWeight

	successChance /= 100.0

	return s.Simulation.RandomFloat() < successChance
}

func (s *SimulationState) evaluateDribble(player models.Player, opposingTeam models.Team) bool {
	// base stats from the player
	skill := float64(player.Technical.Dribbling.Dribbling)
	agility := float64(player.Technical.Dribbling.Agility)
	composure := float64(player.Composure)

	// compute individual ability score
	dribbleScore := skill*0.5 + agility*0.3 + composure*0.2

	// penalize for nearby opponents
	for _, defender := range opposingTeam.Players {
		proximity := s.positionProximityScore(player.Position, defender.Position)
		if proximity > 0 {
			dribbleScore -= float64(defender.Technical.Defending.Interceptions) * proximity
		}
	}

	// adjust for opposing team’s overall defensive acumen
	totalDef := 0.0
	for _, opp := range opposingTeam.Players {
		totalDef += (float64(opp.Technical.Defending.Interceptions) * 0.5) + (float64(opp.Technical.Defending.Blocking) * 0.5)
	}
	avgTeamDef := totalDef / float64(len(opposingTeam.Players))
	dribbleScore -= avgTeamDef * 0.1 // global pressure factor

	// clamp the score
	if dribbleScore < 0 {
		dribbleScore = 0
	} else if dribbleScore > 100 {
		dribbleScore = 100
	}

	dribbleScore /= 100.0

	return s.Simulation.RandomFloat() < dribbleScore
}

func (s *SimulationState) positionProximityScore(attackerPos, defenderPos models.PlayerPosition) float64 {
	if attackerPos == defenderPos {
		return 1.0
	}

	adjacentPositions, ok := models.OpponentAdjacents[attackerPos]
	if !ok {
		return 0.0
	}

	for _, pos := range adjacentPositions {
		if pos == defenderPos {
			return 0.75
		}
	}

	for _, adjPos := range adjacentPositions {
		secondDegreePositions, ok := models.OpponentAdjacents[adjPos]
		if !ok {
			continue
		}
		for _, pos := range secondDegreePositions {
			if pos == defenderPos {
				return 0.4
			}
		}
	}

	return 0.0
}

func (s *SimulationState) nearestTo(attackerPos models.PlayerPosition, opponents []models.Player) *models.Player {
	var (
		bestScore float64
		closest   *models.Player
	)

	for i := range opponents {
		score := s.positionProximityScore(attackerPos, opponents[i].Position)
		if score > bestScore {
			bestScore = score
			closest = &opponents[i]
		}
	}

	return closest
}

func (s *SimulationState) evaluateHold(player models.Player) bool {
	agility := player.Fitness.Agility
	strength := player.Fitness.Strength
	vision := player.TacticalIntelligence.Vision.Passing

	agilityWeight := 0.4
	strengthWeight := 0.4
	visionWeight := 0.2

	successChance := float64(agility)*agilityWeight +
		float64(strength)*strengthWeight +
		float64(vision)*visionWeight

	successChance /= 100.0

	return s.Simulation.RandomFloat() < successChance
}

func (s *SimulationState) makeDecision() Decision {
	lastEvent := s.LastEvent()
	player := lastEvent.FinishingPlayer
	if player == nil {
		return NoDecision
	}

	vision := float64(player.TacticalIntelligence.Vision.Passing)
	dribbling := float64(player.Technical.Dribbling.Dribbling)
	shooting := float64(player.Technical.Shooting.Finishing)
	crossing := float64(player.Technical.Passing.Cross)
	position := player.Position

	rng := s.Simulation.RandomFloat()

	// shooting — more likely if forward or attacking midfielder
	if helpers.IsAttacker(position) && shooting > 70 && rng < shooting/200.0 {
		return DecisionShoot
	}

	// crossing — more likely for wingers and wide backs
	if helpers.IsWinger(position) && crossing > 60 && rng < crossing/200.0 {
		return DecisionCross
	}

	// dribbling — influenced by agility and dribbling
	agility := float64(player.Fitness.Agility)
	dribbleChance := (dribbling*0.6 + agility*0.4) / 100.0
	if rng < dribbleChance*0.6 {
		return DecisionDribble
	}

	// passing — fallback with weighted short/long
	return s.decidePassType(vision)
}

func (s *SimulationState) decidePassType(vision float64) Decision {
	longPassProbability := vision * 0.8 // up to 80% chance

	if s.Simulation.RandomFloat() < longPassProbability {
		return DecisionLongPass
	}
	return DecisionShortPass
}

func (s *SimulationState) logPass(e Event) {
	fmt.Printf("(%s) %s passes to %s\n", s.Timestamp(), e.StartingPlayer.Name, e.FinishingPlayer.Name)
}

func (s *SimulationState) logDribble(e Event) {
	fmt.Printf("(%s) %s is dribbling with the ball\n", s.Timestamp(), e.StartingPlayer.Name)
}

func (s *SimulationState) logInterception(e Event) {
	fmt.Printf("(%s) %s loses the ball to %s\n", s.Timestamp(), e.StartingPlayer.Name, e.FinishingPlayer.Name)
}
