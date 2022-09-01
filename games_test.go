package main

import (
	"fmt"
	"log"
	"testing"
)

func TestGetOpponent(t *testing.T) {

	log.Printf("Running...")

	michigan := Team{
		Name: "Michigan",
	}
	florida := Team{
		Name: "Florida",
	}

	games := getGames(2017)

	opponent := getOpponent(michigan, 1, games)

	//fmt.Printf("Opponent: %s\n", opponent)
	if "Florida" != opponent.Name {
		t.Logf("Opponent: %s\n", opponent)
		t.Fatalf(`Couldn't find Opponent for %s'`, michigan)
	}

	opponent = getOpponent(florida, 1, games)
	if michigan.Name != opponent.Name {
		t.Logf("Opponent: %s\n", opponent)
		t.Fatalf(`Couldn't find Opponent for %s'`, florida)
	}

}

func TestGetOpponentNCState(t *testing.T) {

	log.Printf("Running...")
	var session = openSession()
	defer session.Close()

	// Get Teams
	q := session.QueryCollection("Teams")

	var teams []*Team
	var err = q.GetResults(&teams)
	if err != nil {
		panic(err)
	}

	ncState, err := getTeam("NC State", teams)
	if err != nil {
		panic(err)
	}
	log.Printf("NC State: %s\n", ncState)

	games := getGames(2017)

	opponent := getOpponent(*ncState, 1, games)
	log.Printf("opponent: %s\n", opponent)

	if "South Carolina" != opponent.Name {
		t.Logf("Opponent: %s\n", opponent)
		t.Fatalf(`Couldn't find Opponent for %s'`, ncState)
	}

}

func TestGetOpponentOhioState(t *testing.T) {

	log.Printf("Running...")
	var session = openSession()
	defer session.Close()

	// Get Teams
	q := session.QueryCollection("Teams")

	var teams []*Team
	var err = q.GetResults(&teams)
	if err != nil {
		panic(err)
	}

	ohioState, err := getTeam("Ohio State", teams)
	if err != nil {
		panic(err)
	}
	log.Printf("NC State: %s\n", ohioState)

	games := getGames(2017)

	opponent := getOpponent(*ohioState, 1, games)
	log.Printf("opponent: %s\n", opponent)

	if "Indiana" != opponent.Name {
		t.Logf("Opponent: %s\n", opponent)
		t.Fatalf(`Couldn't find Opponent for %s'`, ohioState)
	}

	games = getGames(2018)

	opponent = getOpponent(*ohioState, 1, games)
	log.Printf("opponent: %s\n", opponent)

	if "Oregon State" != opponent.Name {
		t.Logf("Opponent: %s\n", opponent)
		t.Fatalf(`Couldn't find Opponent for %s'`, ohioState)
	}

}

func TestGetScoreBoard(t *testing.T) {

	log.Printf("Running...")

	scoreBoard := getScoreBoard(2021, 1)
	//fmt.Printf("scoreBoard: %s\n", scoreBoard)
	fmt.Printf("Game[0]: %s\n", scoreBoard.Games[0])

}

func TestGetScoreBoards(t *testing.T) {

	log.Printf("Running...")

	scoreBoards := getScoreBoards(2020)
	fmt.Printf("scoreBoard: %s\n", scoreBoards)
	//fmt.Printf("Game[0]: %s\n", scoreBoard.Games[0])

}

func TestGetGame(t *testing.T) {

	log.Printf("Running...")

	michigan := Team{
		Name: "Michigan",
	}

	games := getGames(2017)

	game := getGame(michigan, 2, games)
	fmt.Printf("Game: %s\n", game)

	log.Printf("Winner: %s", game.Winner())
	log.Printf("Is Winner: %t", game.IsWinner(michigan))

}

func TestGetCFBDGame(t *testing.T) {

	log.Printf("Running...")

	michigan := Team{
		Name: "Michigan",
	}

	game := getCFBDGame(michigan, 2017, 2)
	fmt.Printf("Game: %s\n", game)

	log.Printf("Winner: %s", game.Winner())
	log.Printf("Is Winner: %t", game.IsWinner(michigan))

}

func TestGetGame2(t *testing.T) {

	log.Printf("Running...")

	michigan := Team{
		Name: "Michigan",
	}

	games := getGames(2017)
	game := getGame2(michigan, 2, games)

	fmt.Printf("Game: %s\n", game)

	log.Printf("Winner: %s", game.Winner())
	log.Printf("Is Winner: %t", game.IsWinner(michigan))

}

func TestWinner(t *testing.T) {

	log.Printf("Running...")
	var session = openSession()
	defer session.Close()

	// Get Teams
	q := session.QueryCollection("Teams")

	var teams []*Team
	var err = q.GetResults(&teams)
	if err != nil {
		panic(err)
	}

	ohioState, err := getTeam("Ohio State", teams)
	if err != nil {
		panic(err)
	}
	log.Printf("NC State: %s\n", ohioState)

	games := getGames(2017)

	game := getGame(*ohioState, 1, games)
	log.Printf("game: %s\n", game)

	log.Printf("Winner: %s", game.Winner())
	log.Printf("Is Winner: %t", game.IsWinner(*ohioState))

}

func TestGetMichiganGame2021(t *testing.T) {

	log.Printf("Running...")
	var session = openSession()
	defer session.Close()

	games := getGames(2021)

	michigan := Team{
		Name: "Michigan",
	}

	game := getGameBy(michigan, 1, "regular", games)

	if "Western Michigan" != game.AwayTeam {
		t.Logf("Opponent: %s\n", game)
		t.Fatalf(`Couldn't find Opponent for %s'`, michigan)
	}

	game = getGameBy(michigan, 1, "postseason", games)

	if "Georgia" != game.AwayTeam {
		t.Logf("Opponent: %s\n", game)
		t.Fatalf(`Couldn't find Opponent for %s'`, michigan)
	}

}
