package main

import (
	"errors"
	"github.com/gorilla/mux"
	ravendb "github.com/ravendb/ravendb-go-client"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
)

var store, _ = getDocumentStore("NCAAF")

func getDocumentStore(databaseName string) (*ravendb.DocumentStore, error) {
	serverNodes := []string{"http://dell-4290.local:5050"}
	store := ravendb.NewDocumentStore(serverNodes, databaseName)
	if err := store.Initialize(); err != nil {
		return nil, err
	}
	return store, nil
}

func (p *Team) GetRecord(teamName string, weekNum int) string {
	return "This is record for " + teamName + ", Week: " +
		strconv.Itoa(weekNum)
}

func getAPSeason(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	year := vars["year"]
	log.Print("Year: ", year)

	var session, err = store.OpenSession("")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	q := session.QueryCollection("Polls")
	q = q.WhereEquals("year", year)
	//q = q.OrderBy("week")

	var polls []*Poll
	err = q.GetResults(&polls)
	if err != nil {
		panic(err)
	}

	sort.Slice(polls, func(i, j int) bool {
		return polls[i].Week < polls[j].Week
	})

	// Get Teams
	q = session.QueryCollection("Teams")

	var teams []*Team
	err = q.GetResults(&teams)
	if err != nil {
		panic(err)
	}

	//json.NewEncoder(w).Encode(teams)

	var i, _ = strconv.Atoi(year)

	var season = Season{
		Year: i,
	}

	var xPosition = 0

	for _, poll := range polls {
		xPosition += 250

		var week = Week{
			Number:    poll.Week,
			XPosition: xPosition,
		}
		var yPosition = 50

		for _, teamName := range poll.TeamNames {

			team, err := getTeam(teamName, teams)

			if err != nil {
				log.Println(err)
			} else {
				team.Cx = xPosition
				team.Cy = yPosition
				week.Teams = append(week.Teams, *team)
			}

			yPosition += 75
		}
		season.Weeks = append(season.Weeks, week)
		poll = nil // Free memory??
	}

	getPaths(&season)
	//season.Paths = getPaths(season)
	//json.NewEncoder(w).Encode(season)

	// Files are provided as a slice of strings.
	paths := []string{
		"AP-Season.tmpl",
	}

	funcMap := template.FuncMap{
		// The name "offset" is what the function will be called in the template text.
		"offset": func(i int, offset int) int {
			return i + offset
		},
	}

	season.Title = "NCAAF AP " + year

	t := template.
		Must(template.New("AP-Season.tmpl").
			Funcs(funcMap).
			ParseFiles(paths...))

	err = t.Execute(w, season)
	if err != nil {
		panic(err)
	}

}

func getTeam(name string, teams []*Team) (*Team, error) {
	for i := range teams {
		if teams[i].Name == name {
			return teams[i], nil
		}
	}
	var err = errors.New(name + " not found")
	return nil, err
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/ap/{year}", getAPSeason)

	log.Printf("Running...")
	log.Fatal(http.ListenAndServe(":10000", router))
}
