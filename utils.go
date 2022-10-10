package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func loadAllGames(w http.ResponseWriter, r *http.Request) {

	base := "https://api.collegefootballdata.com/games?year=YEAR&week=WEEK&seasonType=TYPE"
	vars := mux.Vars(r)

	for i := 0; i <= 16; i++ {
		var url = strings.ReplaceAll(base, "YEAR", vars["year"])
		url = strings.ReplaceAll(url, "WEEK", strconv.Itoa(i))
		url = strings.ReplaceAll(url, "TYPE", "regular")
		fmt.Println(url)

	}
}

func loadGames(w http.ResponseWriter, r *http.Request) {

	token := os.Getenv("CFDB_TOKEN")

	base := "https://api.collegefootballdata.com/games?year=YEAR&week=WEEK&seasonType=TYPE"
	vars := mux.Vars(r)

	var url = strings.ReplaceAll(base, "YEAR", vars["year"])
	url = strings.ReplaceAll(url, "WEEK", vars["week"])
	url = strings.ReplaceAll(url, "TYPE", vars["type"])
	fmt.Println(url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(body))

	var games []CFBDGame
	err = json.Unmarshal(body, &games)

	//fmt.Println(games)

	if vars["type"] == "postseason" {
		for i, _ := range games {
			// Can't use 'game' var here as it's a temp variable
			// and doesn't modify original slice
			games[i].Week = 15
		}
	}

	var session = openSession()
	defer session.Close()

	for i, game := range games {
		fmt.Printf("Storing Game: %d\n", game.Week)
		err = session.StoreWithID(&games[i], strconv.Itoa(games[i].ID))
		if err != nil {
			panic(err)
		}

	}

	err = session.SaveChanges()
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("%s", games)))
}

func getRankings(w http.ResponseWriter, r *http.Request) {

	token := os.Getenv("CFDB_TOKEN")

	base := "https://api.collegefootballdata.com/rankings?year=YEAR&week=WEEK&seasonType=TYPE"
	vars := mux.Vars(r)

	var url = strings.ReplaceAll(base, "YEAR", vars["year"])
	url = strings.ReplaceAll(url, "WEEK", vars["week"])
	url = strings.ReplaceAll(url, "TYPE", vars["type"])

	fmt.Println(url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(req)

	//res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(body))

	var weeks []CFBDWeek
	err = json.Unmarshal(body, &weeks)

	fmt.Println(weeks)

	var session = openSession()
	defer session.Close()

	for _, week := range weeks {
		id := strconv.Itoa(week.Season) + "." +
			strings.ToUpper(vars["type"][0:1]) + "." +
			strconv.Itoa(week.Week)

		err = session.StoreWithID(&week, id)
		err = session.Store(&week)
		if err != nil {
			panic(err)
		}
	}

	err = session.SaveChanges()
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("%s", weeks)))
}

func loadGamesForWeek(season int, week int, token string) {

	log.Printf("Loading Games for Week %d/%d\n", season, week)

	base := "https://api.collegefootballdata.com/games?year=YEAR&week=WEEK&seasonType=TYPE"

	var url = strings.ReplaceAll(base, "YEAR", strconv.Itoa(season))
	url = strings.ReplaceAll(url, "WEEK", strconv.Itoa(week))
	url = strings.ReplaceAll(url, "TYPE", "regular")

	fmt.Println(url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(body))

	var games []CFBDGame
	err = json.Unmarshal(body, &games)

	//fmt.Println(games)

	//if vars["type"] == "postseason" {
	//	for i, _ := range games {
	//		// Can't use 'game' var here as it's a temp variable
	//		// and doesn't modify original slice
	//		games[i].Week = 15
	//	}
	//}

	var session = openSession()
	defer session.Close()

	for i, _ := range games {
		err = session.StoreWithID(&games[i], strconv.Itoa(games[i].ID))
		if err != nil {
			panic(err)
		}
	}

	err = session.SaveChanges()
	if err != nil {
		panic(err)
	}

}

func getRankingsForWeek(season int, week int, token string) {

	log.Printf("Loading Rankings for Week %d/%d\n", season, week)

	base := "https://api.collegefootballdata.com/rankings?year=YEAR&week=WEEK&seasonType=TYPE"

	var url = strings.ReplaceAll(base, "YEAR", strconv.Itoa(season))
	url = strings.ReplaceAll(url, "WEEK", strconv.Itoa(week))
	url = strings.ReplaceAll(url, "TYPE", "regular")

	log.Println(url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := client.Do(req)

	//res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(body))

	var weeks []CFBDWeek
	err = json.Unmarshal(body, &weeks)

	fmt.Println(weeks)

	var session = openSession()
	defer session.Close()

	for _, week := range weeks {
		id := strconv.Itoa(week.Season) + "." +
			"R" + "." +
			strconv.Itoa(week.Week)

		err = session.StoreWithID(&week, id)
		err = session.Store(&week)
		if err != nil {
			panic(err)
		}
	}

	err = session.SaveChanges()
	if err != nil {
		panic(err)
	}

}
