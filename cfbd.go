package main

import "encoding/json"

type CFBDWeek struct {
	Season int
	Week   int
	Polls  []CFBDPoll `json:"polls"`
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
	ID     int `json:"id"`
	Season int `json:"season"`
	Week   int `json:"week"`
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
