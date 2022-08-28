package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func getOpponent(team Team, week int, games []*CFBDGame) Team {

	//log.Printf("getOpponent %s, %d, %d", team.Name, year, week)

	game := getGame2(team, week, games)
	//log.Printf("getOpponent: Game = %s", game)

	opponent := Team{}

	if strings.TrimSpace(game.Away.Names["short"]) == team.Name {
		//log.Printf("Found: " + team.Name)
		opponent.Name = game.Home.Names["short"]
		return opponent
	}
	if strings.TrimSpace(game.Home.Names["short"]) == team.Name {
		//log.Printf("Found: " + team.Name)
		opponent.Name = game.Away.Names["short"]
		return opponent
	}

	// Check aliases
	for _, alias := range team.Names {
		if strings.TrimSpace(game.Away.Names["short"]) == alias {
			//log.Printf("getGame Found: " + team.Name)
			opponent.Name = game.Home.Names["short"]
			return opponent
		}
		if strings.TrimSpace(game.Home.Names["short"]) == alias {
			//log.Printf("getGame Found: " + team.Name)
			opponent.Name = game.Away.Names["short"]
			return opponent
		}
	}

	//log.Printf("getOpponent %s, %d Not found", team.Name, week)

	return Team{}

}

func getCFBDGame(team Team, year int, week int) Game {
	var session = openSession()
	defer session.Close()

	q := session.QueryCollection("CFBDGames")
	q = q.WhereEquals("season", year)
	q = q.WhereEquals("week", week)
	q = q.WhereEquals("away_team", team.Name)
	q = q.OrElse()
	q = q.WhereEquals("season", year)
	q = q.WhereEquals("week", week)
	q = q.WhereEquals("home_team", team.Name)

	var games []*CFBDGame
	var err = q.GetResults(&games)
	if err != nil {
		panic(err)
	}

	game := games[0]

	home := Side{
		Score: strconv.Itoa(game.HomePoints),
		Names: map[string]string{
			"short": game.HomeTeam,
		},
	}

	away := Side{
		Score: strconv.Itoa(game.AwayPoints),
		Names: map[string]string{
			"short": game.AwayTeam,
		},
	}

	if game.HomePoints > game.AwayPoints {
		home.Winner = true
	} else {
		away.Winner = true
	}

	return Game{
		GameID: strconv.Itoa(games[0].ID),
		Home:   home,
		Away:   away,
	}
}

func getGame(team Team, week int, scoreboards []*ScoreBoard) Game {

	//log.Printf("getGame %s, %d", team.Name, week)

	ws := fmt.Sprintf("%02d", week)
	game := Game{}

	for _, scoreboard := range scoreboards {
		if strings.HasSuffix(scoreboard.ID, ws) {
			//log.Printf("getGame Found: %s\n", scoreboard.ID)
			//log.Printf("getGame Found: %s\n", scoreboard.InputMD5Sum)

			for _, wrapper := range scoreboard.Games {
				game := wrapper.Game

				//log.Printf("Game: %s\n", game.Title)
				if strings.TrimSpace(game.Away.Names["short"]) == team.Name {
					//log.Printf("getGame Found: " + team.Name)
					team.Name = game.Home.Names["short"]
					return game
				}
				if strings.TrimSpace(game.Home.Names["short"]) == team.Name {
					//log.Printf("getGame Found: " + team.Name)
					team.Name = game.Away.Names["short"]
					return game
				}

				// Check aliases
				for _, alias := range team.Names {
					if strings.TrimSpace(game.Away.Names["short"]) == alias {
						//log.Printf("getGame Found: " + team.Name)
						team.Name = game.Home.Names["short"]
						return game
					}
					if strings.TrimSpace(game.Home.Names["short"]) == alias {
						//log.Printf("getGame Found: " + team.Name)
						team.Name = game.Away.Names["short"]
						return game
					}
				}
			}

		}
	}

	return game
}

func getScoreBoard(year int, week int) *ScoreBoard {

	session := openSession()
	defer session.Close()

	id := fmt.Sprintf("scoreboard.%d.%02d", year, week)

	log.Println("getScoreBoard ID: " + id)

	q := session.QueryCollection("ScoreBoards")
	q = q.Where("ID", "==", id)
	var scoreBoard *ScoreBoard
	var err = q.Single(&scoreBoard)
	if err != nil {
		panic(err)
	}

	return scoreBoard
}

func getScoreBoards(year int) []*ScoreBoard {

	session := openSession()
	defer session.Close()

	id := fmt.Sprintf("scoreboard.%d", year)

	log.Println("ID: " + id)

	q := session.QueryCollection("ScoreBoards")
	q.WhereStartsWith("ID", id)
	var scoreBoards []*ScoreBoard
	var err = q.GetResults(&scoreBoards)
	if err != nil {
		panic(err)
	}

	log.Printf("Found %d ScoreBoards for %d",
		len(scoreBoards), year)

	return scoreBoards
}

func getGames(season int) []*CFBDGame {

	session := openSession()
	defer session.Close()
	q := session.QueryCollection("CFBDGames")
	q = q.WhereEquals("season", season)

	var games []*CFBDGame
	var err = q.GetResults(&games)
	if err != nil {
		panic(err)
	}

	return games
}

func getGame2(team Team, week int, games []*CFBDGame) Game {

	game := games[100]

	for _, cfbdGame := range games {
		if cfbdGame.Week == week {
			if cfbdGame.AwayTeam == team.Name {
				game = cfbdGame
				break
			}
			if cfbdGame.HomeTeam == team.Name {
				game = cfbdGame
				break
			}
		}
	}
	home := Side{
		Score: strconv.Itoa(game.HomePoints),
		Names: map[string]string{
			"short": game.HomeTeam,
		},
	}

	away := Side{
		Score: strconv.Itoa(game.AwayPoints),
		Names: map[string]string{
			"short": game.AwayTeam,
		},
	}

	if game.HomePoints > game.AwayPoints {
		home.Winner = true
	} else {
		away.Winner = true
	}

	return Game{
		GameID: strconv.Itoa(game.ID),
		Home:   home,
		Away:   away,
	}
}
