package simulation

import (
	"fmt"
	"math/rand"
	"slices"
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
	Pitch             *Pitch
}

func (sim *Simulation) Run() {
	done := make(chan struct{})

	go func() {
		sim.State.Listen()
		close(done)
	}()

	sim.Pitch = NewPitch(&sim.Match)
	sim.Pitch.Draw()

	sim.State.CaptureEvent(
		sim.State.startingEvent(sim.KickoffTeam),
	)

	sim.State.runGame()

	close(sim.State.EventQueue)

	<-done

	sim.State.Outcome = &Outcome{
		HomeScore: sim.State.HomeScore,
		AwayScore: sim.State.AwayScore,
	}
}

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
		FullTime:          false,
		Events:            make([]Event, 0),
		EventQueue:        make(chan Event, 100),
	}

	state.registerTriggers()

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

	return sim
}

// things that change
type SimulationState struct {
	Simulation           *Simulation
	Start                time.Time
	Time                 time.Time
	FirstHalfEnded       bool
	FirstHalfExtraEnded  bool
	SecondHalfStarted    bool
	SecondHalfEnded      bool
	SecondHalfExtraEnded bool
	FullTime             bool
	HomeScore            int
	AwayScore            int
	HomeYellowCards      int
	AwayYellowCards      int
	HomeRedCards         int
	AwayRedCards         int
	HomeMomentum         float64
	AwayMomentum         float64
	HomeTeamAttacking    bool
	AwayTeamAttacking    bool
	Stalemate            bool
	FirstHalfExtraTime   time.Duration // seconds
	SecondHalfExtraTime  time.Duration // seconds
	RandomFloat          func() float64
	Triggers             map[EventType]func(e Event)
	EventQueue           chan Event
	Events               []Event
	Outcome              *Outcome
}

func (s *SimulationState) Listen() {
	for event := range s.EventQueue {
		s.Events = append(s.Events, event)
		if trigger, exists := s.Triggers[event.Type]; exists {
			trigger(event)
		} else {
			s.log(event)
		}
	}
}

func (s *SimulationState) CaptureEvent(e Event) {
	if s.FullTime {
		return
	}

	select {
	case s.EventQueue <- e:
		// sent successfully
	default:
		// could log a warning if needed
	}
}

func (s *SimulationState) LastEvent() Event {
	return s.Events[len(s.Events)-1]
}

func (s *SimulationState) LastEvents(count int) []Event {
	if len(s.Events) <= count {
		return s.Events
	}

	return s.Events[len(s.Events)-count-1 : len(s.Events)-1]
}

func (s *SimulationState) Timestamp() string {
	duration := s.Time.Sub(s.Start)
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", int(duration.Minutes()), seconds)
}

func (s *SimulationState) runGame() {
	firstHalfDuration := 45 * time.Minute
	secondHalfDuration := 45 * time.Minute

	for !s.FullTime {
		elapsed := s.Time.Sub(s.Start)

		switch {
		case !s.FirstHalfEnded && elapsed >= firstHalfDuration:
			s.CaptureEvent(Event{Type: ETEndOfFirstHalf})
			s.FirstHalfEnded = true
		case s.FirstHalfEnded && !s.FirstHalfExtraEnded && elapsed >= firstHalfDuration+s.FirstHalfExtraTime:
			s.CaptureEvent(Event{Type: ETEndOfFirstHalfExtraTime})
			s.FirstHalfExtraEnded = true
			s.SecondHalfStarted = true
		case s.SecondHalfStarted && !s.SecondHalfEnded && elapsed >= firstHalfDuration+secondHalfDuration:
			s.CaptureEvent(Event{Type: ETEndOfSecondHalf})
			s.SecondHalfEnded = true
		case s.SecondHalfEnded && !s.SecondHalfExtraEnded && elapsed >= firstHalfDuration+secondHalfDuration+s.SecondHalfExtraTime:
			s.CaptureEvent(Event{Type: ETEndOfSecondHalfExtraTime})
			s.SecondHalfExtraEnded = true
			s.FullTime = true
		}
	}
}

func (s *SimulationState) registerTriggers() {
	s.Triggers = make(map[EventType]func(e Event))
	s.Triggers[ETEndOfFirstHalf] = func(e Event) {
		s.log(e)
	}
	s.Triggers[ETEndOfFirstHalfExtraTime] = func(e Event) {
		s.log(e)
		// time is reset
		s.Time = s.Start.Add(time.Minute * 45)
	}
	s.Triggers[ETEndOfSecondHalfExtraTime] = func(e Event) {
		s.log(e)
		if s.HomeScore < s.AwayScore {
			s.HomeMomentum += 0.5
		} else {
			s.AwayMomentum += 0.5
		}
	}
	s.Triggers[ETPass] = func(e Event) {
		s.log(e)
		s.addTime(time.Second * 3)
		s.CaptureEvent(
			s.action(e),
		)
	}
	s.Triggers[ETCross] = func(e Event) {
		s.log(e)
		s.addTime(time.Second * 3)
		s.CaptureEvent(
			s.action(e),
		)
	}
	s.Triggers[ETDribble] = func(e Event) {
		s.log(e)
		s.addTime(time.Second * 5)
		s.CaptureEvent(
			s.action(e),
		)
	}
	s.Triggers[ETGoal] = func(e Event) {
		s.log(e)
		s.addTime(time.Minute * 2)
		if s.isHome(e.Team) {
			s.HomeScore++
			s.HomeMomentum += 0.1
			s.AwayMomentum -= 0.2
		} else {
			s.AwayScore++
			s.AwayMomentum += 0.1
			s.HomeMomentum -= 0.2
		}
		s.CaptureEvent(
			s.goal(e),
		)
	}
	s.Triggers[ETReset] = func(e Event) {
		s.log(e)
		s.addTime(time.Second * 3)
		s.CaptureEvent(
			s.reset(e),
		)
	}
	s.Triggers[ETInterception] = func(e Event) {
		s.log(e)
		s.addTime(time.Second * 3)
		s.CaptureEvent(
			s.action(e),
		)
	}
	s.Triggers[ETPossession] = func(e Event) {
		s.log(e)
		s.addTime(time.Second * 3)
		s.CaptureEvent(
			s.action(e),
		)
	}
	s.Triggers[ETYellowCard] = func(e Event) {
		s.log(e)
		s.CaptureEvent(
			s.freeKick(e),
		)
		duration := time.Second * 3
		s.addTime(duration)
		s.addExtraTime(duration)
		if s.isHome(e.Team) {
			s.HomeYellowCards++
		} else {
			s.AwayYellowCards++
		}
	}
	s.Triggers[ETRedCard] = func(e Event) {
		s.CaptureEvent(
			s.freeKick(e),
		)
		duration := time.Second * 20
		s.addTime(duration)
		s.addExtraTime(duration)
		if s.isHome(e.Team) {
			s.HomeRedCards++
		} else {
			s.HomeRedCards++
		}
	}
	s.Triggers[ETSave] = func(e Event) {
		s.log(e)
		s.addTime(time.Second * 3)
		s.addExtraTime(time.Second * 1)
		s.CaptureEvent(
			s.goalKeeperKick(e),
		)
	}
}

func (s *SimulationState) addTime(d time.Duration) {
	rand := s.Simulation.RandomFloat() * 2.0 // random number between 0 and 2
	s.Time = s.Time.Add(time.Duration(d.Seconds()*rand) * time.Second)
}

func (s *SimulationState) addExtraTime(d time.Duration) {
	if !s.FirstHalfEnded {
		s.FirstHalfExtraTime += d
	}
	if s.SecondHalfStarted && !s.SecondHalfEnded {
		s.SecondHalfExtraTime += d
	}
}

func (s *SimulationState) log(e Event) {
	switch e.Type {
	case ETPass:
		fmt.Printf("(%s) %s passes to %s\n", s.Timestamp(), e.StartingPlayer.Name, e.FinishingPlayer.Name)
	case ETGoal:
		fmt.Printf("(%s) %s shoots and scores!\n", s.Timestamp(), e.FinishingPlayer.Name)
	case ETReset:
		fmt.Printf("(%s) The game restarts after the goal\n", s.Timestamp())
	case ETCross:
		fmt.Printf("(%s) %s crosses to %s\n", s.Timestamp(), e.StartingPlayer.Name, e.FinishingPlayer.Name)
	case ETDribble:
		fmt.Printf("(%s) %s is dribbling with the ball\n", s.Timestamp(), e.StartingPlayer.Name)
	case ETInterception:
		fmt.Printf("(%s) %s loses the ball to %s\n", s.Timestamp(), e.StartingPlayer.Name, e.FinishingPlayer.Name)
	case ETPossession:
		fmt.Printf("(%s) %s has the ball\n", s.Timestamp(), e.FinishingPlayer.Name)
	case ETYellowCard:
		fmt.Printf("(%s) %s is given a yellow card for a foul\n", s.Timestamp(), e.FinishingPlayer.Name)
	case ETRedCard:
		fmt.Printf("(%s) %s is shown a red card for a bad foul\n", s.Timestamp(), e.FinishingPlayer.Name)
	case ETSave:
		fmt.Printf("(%s) %s took a shot but it was saved by %s\n", s.Timestamp(), e.StartingPlayer.Name, e.FinishingPlayer.Name)
	case ETEndOfFirstHalf:
		fmt.Printf("(%s) %d minutes extra time added in the first half\n", s.Timestamp(), int(s.FirstHalfExtraTime.Minutes()))
	case ETEndOfFirstHalfExtraTime:
		fmt.Printf("(%s) The whistle blows for the end of the first half\n", s.Timestamp())
	case ETEndOfSecondHalf:
		fmt.Printf("(%s) %d minutes extra time added in the second half\n", s.Timestamp(), int(s.SecondHalfExtraTime.Minutes()))
	case ETEndOfSecondHalfExtraTime:
		fmt.Printf("(%s) The full time whistle blows\n", s.Timestamp())
		// case ETRestart:
		// 	// Assuming this special case doesn't need a full Event object
		// 	fmt.Printf("(%s) The game restarts.\n", s.Timestamp())
	}
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
		Positions: []models.PlayerPosition{models.Striker},
	})
	receivingPlayer := team.SearchPlayers(models.PlayerSearchOptions{
		Positions: []models.PlayerPosition{models.CentralMidfielder},
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
	coinFlip := s.Simulation.RandomFloat()
	var receivingPlayer models.Player
	if coinFlip < 0.5 {
		receivingPlayer = e.Team.SearchPlayers(models.PlayerSearchOptions{
			Positions: []models.PlayerPosition{models.Striker, models.LeftWinger, models.LeftMidfielder, models.RightWinger, models.RightMidfielder, models.CentreForward, models.CentralAttackingMidfielder},
		})
	} else {
		receivingPlayer = e.Team.SearchPlayers(models.PlayerSearchOptions{
			Positions: []models.PlayerPosition{models.LeftBack, models.LeftCentreBack, models.RightCentreBack, models.RightBack},
		})
	}
	return Event{
		Type:            ETPass,
		Team:            e.Team,
		StartingPlayer:  e.FinishingPlayer,
		FinishingPlayer: &receivingPlayer,
	}
}

func (s *SimulationState) goal(e Event) Event {
	opposingTeam := s.Simulation.opposingTeam(e.Team)
	return Event{
		Type: ETReset,
		Team: opposingTeam,
	}
}

func (s *SimulationState) reset(e Event) Event {
	startingPlayer := e.Team.SearchPlayers(models.PlayerSearchOptions{
		Positions: []models.PlayerPosition{models.Striker},
	})
	receivingPlayer := e.Team.SearchPlayers(models.PlayerSearchOptions{
		Positions: []models.PlayerPosition{models.CentralMidfielder},
	})
	return Event{
		Type:            ETPass,
		Team:            e.Team,
		StartingPlayer:  &startingPlayer,
		FinishingPlayer: &receivingPlayer,
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
		if s.evaluateLongPass(team, *player) {
			receivingPlayer := team.SearchPlayers(models.PlayerSearchOptions{
				Positions:  []models.PlayerPosition{models.Striker, models.RightWinger, models.LeftWinger},
				Exclusions: map[models.PlayerNumber]string{player.Number: player.Initials()},
			})
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
		nearestPlayer := s.teamMateNearest(player.Position, team.Players)
		receivingPlayer := team.SearchPlayers(models.PlayerSearchOptions{
			Positions:  []models.PlayerPosition{nearestPlayer.Position},
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
			receivingPlayer := team.SearchPlayers(models.PlayerSearchOptions{
				Positions:  []models.PlayerPosition{nearestTeammate.Position, models.LeftWinger, models.RightWinger, models.Striker, models.CentralAttackingMidfielder, models.CentreForward},
				Exclusions: map[models.PlayerNumber]string{player.Number: player.Initials()},
			})
			return Event{
				Type:            ETCross,
				Team:            team,
				StartingPlayer:  player,
				FinishingPlayer: &receivingPlayer,
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
		Positions: []models.PlayerPosition{models.Goalkeeper},
	})
	return Event{
		Type:            ETSave,
		Team:            opposingTeam,
		StartingPlayer:  player,
		FinishingPlayer: &goalKeeper,
	}
}

func (s *SimulationState) evaluateLongPass(team models.Team, player models.Player) bool {
	skill := player.Technical.Passing.LongPass
	vision := player.TacticalIntelligence.Vision.Passing
	agility := player.Fitness.Agility

	var momentum float64
	if s.isHome(team) {
		momentum = s.HomeMomentum
	} else {
		momentum = s.AwayMomentum
	}

	visionWeight := 0.4
	executionWeight := 0.4
	agilityWeight := 0.3
	momentumWeight := 0.1

	successChance := float64(vision)*visionWeight +
		float64(skill)*executionWeight +
		float64(agility)*agilityWeight +
		float64(momentum)*momentumWeight

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

	// make subsequent dribbles harder
	if s.LastEvent().Type == ETDribble {
		dribbleScore *= 0.5
	}

	// penalize for nearby opponents
	// for _, defender := range opposingTeam.Players {
	// 	proximity := s.positionProximityScore(player.Position, defender.Position, models.OpponentAdjacents)
	// 	if proximity > 0 {
	// 		dribbleScore -= float64(defender.Technical.Defending.Interceptions) * proximity
	// 	}
	// }

	// adjust for opposing team’s overall defensive acumen
	// totalDef := 0.0
	// for _, opp := range opposingTeam.Players {
	// 	totalDef += (float64(opp.Technical.Defending.Interceptions) * 0.5) + (float64(opp.Technical.Defending.Blocking) * 0.5)
	// }
	// avgTeamDef := totalDef / float64(len(opposingTeam.Players))
	// dribbleScore -= avgTeamDef * 0.1 // global pressure factor

	// clamp the score
	if dribbleScore < 0 {
		dribbleScore = 0
	} else if dribbleScore > 100 {
		dribbleScore = 100
	}

	dribbleScore /= 100.0

	return s.Simulation.RandomFloat() < dribbleScore
}

func (s *SimulationState) positionProximityScore(playerPos, adjacentPos models.PlayerPosition, adjacents map[models.PlayerPosition][]models.PlayerPosition) float64 {
	if playerPos == adjacentPos {
		return 1.0
	}

	adjacentPositions, ok := adjacents[playerPos]
	if !ok {
		return 0.0
	}

	for _, pos := range adjacentPositions {
		if pos == adjacentPos {
			return 0.75
		}
	}

	for _, adjPos := range adjacentPositions {
		secondDegreePositions, ok := adjacents[adjPos]
		if !ok {
			continue
		}
		if slices.Contains(secondDegreePositions, adjacentPos) {
			return 0.5
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

		score := s.positionProximityScore(attackerPos, players[i].Position, models.TeammateAdjacents)
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
	power := float64(player.Technical.Shooting.Power)
	finishing := float64(player.Technical.Shooting.Finishing)
	curve := float64(player.Technical.Shooting.Curve)
	composure := float64(player.Composure)

	shotScore := power*0.3 + finishing*0.3 + curve*0.2 + composure*0.2

	if shotScore < 0 {
		shotScore = 0
	} else if shotScore > 100 {
		shotScore = 100
	}

	opponentTeam := s.Simulation.opposingTeam(team)
	opponentKeeper := opponentTeam.SearchPlayers(models.PlayerSearchOptions{
		Positions:  []models.PlayerPosition{models.Goalkeeper},
		Exclusions: map[models.PlayerNumber]string{player.Number: player.Initials()},
	})
	goalkeeperFactor := float64(opponentKeeper.Technical.Goalkeeping.Reflexes)*0.4 +
		float64(opponentKeeper.Technical.Goalkeeping.Positioning)*0.4 +
		float64(opponentKeeper.Technical.Goalkeeping.Reactions)*0.2
	goalkeeperDifficulty := helpers.Sigmoid((goalkeeperFactor - 50) / 10)

	shotScore *= 1 - goalkeeperDifficulty

	normalized := (shotScore - 50) / 10
	probability := helpers.Sigmoid(normalized)

	return s.Simulation.RandomFloat() < probability
}

func (s *SimulationState) makeDecision(event Event) Decision {
	player := event.FinishingPlayer
	if player == nil {
		return NoDecision
	}

	dribbling := float64(player.Technical.Dribbling.Dribbling)
	shooting := float64(player.Technical.Shooting.Finishing)*0.5 +
		float64(player.TacticalIntelligence.Vision.Shooting)*0.5
	crossing := float64(player.Technical.Passing.Cross)
	position := player.Position

	rng := s.Simulation.RandomFloat()

	// shooting — more likely if forward or attacking midfielder
	if helpers.IsAttacker(position) && shooting > 70 && rng < shooting/200.0 {
		return DecisionShoot
	}

	// shoot if attacker and passing in box
	var involvedInRecentEvents bool
	for _, event := range s.LastEvents(5) {
		if event.FinishingPlayer == nil {
			continue
		}
		if event.FinishingPlayer.Position == position || event.StartingPlayer.Position == position {
			involvedInRecentEvents = true
		}
	}
	if (helpers.IsAttacker(position) || helpers.IsWinger(position)) && involvedInRecentEvents {
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
	return s.decidePassType(*player)
}

func (s *SimulationState) decidePassType(player models.Player) Decision {
	longPassProbability := float64(player.TacticalIntelligence.Vision.Passing) * 0.8 // up to 80% chance

	if s.Simulation.RandomFloat() < longPassProbability {
		return DecisionLongPass
	}
	return DecisionShortPass
}
