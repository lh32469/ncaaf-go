package main

import (
	"encoding/json"
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
		Score: strconv.Itoa(game.HomePoints()),
		Names: map[string]string{
			"short": game.HomeTeam(),
		},
	}

	away := Side{
		Score: strconv.Itoa(game.AwayPoints()),
		Names: map[string]string{
			"short": game.AwayTeam(),
		},
	}

	if game.HomePoints() > game.AwayPoints() {
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

func getGame(team Team, week int, games []*CFBDGame) CFBDGame {
	for _, cfbdGame := range games {
		if cfbdGame.Week == week {
			if contains(team.Names, cfbdGame.AwayTeam()) {
				return *cfbdGame
			}
			if contains(team.Names, cfbdGame.HomeTeam()) {
				return *cfbdGame
			}
		}
	}

	return CFBDGame{
		ID: -1,
	}
}

func getGame2(team Team, week int, games []*CFBDGame) Game {

	game := games[0]

	for _, cfbdGame := range games {
		if cfbdGame.Week == week {
			if contains(team.Names, cfbdGame.AwayTeam()) {
				game = cfbdGame
				break
			}
			if contains(team.Names, cfbdGame.HomeTeam()) {
				game = cfbdGame
				break
			}
		}
	}

	home := Side{
		Score: strconv.Itoa(game.HomePoints()),
		Names: map[string]string{
			"short": game.HomeTeam(),
		},
	}

	away := Side{
		Score: strconv.Itoa(game.AwayPoints()),
		Names: map[string]string{
			"short": game.AwayTeam(),
		},
	}

	if game.HomePoints() > game.AwayPoints() {
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

type Game struct {
	GameID         string `json:"gameID"`
	Title          string
	StartTime      string
	StartTimeEpoch string
	Home           Side
	Away           Side
}

type Side struct {
	Score  string
	Winner bool
	Rank   string
	Names  map[string]string
}

func (t Game) String() string {

	out, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(out)
}

func (game Game) Result() string {

	return game.Home.Names["short"] + " " +
		game.Home.Score + " " +
		game.Away.Names["short"] + " " +
		game.Away.Score
}

func (game Game) Winner() string {
	if game.Home.Winner {
		return game.Home.Names["short"]
	} else {
		return game.Away.Names["short"]
	}
}

func (game Game) IsWinner(team Team) bool {

	winner := game.Winner()

	if team.Name == winner {
		return true
	}

	for _, alias := range team.Names {
		if strings.TrimSpace(winner) == alias {
			return true
		}
	}

	return false
}
