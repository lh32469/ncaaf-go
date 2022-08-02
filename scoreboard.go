package main

import (
	"encoding/json"
	"strings"
)

type ScoreBoard struct {
	ID          string
	UpdatedAt   string        `json:"updated_at"`
	InputMD5Sum string        `json:"inputMD5Sum"`
	Games       []GameWrapper `json:"games"`
}

type GameWrapper struct {
	Game Game `json:"game"`
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

func (t ScoreBoard) String() string {

	out, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(out)
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
