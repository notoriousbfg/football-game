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

// func (sim *Simulation) runPeriod(start, seconds int) {
// 	for i := start; i <= start+seconds; i++ {
// 		sim.State.Time = sim.State.Time.Add(time.Second)
// 	}
// }

func (sim *Simulation) opposingTeam(team models.Team) models.Team {
	if sim.Match.A.Name == team.Name {
		return sim.Match.H
	} else {
		return sim.Match.A
	}
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

	sim.registerTriggers()

	// triggers[HomeTeamAttacking] = func(e Event, s *SimulationState) {
	// 	if s.evaluateAttack(e, randomFloat) {
	// 		s.HomeScore++
	// 		s.HomeMomentum += 0.1
	// 		s.AwayMomentum -= 0.1
	// 	}
	// 	e.Log.log(s)
	// }

	// triggers[AwayTeamAttacking] = func(e Event, s *SimulationState) {
	// 	if s.evaluateAttack(e, randomFloat) {
	// 		s.AwayScore++
	// 		s.AwayMomentum += 0.1
	// 		s.HomeMomentum -= 0.1
	// 	}
	// 	e.Log.log(s)
	// }

	return sim
}

func (sim *Simulation) registerTriggers() {
	sim.EventTriggers = make(map[EventType]EventTrigger)

	sim.EventTriggers[ETEndOfFirstHalf] = func(e Event, s *SimulationState) {

	}

	sim.EventTriggers[ETPass] = func(e Event, s *SimulationState) {
		s.Log.logPass(s, e)
		s.Time = s.Time.Add(time.Second * 3)
		s.CaptureEvent(
			s.action(e),
		)
	}

	sim.EventTriggers[ETCross] = func(e Event, s *SimulationState) {
		s.Log.logCross(s, e)
		s.Time = s.Time.Add(time.Second * 3)
		s.CaptureEvent(
			s.action(e),
		)
	}

	sim.EventTriggers[ETDribble] = func(e Event, s *SimulationState) {
		s.Log.logDribble(s, e)
		s.Time = s.Time.Add(time.Second * 5)
		s.CaptureEvent(
			s.action(e),
		)
	}

	sim.EventTriggers[ETGoal] = func(e Event, s *SimulationState) {
		s.Log.logGoal(s, e)
		s.Time = s.Time.Add(time.Minute * 2)
		if s.isHome(e.Team) {
			s.HomeScore++
			s.HomeMomentum += 0.1
		} else {
			s.AwayScore++
			s.AwayMomentum += 0.1
		}
		s.Log.logRestart(s)
		s.CaptureEvent(
			s.reset(e),
		)
	}

	sim.EventTriggers[ETInterception] = func(e Event, s *SimulationState) {
		s.Log.logInterception(s, e)
		s.Time = s.Time.Add(time.Second * 3)
		s.CaptureEvent(
			s.action(e),
		)
	}

	sim.EventTriggers[ETPossession] = func(e Event, s *SimulationState) {
		s.Log.logPossession(s, e)
		s.Time = s.Time.Add(time.Second * 3)
		s.CaptureEvent(
			s.action(e),
		)
	}

	sim.EventTriggers[ETYellowCard] = func(e Event, s *SimulationState) {
		s.Log.logYellowCard(s, e)
		s.CaptureEvent(
			s.freeKick(e),
		)
		s.Time = s.Time.Add(time.Second * 3)
		if sim.State.isHome(e.Team) {
			sim.State.HomeYellowCards++
		} else {
			sim.State.AwayYellowCards++
		}
	}

	sim.EventTriggers[ETRedCard] = func(e Event, s *SimulationState) {
		s.CaptureEvent(
			s.freeKick(e),
		)
		if sim.State.isHome(e.Team) {
			sim.State.HomeRedCards++
		} else {
			sim.State.HomeRedCards++
		}
	}

	sim.EventTriggers[ETSave] = func(e Event, s *SimulationState) {
		s.Log.logSave(s, e)
		s.CaptureEvent(
			s.goalKeeperKick(e),
		)
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
	FirstHalfExtraTime  time.Duration // seconds
	SecondHalfExtraTime time.Duration // seconds
	RandomFloat         func() float64
	Events              []Event
	Outcome             *Outcome
	Log                 *Log
}

func (s *SimulationState) CaptureEvent(e Event) {
	s.Events = append(s.Events, e)

	if trigger, exists := s.Simulation.EventTriggers[e.Type]; exists {
		trigger(e, s)
	} else {
		s.Log.logMissingTrigger(s, e)
	}
}

func (s *SimulationState) LastEvent() Event {
	return s.Events[len(s.Events)-1]
}

func (s *SimulationState) Timestamp() string {
	duration := s.Time.Sub(s.Start)

	minutes := duration.Minutes()
	if duration.Hours() > 0 {
		minutes += duration.Hours() * 60
	}

	return fmt.Sprintf("%02d:%02d", int(minutes)%90, int(duration.Seconds())%60)
}

func (s *SimulationState) trackProgress() *Event {
	if s.Time.After(s.Start.Add((time.Second * 90 * 60) + s.SecondHalfExtraTime)) {
		return &Event{Type: ETEndOfSecondHalfExtraTime}
	} else if s.Time.After(s.Start.Add(time.Second * 90 * 60)) {
		return &Event{Type: ETEndOfSecondHalf}
	} else if s.Time.After(s.Start.Add((time.Second * 45 * 60) + s.FirstHalfExtraTime)) {
		return &Event{Type: ETEndOfFirstHalfExtraTime}
	} else if s.Time.After(s.Start.Add((time.Second * 45 * 60))) {
		return &Event{Type: ETEndOfFirstHalf}
	}
	return nil
}

func (s *SimulationState) isHome(team models.Team) bool {
	return team.Name == s.Simulation.Match.H.Name
}

type Match struct {
	H models.Team
	A models.Team
}

type TacticalCounter struct{}

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

func (s *SimulationState) goalKeeperKick(e Event) Event {
	receivingPlayer := e.Team.SearchPlayers(models.PlayerSearchOptions{
		Position: models.CentralMidfielder,
	})
	return Event{
		Type:            ETPass,
		Team:            e.Team,
		StartingPlayer:  e.FinishingPlayer,
		FinishingPlayer: &receivingPlayer,
	}
}

func (s *SimulationState) reset(e Event) Event {
	opposingTeam := s.Simulation.opposingTeam(e.Team)
	startingPlayer := e.Team.SearchPlayers(models.PlayerSearchOptions{
		Position: models.Striker,
	})
	receivingPlayer := e.Team.SearchPlayers(models.PlayerSearchOptions{
		Position: models.CentralMidfielder,
	})
	return Event{
		Type:            ETPass,
		Team:            opposingTeam,
		StartingPlayer:  &startingPlayer,
		FinishingPlayer: &receivingPlayer,
		EventMeta: EventMeta{
			"quality": 100,
		},
	}
}

func (s *SimulationState) freeKick(e Event) Event {
	opposingTeam := s.Simulation.opposingTeam(e.Team)
	kickTaker := s.opponentNearestTo(e.FinishingPlayer.Position, opposingTeam.Players)
	return Event{
		Type:            ETFreeKickOnGoal,
		Team:            opposingTeam,
		StartingPlayer:  kickTaker,
		FinishingPlayer: kickTaker,
	}
}

func (s *SimulationState) action(e Event) Event {
	player := e.FinishingPlayer

	decision := s.makeDecision(e)

	return s.evaluateDecision(e.Team, player, decision)
}

func (s *SimulationState) evaluateDecision(team models.Team, player *models.Player, decision Decision) Event {
	switch decision {
	case DecisionLongPass:
		receivingPlayer := team.SearchPlayers(models.PlayerSearchOptions{
			Position:   models.Striker,
			Exclusions: map[models.PlayerNumber]string{player.Number: player.Initials()},
		})
		if s.evaluateLongPass(*player) {
			return Event{
				Type:            ETPass,
				Team:            team,
				StartingPlayer:  player,
				FinishingPlayer: &receivingPlayer,
			}
		} else {
			return s.turnover(team, player)
		}
	case DecisionShortPass:
		receivingPlayer := team.SearchPlayers(models.PlayerSearchOptions{
			Position:   models.Striker,
			Exclusions: map[models.PlayerNumber]string{player.Number: player.Initials()},
		})
		if s.evaluateShortPass(*player) {
			return Event{
				Type:            ETPass,
				Team:            team,
				StartingPlayer:  player,
				FinishingPlayer: &receivingPlayer,
			}
		} else {
			return s.turnover(team, player)
		}
	case DecisionDribble:
		opposingTeam := s.Simulation.opposingTeam(team)
		if s.evaluateDribble(*player, opposingTeam) {
			return Event{
				Type:            ETDribble,
				Team:            team,
				StartingPlayer:  player,
				FinishingPlayer: player,
			}
		} else {
			return s.turnover(team, player)
		}
	case DecisionCross:
		if s.evaluateCross(*player) {
			nearestTeammate := s.teamMateNearest(player.Position, team.Players)
			return Event{
				Type:            ETCross,
				Team:            team,
				StartingPlayer:  player,
				FinishingPlayer: nearestTeammate,
			}
		} else {
			return s.turnover(team, player)
		}
	case DecisionShoot:
		if s.evaluateShot(team, *player) {
			return Event{
				Type:            ETGoal,
				Team:            team,
				StartingPlayer:  player,
				FinishingPlayer: player,
			}
		} else {
			return s.save(player, team)
		}
	case NoDecision:
		if s.evaluateHold(*player) {
			return Event{
				Type:            ETPossession,
				Team:            team,
				StartingPlayer:  player,
				FinishingPlayer: player,
			}
		} else {
			return s.turnover(team, player)
		}
	default:
		panic(fmt.Errorf("decision '%s' not handled", decision.String()))
	}
}

func (s *SimulationState) turnover(team models.Team, player *models.Player) Event {
	opposingTeam := s.Simulation.opposingTeam(team)
	interceptor := s.opponentNearestTo(player.Position, opposingTeam.Players)
	return Event{
		Type:            ETInterception,
		Team:            opposingTeam,
		StartingPlayer:  player,
		FinishingPlayer: interceptor,
	}
}

func (s *SimulationState) save(player *models.Player, team models.Team) Event {
	opposingTeam := s.Simulation.opposingTeam(team)
	goalKeeper := opposingTeam.SearchPlayers(models.PlayerSearchOptions{
		Position: models.Goalkeeper,
	})
	return Event{
		Type:            ETSave,
		Team:            opposingTeam,
		StartingPlayer:  player,
		FinishingPlayer: &goalKeeper,
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

func (s *SimulationState) evaluateCross(player models.Player) bool {
	return true
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
		proximity := s.positionProximityScore(player.Position, defender.Position, models.OpponentAdjacents)
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

func (s *SimulationState) positionProximityScore(attackerPos, defenderPos models.PlayerPosition, adjacents map[models.PlayerPosition][]models.PlayerPosition) float64 {
	if attackerPos == defenderPos {
		return 1.0
	}

	adjacentPositions, ok := adjacents[attackerPos]
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

func (s *SimulationState) teamMateNearest(attackerPos models.PlayerPosition, players []models.Player) *models.Player {
	var (
		bestScore float64
		closest   *models.Player
	)

	for i, player := range players {
		if player.Position == attackerPos {
			continue
		}

		score := s.positionProximityScore(attackerPos, players[i].Position, models.SimilarPositions)
		if score > bestScore {
			bestScore = score
			closest = &players[i]
		}
	}

	return closest
}

func (s *SimulationState) opponentNearestTo(attackerPos models.PlayerPosition, opponents []models.Player) *models.Player {
	var (
		bestScore float64
		closest   *models.Player
	)

	for i := range opponents {
		score := s.positionProximityScore(attackerPos, opponents[i].Position, models.OpponentAdjacents)
		if score > bestScore {
			bestScore = score
			closest = &opponents[i]
		}
	}

	return closest
}

func (s *SimulationState) evaluateHold(player models.Player) bool {
	if player.Position == models.Goalkeeper {
		return true
	}

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

func (s *SimulationState) evaluateShot(team models.Team, player models.Player) bool {
	power := player.Technical.Shooting.Power
	finishing := player.Technical.Shooting.Finishing
	curve := player.Technical.Shooting.Curve
	composure := player.Composure

	finishingWeight := 0.3
	powerWeight := 0.3
	curveWeight := 0.2
	composureWeight := 0.2

	shotScore := float64(power)*powerWeight +
		float64(finishing)*finishingWeight +
		float64(curve)*curveWeight +
		float64(composure)*composureWeight

	if shotScore < 0 {
		shotScore = 0
	} else if shotScore > 100 {
		shotScore = 100
	}

	normalized := (shotScore - 50) / 10
	shotScore = helpers.Sigmoid(normalized)

	shotScore /= 100.0

	return s.Simulation.RandomFloat() < shotScore
}

func (s *SimulationState) makeDecision(event Event) Decision {
	player := event.FinishingPlayer
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
