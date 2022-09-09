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
	ID         int       `json:"id"`
	Season     int       `json:"season"`
	SeasonType string    `json:"season_type"`
	Week       int       `json:"week"`
	StartDate  time.Time `json:"start_date"`
	Venue      string    `json:"venue"`

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
	HomeTeam string `json:"home_team"`
	//homeConference  string
	//homeDivision    string
	HomePoints int `json:"home_points"`
	//HomeLineScores []int `json:"home_line_scores"`
	//homePostWinProb float32
	//homePregameElo  int
	//homePostgameElo int
	//awayId          int
	AwayTeam string `json:"away_team"`
	//awayConference  string
	//awayDivision    string
	AwayPoints int `json:"away_points"`
	//AwayLineScores []int `json:"away_line_scores"`
	//awayPostWinProb float32
	//awayPregameElo  int
	//awayPostgameElo int
	//excitementIndex int
	//highlights      string
	//notes           string
}

func (game CFBDGame) Result() string {

	if game.ID != -1 {
		return game.AwayTeam + " (" +
			strconv.Itoa(game.AwayPoints) + ") @ " +
			game.HomeTeam + " (" +
			strconv.Itoa(game.HomePoints) + ")\n" +
			game.StartDate.Format("01-02-2006")
	} else {
		return "Bye Week"
	}
}

func (game CFBDGame) Winner() string {
	if game.HomePoints > game.AwayPoints {
		return game.HomeTeam
	} else {
		return game.AwayTeam
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
