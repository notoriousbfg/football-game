package helpers

import "github.com/notoriousbfg/football-game/models"

func IsWinger(pos models.PlayerPosition) bool {
	return pos == models.LeftWinger || pos == models.RightWinger ||
		pos == models.LeftMidfielder || pos == models.RightMidfielder ||
		pos == models.LeftWingBack || pos == models.RightWingBack
}

func IsAttacker(pos models.PlayerPosition) bool {
	return pos == models.Striker || pos == models.CentreForward ||
		pos == models.CentralAttackingMidfielder
}
