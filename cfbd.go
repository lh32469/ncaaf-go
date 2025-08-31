package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type CFBDWeek struct {
	Season     int        `json:"season"`
	SeasonType string     `json:"seasonType"`
	Week       int        `json:"week"`
	Polls      []CFBDPoll `json:"polls"`
}
type CFBDPoll struct {
	Poll  string
	Ranks []CFBDRank `json:"ranks"`
}

type CFBDRank struct {
	Rank       int
	School     string
	Conference string
}

type CFBDGame struct {
	ID          int       `json:"id"`
	Season      int       `json:"season"`
	SeasonType  string    `json:"season_type"`
	Week        int       `json:"week"`
	StartDateV1 time.Time `json:"start_date"`
	StartDateV2 time.Time `json:"startDate"`
	Venue       string    `json:"venue"`

	//seasonType     string
	//startDate      string
	//startTimeTbd   int
	//neutralSite    bool
	//conferenceGame bool
	//
	//attendance      int
	//venueId         int
	//venue           string
	//homeId          int
	HomeTeamV1 string `json:"home_team"`
	HomeTeamV2 string `json:"homeTeam"`
	//homeConference  string
	//homeDivision    string
	HomePointsV1 int `json:"home_points"`
	HomePointsV2 int `json:"homePoints"`
	//HomeLineScores []int `json:"home_line_scores"`
	//homePostWinProb float32
	//homePregameElo  int
	//homePostgameElo int
	//awayId          int
	AwayTeamV1 string `json:"away_team"`
	AwayTeamV2 string `json:"awayTeam"`
	//awayConference  string
	//awayDivision    string
	AwayPointsV1 int `json:"away_points"`
	AwayPointsV2 int `json:"awayPoints"`
	//AwayLineScores []int `json:"away_line_scores"`
	//awayPostWinProb float32
	//awayPregameElo  int
	//awayPostgameElo int
	//excitementIndex int
	//highlights      string
	Notes string `json:"notes"`
}

func (game *CFBDGame) StartDate() time.Time {
	if game.StartDateV1.IsZero() {
		return game.StartDateV2
	}
	return game.StartDateV1
}

func (game *CFBDGame) HomeTeam() string {
	if game.HomeTeamV1 != "" {
		return game.HomeTeamV1
	}
	return game.HomeTeamV2
}

func (game *CFBDGame) HomePoints() int {
	if game.HomePointsV1 == 0 {
		return game.HomePointsV2
	}
	return game.HomePointsV1
}

func (game *CFBDGame) AwayTeam() string {
	if game.AwayTeamV1 != "" {
		return game.AwayTeamV1
	}
	return game.AwayTeamV2
}

func (game *CFBDGame) AwayPoints() int {
	if game.AwayPointsV1 == 0 {
		return game.AwayPointsV2
	}
	return game.AwayPointsV1
}

func (game CFBDGame) Result() string {

	if game.ID != -1 {
		var loc, err = time.LoadLocation("America/Los_Angeles")
		if err != nil {
			panic(err)
		}
		return game.AwayTeam() + " (" +
			strconv.Itoa(game.AwayPoints()) + ") @ " +
			game.HomeTeam() + " (" +
			strconv.Itoa(game.HomePoints()) + ")\n" +
			game.StartDate().In(loc).Format("01-02; 15:04 MST") + "\n" +
			game.Venue + "\n" +
			game.Notes
	} else {
		return "Bye Week"
	}
}

func (game CFBDGame) Winner() string {
	if game.HomePoints() > game.AwayPoints() {
		return game.HomeTeam()
	} else {
		return game.AwayTeam()
	}
}

func (game CFBDGame) IsWinner(team Team) bool {

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

func (t CFBDWeek) String() string {

	out, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(out)
}

func (t CFBDGame) String() string {

	out, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(out)
}
