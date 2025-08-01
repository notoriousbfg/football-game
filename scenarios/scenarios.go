package scenarios

import "github.com/notoriousbfg/football-game/models"

func HomeTeam() models.Team {
	return models.Team{
		Name:      "AFC Bournemouth",
		Morale:    80,
		Fitness:   76,
		Chemistry: 85,
		Strategy: models.Strategy{
			Tactic:    models.TacticPressing,
			Formation: models.FormationFourFourTwo,
			PlayStyle: models.PlayStyleCreative,
			PlayerInstructions: map[models.PlayerNumber]models.Instruction{
				1:  {Position: models.PositionCenter}, // GK
				2:  {Position: models.PositionCenter}, // RB
				3:  {Position: models.PositionCenter}, // LB
				4:  {Position: models.PositionCenter}, // CB
				5:  {Position: models.PositionCenter}, // CB
				6:  {Position: models.PositionCenter}, // CM
				7:  {Position: models.PositionWing},   // RW
				8:  {Position: models.PositionCenter}, // CM
				9:  {Position: models.PositionCenter}, // ST
				10: {Position: models.PositionCenter}, // CAM
				11: {Position: models.PositionWing},   // LW
			},
		},
		Training: models.Training{
			Focus: models.Passing,
		},
		Players: []models.Player{
			{
				Name:     "Kepa",
				Number:   1,
				Position: models.Goalkeeper,
				Form:     75, Adaptability: 72, Composure: 78,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 50, Acceleration: 48},
					Passing:  models.PassingSkill{ShortPass: 60, LongPass: 62, Cross: 55, Lob: 60, ThroughBall: 58, Chip: 55},
					Shooting: models.ShootingSkill{Power: 50, Curve: 48, Finishing: 45, Spin: 40},
					Defending: models.DefendingSkill{
						Jumping:       75,
						Interceptions: 60,
						Heading:       models.HeadingSkill{Accuracy: 62, Power: 65},
						Blocking:      68,
					},
					FreeKicks: 45, Penalties: 40,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 80,
					Vision:      models.TacticalVision{Passing: 65, Shooting: 40, Defence: 85},
				},
				Stamina: models.Stamina{Stamina: 65},
				Fitness: models.Fitness{Strength: 70, Agility: 66, InjuryTolerance: 80, InjuryResistance: 78},
			},
			{
				Name:     "Adam Smith",
				Number:   2,
				Position: models.RightBack,
				Form:     72, Adaptability: 75, Composure: 70,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 78, Acceleration: 75},
					Passing:  models.PassingSkill{ShortPass: 68, LongPass: 65, Cross: 70, Lob: 60, ThroughBall: 55, Chip: 50},
					Shooting: models.ShootingSkill{Power: 60, Curve: 58, Finishing: 55, Spin: 50},
					Defending: models.DefendingSkill{
						Jumping:       70,
						Interceptions: 75,
						Heading:       models.HeadingSkill{Accuracy: 60, Power: 65},
						Blocking:      72,
					},
					FreeKicks: 45, Penalties: 40,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 75,
					Vision:      models.TacticalVision{Passing: 68, Shooting: 50, Defence: 75},
				},
				Stamina: models.Stamina{Stamina: 78},
				Fitness: models.Fitness{Strength: 70, Agility: 72, InjuryTolerance: 75, InjuryResistance: 77},
			},
			{
				Name:     "Ilya Zabarnyi",
				Number:   5,
				Position: models.LeftCentreBack,
				Form:     73, Adaptability: 70, Composure: 75,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 68, Acceleration: 65},
					Passing:  models.PassingSkill{ShortPass: 70, LongPass: 72, Cross: 55, Lob: 58, ThroughBall: 60, Chip: 50},
					Shooting: models.ShootingSkill{Power: 60, Curve: 50, Finishing: 48, Spin: 45},
					Defending: models.DefendingSkill{
						Jumping:       75,
						Interceptions: 78,
						Heading:       models.HeadingSkill{Accuracy: 78, Power: 80},
						Blocking:      72,
					},
					FreeKicks: 40, Penalties: 30,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 78,
					Vision:      models.TacticalVision{Passing: 68, Shooting: 45, Defence: 80},
				},
				Stamina: models.Stamina{Stamina: 72},
				Fitness: models.Fitness{Strength: 78, Agility: 68, InjuryTolerance: 75, InjuryResistance: 80},
			},
			{
				Name:         "Dean Huijsen",
				Position:     models.RightCentreBack,
				Number:       15,
				Form:         78,
				Adaptability: 80,
				Composure:    82,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 70, Acceleration: 68},
					Passing:  models.PassingSkill{ShortPass: 78, LongPass: 80, Cross: 60, Lob: 75, ThroughBall: 70, Chip: 65},
					Shooting: models.ShootingSkill{Power: 65, Curve: 60, Finishing: 58, Spin: 60},
					Defending: models.DefendingSkill{
						Jumping:       85,
						Interceptions: 82,
						Heading:       models.HeadingSkill{Accuracy: 84, Power: 80},
						Blocking:      80,
					},
					FreeKicks: 50,
					Penalties: 60,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 80,
					Vision: models.TacticalVision{
						Passing:  78,
						Shooting: 55,
						Defence:  80,
					},
				},
				Stamina: models.Stamina{
					Stamina: 76,
				},
				Fitness: models.Fitness{
					Strength:         85,
					Agility:          72,
					InjuryTolerance:  78,
					InjuryResistance: 75,
				},
			},
			{
				Name:         "Justin Kluivert",
				Position:     models.RightWinger,
				Number:       19,
				Form:         80,
				Adaptability: 84,
				Composure:    78,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 90, Acceleration: 92},
					Passing:  models.PassingSkill{ShortPass: 78, LongPass: 70, Cross: 80, Lob: 72, ThroughBall: 75, Chip: 70},
					Shooting: models.ShootingSkill{Power: 75, Curve: 80, Finishing: 78, Spin: 77},
					Defending: models.DefendingSkill{
						Jumping:       65,
						Interceptions: 60,
						Heading:       models.HeadingSkill{Accuracy: 60, Power: 62},
						Blocking:      55,
					},
					FreeKicks: 68,
					Penalties: 70,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 82,
					Vision: models.TacticalVision{
						Passing:  78,
						Shooting: 80,
						Defence:  60,
					},
				},
				Stamina: models.Stamina{
					Stamina: 84,
				},
				Fitness: models.Fitness{
					Strength:         70,
					Agility:          90,
					InjuryTolerance:  75,
					InjuryResistance: 78,
				},
			},
			{
				Name:     "Milos Kerkez",
				Number:   3,
				Position: models.LeftBack,
				Form:     70, Adaptability: 68, Composure: 72,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 75, Acceleration: 72},
					Passing:  models.PassingSkill{ShortPass: 70, LongPass: 68, Cross: 72, Lob: 65, ThroughBall: 60, Chip: 55},
					Shooting: models.ShootingSkill{Power: 58, Curve: 55, Finishing: 52, Spin: 50},
					Defending: models.DefendingSkill{
						Jumping:       68,
						Interceptions: 70,
						Heading:       models.HeadingSkill{Accuracy: 65, Power: 65},
						Blocking:      68,
					},
					FreeKicks: 45, Penalties: 40,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 70,
					Vision:      models.TacticalVision{Passing: 68, Shooting: 50, Defence: 70},
				},
				Stamina: models.Stamina{Stamina: 75},
				Fitness: models.Fitness{Strength: 68, Agility: 70, InjuryTolerance: 75, InjuryResistance: 77},
			},
			{
				Name:     "Philip Billing",
				Number:   7,
				Position: models.CentralMidfielder,
				Form:     75, Adaptability: 73, Composure: 72,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 72, Acceleration: 70},
					Passing:  models.PassingSkill{ShortPass: 75, LongPass: 70, Cross: 65, Lob: 60, ThroughBall: 62, Chip: 58},
					Shooting: models.ShootingSkill{Power: 78, Curve: 70, Finishing: 68, Spin: 60},
					Defending: models.DefendingSkill{
						Jumping:       68,
						Interceptions: 65,
						Heading:       models.HeadingSkill{Accuracy: 60, Power: 65},
						Blocking:      62,
					},
					FreeKicks: 70, Penalties: 60,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 75,
					Vision:      models.TacticalVision{Passing: 70, Shooting: 70, Defence: 65},
				},
				Stamina: models.Stamina{Stamina: 75},
				Fitness: models.Fitness{Strength: 72, Agility: 70, InjuryTolerance: 75, InjuryResistance: 73},
			},
			{
				Name:     "Lewis Cook",
				Number:   8,
				Position: models.CentralMidfielder,
				Form:     74, Adaptability: 70, Composure: 75,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 68, Acceleration: 65},
					Passing:  models.PassingSkill{ShortPass: 78, LongPass: 72, Cross: 65, Lob: 60, ThroughBall: 68, Chip: 62},
					Shooting: models.ShootingSkill{Power: 65, Curve: 60, Finishing: 60, Spin: 55},
					Defending: models.DefendingSkill{
						Jumping:       65,
						Interceptions: 68,
						Heading:       models.HeadingSkill{Accuracy: 58, Power: 60},
						Blocking:      62,
					},
					FreeKicks: 60, Penalties: 55,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 72,
					Vision:      models.TacticalVision{Passing: 75, Shooting: 62, Defence: 68},
				},
				Stamina: models.Stamina{Stamina: 72},
				Fitness: models.Fitness{Strength: 68, Agility: 66, InjuryTolerance: 75, InjuryResistance: 72},
			},
			{
				Name:     "Evanilson",
				Number:   9,
				Position: models.Striker,
				Form:     76, Adaptability: 73, Composure: 74,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 72, Acceleration: 70},
					Passing:  models.PassingSkill{ShortPass: 68, LongPass: 65, Cross: 55, Lob: 50, ThroughBall: 55, Chip: 50},
					Shooting: models.ShootingSkill{Power: 80, Curve: 70, Finishing: 78, Spin: 65},
					Defending: models.DefendingSkill{
						Jumping:       70,
						Interceptions: 50,
						Heading:       models.HeadingSkill{Accuracy: 65, Power: 70},
						Blocking:      55,
					},
					FreeKicks: 40, Penalties: 50,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 78,
					Vision:      models.TacticalVision{Passing: 60, Shooting: 75, Defence: 50},
				},
				Stamina: models.Stamina{Stamina: 70},
				Fitness: models.Fitness{Strength: 75, Agility: 70, InjuryTolerance: 75, InjuryResistance: 74},
			},
			{
				Name:     "Ryan Christie",
				Number:   10,
				Position: models.CentralAttackingMidfielder,
				Form:     74, Adaptability: 72, Composure: 73,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 70, Acceleration: 68},
					Passing:  models.PassingSkill{ShortPass: 78, LongPass: 70, Cross: 68, Lob: 65, ThroughBall: 70, Chip: 65},
					Shooting: models.ShootingSkill{Power: 72, Curve: 70, Finishing: 70, Spin: 65},
					Defending: models.DefendingSkill{
						Jumping:       60,
						Interceptions: 62,
						Heading:       models.HeadingSkill{Accuracy: 58, Power: 60},
						Blocking:      60,
					},
					FreeKicks: 70, Penalties: 68,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 75,
					Vision:      models.TacticalVision{Passing: 78, Shooting: 70, Defence: 60},
				},
				Stamina: models.Stamina{Stamina: 70},
				Fitness: models.Fitness{Strength: 68, Agility: 70, InjuryTolerance: 72, InjuryResistance: 70},
			},
			{
				Name:     "Dango Ouattara",
				Number:   11,
				Position: models.LeftWinger,
				Form:     73, Adaptability: 70, Composure: 70,
				Technical: models.TechnicalSkill{
					Speed:    models.SpeedSkill{Speed: 78, Acceleration: 76},
					Passing:  models.PassingSkill{ShortPass: 70, LongPass: 68, Cross: 75, Lob: 70, ThroughBall: 68, Chip: 65},
					Shooting: models.ShootingSkill{Power: 68, Curve: 65, Finishing: 65, Spin: 60},
					Defending: models.DefendingSkill{
						Jumping:       55,
						Interceptions: 55,
						Heading:       models.HeadingSkill{Accuracy: 50, Power: 50},
						Blocking:      55,
					},
					FreeKicks: 55, Penalties: 50,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 70,
					Vision:      models.TacticalVision{Passing: 70, Shooting: 65, Defence: 55},
				},
				Stamina: models.Stamina{Stamina: 75},
				Fitness: models.Fitness{Strength: 65, Agility: 72, InjuryTolerance: 70, InjuryResistance: 72},
			},
		},
	}
}

func AwayTeam() models.Team {
	return models.Team{
		Name:      "Arsenal",
		Morale:    85,
		Fitness:   90,
		Chemistry: 88,
		Strategy: models.Strategy{
			Tactic:    models.TacticPressing,
			Formation: models.FormationFourThreeThree,
			PlayerInstructions: map[models.PlayerNumber]models.Instruction{
				1:  {Position: models.PositionCenter},
				4:  {Position: models.PositionCenter},
				6:  {Position: models.PositionCenter},
				35: {Position: models.PositionWing},
				18: {Position: models.PositionWing},
				5:  {Position: models.PositionCenter},
				8:  {Position: models.PositionCenter},
				41: {Position: models.PositionCenter},
				7:  {Position: models.PositionWing},
				11: {Position: models.PositionWing},
				9:  {Position: models.PositionCenter},
			},
			PlayStyle: models.PlayStyleDriven,
		},
		Training: models.Training{
			Focus: models.Passing,
		},
		Players: []models.Player{
			{
				Name:         "Aaron Ramsdale",
				Position:     models.Goalkeeper,
				Number:       1,
				Form:         82,
				Adaptability: 80,
				Composure:    84,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 55, Acceleration: 60},
					Passing:   models.PassingSkill{ShortPass: 75, LongPass: 78, Cross: 50, Lob: 65, ThroughBall: 60, Chip: 55},
					Shooting:  models.ShootingSkill{Power: 50, Curve: 45, Finishing: 30, Spin: 40},
					Defending: models.DefendingSkill{Jumping: 85, Interceptions: 60, Heading: models.HeadingSkill{Accuracy: 45, Power: 60}, Blocking: 90},
					FreeKicks: 30,
					Penalties: 35,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 85,
					Vision:      models.TacticalVision{Passing: 75, Shooting: 40, Defence: 90},
				},
				Stamina: models.Stamina{Stamina: 70},
				Fitness: models.Fitness{Strength: 78, Agility: 85, InjuryTolerance: 85, InjuryResistance: 88},
			},
			{
				Name:         "Ben White",
				Position:     models.RightBack,
				Number:       4,
				Form:         84,
				Adaptability: 85,
				Composure:    82,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 82, Acceleration: 80},
					Passing:   models.PassingSkill{ShortPass: 80, LongPass: 78, Cross: 75, Lob: 70, ThroughBall: 72, Chip: 68},
					Shooting:  models.ShootingSkill{Power: 65, Curve: 60, Finishing: 55, Spin: 58},
					Defending: models.DefendingSkill{Jumping: 78, Interceptions: 85, Heading: models.HeadingSkill{Accuracy: 75, Power: 72}, Blocking: 80},
					FreeKicks: 45,
					Penalties: 50,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 82,
					Vision:      models.TacticalVision{Passing: 78, Shooting: 50, Defence: 85},
				},
				Stamina: models.Stamina{Stamina: 85},
				Fitness: models.Fitness{Strength: 80, Agility: 82, InjuryTolerance: 85, InjuryResistance: 86},
			},
			{
				Name:         "William Saliba",
				Position:     models.LeftCentreBack,
				Number:       6,
				Form:         88,
				Adaptability: 88,
				Composure:    90,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 80, Acceleration: 78},
					Passing:   models.PassingSkill{ShortPass: 82, LongPass: 80, Cross: 60, Lob: 70, ThroughBall: 68, Chip: 65},
					Shooting:  models.ShootingSkill{Power: 70, Curve: 60, Finishing: 50, Spin: 55},
					Defending: models.DefendingSkill{Jumping: 85, Interceptions: 90, Heading: models.HeadingSkill{Accuracy: 88, Power: 84}, Blocking: 92},
					FreeKicks: 40,
					Penalties: 55,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 90,
					Vision:      models.TacticalVision{Passing: 80, Shooting: 45, Defence: 92},
				},
				Stamina: models.Stamina{Stamina: 88},
				Fitness: models.Fitness{Strength: 85, Agility: 80, InjuryTolerance: 90, InjuryResistance: 90},
			},
			{
				Name:         "Oleksandr Zinchenko",
				Position:     models.LeftBack,
				Number:       35,
				Form:         80,
				Adaptability: 84,
				Composure:    78,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 78, Acceleration: 80},
					Passing:   models.PassingSkill{ShortPass: 85, LongPass: 84, Cross: 80, Lob: 75, ThroughBall: 82, Chip: 70},
					Shooting:  models.ShootingSkill{Power: 65, Curve: 72, Finishing: 60, Spin: 68},
					Defending: models.DefendingSkill{Jumping: 70, Interceptions: 75, Heading: models.HeadingSkill{Accuracy: 60, Power: 60}, Blocking: 70},
					FreeKicks: 60,
					Penalties: 55,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 80,
					Vision:      models.TacticalVision{Passing: 85, Shooting: 65, Defence: 78},
				},
				Stamina: models.Stamina{Stamina: 85},
				Fitness: models.Fitness{Strength: 74, Agility: 85, InjuryTolerance: 80, InjuryResistance: 82},
			},
			{
				Name:         "Takehiro Tomiyasu",
				Position:     models.RightCentreBack,
				Number:       18,
				Form:         82,
				Adaptability: 80,
				Composure:    83,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 76, Acceleration: 78},
					Passing:   models.PassingSkill{ShortPass: 78, LongPass: 75, Cross: 70, Lob: 70, ThroughBall: 65, Chip: 65},
					Shooting:  models.ShootingSkill{Power: 60, Curve: 58, Finishing: 50, Spin: 55},
					Defending: models.DefendingSkill{Jumping: 80, Interceptions: 85, Heading: models.HeadingSkill{Accuracy: 80, Power: 78}, Blocking: 85},
					FreeKicks: 45,
					Penalties: 50,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 84,
					Vision:      models.TacticalVision{Passing: 75, Shooting: 50, Defence: 85},
				},
				Stamina: models.Stamina{Stamina: 82},
				Fitness: models.Fitness{Strength: 82, Agility: 80, InjuryTolerance: 80, InjuryResistance: 80},
			},
			{
				Name:         "Thomas Partey",
				Position:     models.CentralDefensiveMidfielder,
				Number:       5,
				Form:         83,
				Adaptability: 82,
				Composure:    85,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 78, Acceleration: 75},
					Passing:   models.PassingSkill{ShortPass: 85, LongPass: 82, Cross: 70, Lob: 75, ThroughBall: 78, Chip: 70},
					Shooting:  models.ShootingSkill{Power: 72, Curve: 65, Finishing: 60, Spin: 60},
					Defending: models.DefendingSkill{Jumping: 75, Interceptions: 85, Heading: models.HeadingSkill{Accuracy: 78, Power: 75}, Blocking: 85},
					FreeKicks: 55,
					Penalties: 60,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 86,
					Vision:      models.TacticalVision{Passing: 82, Shooting: 68, Defence: 88},
				},
				Stamina: models.Stamina{Stamina: 86},
				Fitness: models.Fitness{Strength: 85, Agility: 80, InjuryTolerance: 80, InjuryResistance: 78},
			},
			{
				Name:         "Martin Ødegaard",
				Position:     models.CentralAttackingMidfielder,
				Number:       8,
				Form:         90,
				Adaptability: 88,
				Composure:    90,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 82, Acceleration: 85},
					Passing:   models.PassingSkill{ShortPass: 92, LongPass: 88, Cross: 80, Lob: 85, ThroughBall: 95, Chip: 82},
					Shooting:  models.ShootingSkill{Power: 75, Curve: 90, Finishing: 85, Spin: 85},
					Defending: models.DefendingSkill{Jumping: 65, Interceptions: 72, Heading: models.HeadingSkill{Accuracy: 60, Power: 58}, Blocking: 60},
					FreeKicks: 88,
					Penalties: 85,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 92,
					Vision:      models.TacticalVision{Passing: 95, Shooting: 90, Defence: 75},
				},
				Stamina: models.Stamina{Stamina: 85},
				Fitness: models.Fitness{Strength: 72, Agility: 90, InjuryTolerance: 80, InjuryResistance: 82},
			},
			{
				Name:         "Declan Rice",
				Position:     models.CentralMidfielder,
				Number:       41,
				Form:         88,
				Adaptability: 86,
				Composure:    86,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 80, Acceleration: 80},
					Passing:   models.PassingSkill{ShortPass: 85, LongPass: 85, Cross: 75, Lob: 80, ThroughBall: 82, Chip: 78},
					Shooting:  models.ShootingSkill{Power: 80, Curve: 70, Finishing: 72, Spin: 70},
					Defending: models.DefendingSkill{Jumping: 80, Interceptions: 88, Heading: models.HeadingSkill{Accuracy: 82, Power: 80}, Blocking: 85},
					FreeKicks: 60,
					Penalties: 70,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 88,
					Vision:      models.TacticalVision{Passing: 85, Shooting: 75, Defence: 88},
				},
				Stamina: models.Stamina{Stamina: 90},
				Fitness: models.Fitness{Strength: 88, Agility: 85, InjuryTolerance: 88, InjuryResistance: 88},
			},
			{
				Name:         "Bukayo Saka",
				Position:     models.RightWinger,
				Number:       7,
				Form:         92,
				Adaptability: 90,
				Composure:    88,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 90, Acceleration: 92},
					Passing:   models.PassingSkill{ShortPass: 88, LongPass: 82, Cross: 90, Lob: 80, ThroughBall: 88, Chip: 85},
					Shooting:  models.ShootingSkill{Power: 85, Curve: 90, Finishing: 90, Spin: 88},
					Defending: models.DefendingSkill{Jumping: 65, Interceptions: 70, Heading: models.HeadingSkill{Accuracy: 60, Power: 58}, Blocking: 65},
					FreeKicks: 80,
					Penalties: 85,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 90,
					Vision:      models.TacticalVision{Passing: 90, Shooting: 88, Defence: 72},
				},
				Stamina: models.Stamina{Stamina: 92},
				Fitness: models.Fitness{Strength: 78, Agility: 95, InjuryTolerance: 85, InjuryResistance: 88},
			},
			{
				Name:         "Gabriel Martinelli",
				Position:     models.LeftWinger,
				Number:       11,
				Form:         88,
				Adaptability: 85,
				Composure:    84,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 92, Acceleration: 94},
					Passing:   models.PassingSkill{ShortPass: 82, LongPass: 78, Cross: 88, Lob: 80, ThroughBall: 80, Chip: 78},
					Shooting:  models.ShootingSkill{Power: 82, Curve: 85, Finishing: 86, Spin: 84},
					Defending: models.DefendingSkill{Jumping: 70, Interceptions: 72, Heading: models.HeadingSkill{Accuracy: 68, Power: 70}, Blocking: 72},
					FreeKicks: 70,
					Penalties: 75,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 88,
					Vision:      models.TacticalVision{Passing: 85, Shooting: 85, Defence: 70},
				},
				Stamina: models.Stamina{Stamina: 88},
				Fitness: models.Fitness{Strength: 78, Agility: 90, InjuryTolerance: 82, InjuryResistance: 85},
			},
			{
				Name:         "Gabriel Jesus",
				Position:     models.Striker,
				Number:       9,
				Form:         86,
				Adaptability: 86,
				Composure:    86,
				Technical: models.TechnicalSkill{
					Speed:     models.SpeedSkill{Speed: 88, Acceleration: 90},
					Passing:   models.PassingSkill{ShortPass: 85, LongPass: 78, Cross: 80, Lob: 78, ThroughBall: 85, Chip: 82},
					Shooting:  models.ShootingSkill{Power: 85, Curve: 88, Finishing: 90, Spin: 85},
					Defending: models.DefendingSkill{Jumping: 78, Interceptions: 70, Heading: models.HeadingSkill{Accuracy: 78, Power: 75}, Blocking: 72},
					FreeKicks: 75,
					Penalties: 82,
				},
				TacticalIntelligence: models.TacticalIntelligence{
					Positioning: 88,
					Vision:      models.TacticalVision{Passing: 88, Shooting: 90, Defence: 72},
				},
				Stamina: models.Stamina{Stamina: 88},
				Fitness: models.Fitness{Strength: 80, Agility: 90, InjuryTolerance: 82, InjuryResistance: 80},
			},
		},
	}
}
