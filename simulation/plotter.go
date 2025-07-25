package simulation

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/notoriousbfg/football-game/models"
)

type Pitch struct{}

func (p *Pitch) Draw(match *Match) {
	homeTeam := p.drawHomeTeam(match.H)
	awayTeam := p.drawAwayTeam(match.A)

	p.drawPitch(homeTeam, awayTeam)
}

func (p *Pitch) drawHomeTeam(team models.Team) []string {
	var rows []string
	switch team.Strategy.Formation {
	case models.FormationFourThreeThree:
		rows = make([]string, 4)
		rows[0] = p.renderRow(team, "three", []models.PlayerPosition{models.LeftWinger, models.Striker, models.RightWinger})
		rows[1] = p.renderRow(team, "three", []models.PlayerPosition{models.CentralMidfielder, models.CentralMidfielder, models.CentralMidfielder})
		rows[2] = p.renderRow(team, "four", []models.PlayerPosition{models.LeftBack, models.CentreBack, models.CentreBack, models.RightBack})
		rows[3] = p.renderRow(team, "one", []models.PlayerPosition{models.Goalkeeper})
	}
	return rows
}

func (p *Pitch) drawAwayTeam(team models.Team) []string {
	var rows []string
	switch team.Strategy.Formation {
	case models.FormationFourThreeThree:
		rows = make([]string, 4)
		rows[0] = p.renderRow(team, "one", []models.PlayerPosition{models.Goalkeeper})
		rows[1] = p.renderRow(team, "four", []models.PlayerPosition{models.LeftBack, models.CentreBack, models.CentreBack, models.RightBack})
		rows[2] = p.renderRow(team, "three", []models.PlayerPosition{models.CentralMidfielder, models.CentralMidfielder, models.CentralMidfielder})
		rows[3] = p.renderRow(team, "three", []models.PlayerPosition{models.LeftWinger, models.Striker, models.RightWinger})

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
	foundPlayers := make(map[models.PlayerNumber]bool, 0)
	for i, position := range positions {
		player := team.SearchPlayers(models.PlayerSearchOptions{
			Position:   position,
			Exclusions: foundPlayers, // prevents adding the same place twice
		})
		bodyStr = strings.Replace(
			bodyStr,
			fmt.Sprintf("P%d", i),
			player.Initials(),
			1,
		)
		foundPlayers[player.Number] = true
	}
	return bodyStr
}

type Formation struct {
	Rows []int
}
