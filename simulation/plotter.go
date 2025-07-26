package simulation

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/notoriousbfg/football-game/models"
)

type Pitch struct {
	Match      *Match
	Exclusions map[string]map[models.PlayerNumber]bool
}

func NewPitch(match *Match) *Pitch {
	return &Pitch{
		Match: match,
		Exclusions: map[string]map[models.PlayerNumber]bool{
			match.H.Name: make(map[models.PlayerNumber]bool),
			match.A.Name: make(map[models.PlayerNumber]bool),
		},
	}
}

func (p *Pitch) Draw() {
	homeTeam := p.drawHomeTeam(p.Match.H)
	awayTeam := p.drawAwayTeam(p.Match.A)

	p.drawPitch(homeTeam, awayTeam)
}

func (p *Pitch) drawHomeTeam(team models.Team) []string {
	var rows []string
	switch team.Strategy.Formation {
	case models.FormationFourThreeThree:
		rows = make([]string, 4)
		rows[0] = p.renderRow(team, "three", []models.PlayerPosition{
			models.LeftWinger,
			models.Striker,
			models.RightWinger,
		})
		rows[1] = p.renderRow(team, "three", []models.PlayerPosition{
			models.CentralMidfielder,
			models.CentralAttackingMidfielder,
			models.CentralMidfielder,
		})
		rows[2] = p.renderRow(team, "four", []models.PlayerPosition{
			models.LeftBack,
			models.LeftCentreBack,
			models.RightCentreBack,
			models.RightBack,
		})
		rows[3] = p.renderRow(team, "one", []models.PlayerPosition{
			models.Goalkeeper,
		})
	}
	return rows
}

func (p *Pitch) drawAwayTeam(team models.Team) []string {
	var rows []string
	switch team.Strategy.Formation {
	case models.FormationFourThreeThree:
		rows = make([]string, 4)
		rows[0] = p.renderRow(team, "one", []models.PlayerPosition{
			models.Goalkeeper,
		})
		rows[1] = p.renderRow(team, "four", []models.PlayerPosition{
			models.RightBack,
			models.RightCentreBack,
			models.LeftCentreBack,
			models.LeftBack,
		})
		rows[2] = p.renderRow(team, "three", []models.PlayerPosition{
			models.CentralMidfielder,
			models.CentralAttackingMidfielder | models.CentralMidfielder,
			models.CentralMidfielder,
		})
		rows[3] = p.renderRow(team, "three", []models.PlayerPosition{
			models.RightWinger,
			models.Striker,
			models.LeftWinger,
		})

	}
	return rows
}

func (p *Pitch) drawPitch(home, away []string) {
	templ, err := os.Open("./simulation/templates/pitch.txt")
	if err != nil {
		panic(err)
	}
	defer templ.Close()
	body, err := io.ReadAll(templ)
	if err != nil {
		panic(err)
	}
	bodyStr := string(body)
	pitchParts := strings.Split(bodyStr, "\n")
	result := make([]string, 0)
	result = append(result, pitchParts[0])
	result = append(result, away...)
	result = append(result, pitchParts[1])
	result = append(result, home...)
	result = append(result, pitchParts[2])
	for _, row := range result {
		fmt.Println(row)
	}
}

func (p *Pitch) renderRow(team models.Team, templateName string, positions []models.PlayerPosition) string {
	templ, err := os.Open(fmt.Sprintf("./simulation/templates/%s.txt", templateName))
	if err != nil {
		panic(err)
	}
	defer templ.Close()
	body, err := io.ReadAll(templ)
	if err != nil {
		panic(err)
	}
	bodyStr := string(body)
	exclusions := p.Exclusions[team.Name]
	for i, position := range positions {
		player := team.SearchPlayers(models.PlayerSearchOptions{
			Position:   position,
			Exclusions: exclusions, // prevents adding the same place twice
		})
		initials := player.Initials()
		bodyStr = strings.Replace(
			bodyStr,
			fmt.Sprintf("P%d", i),
			initials,
			1,
		)
		exclusions[player.Number] = true
	}
	return bodyStr
}

type Formation struct {
	Rows []int
}
