package main

import (
	"errors"
	"github.com/gorilla/mux"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func (p *Team) GetRecord(teamName string, weekNum int) string {
	return "This is record for " + teamName + ", Week: " +
		strconv.Itoa(weekNum)
}

func getAPSeason(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	year := vars["year"]
	log.Print("Year: ", year)

	var session = openSession()
	defer session.Close()

	q := session.QueryCollection("Polls")
	q = q.WhereEquals("year", year)
	//q = q.OrderBy("week")

	var polls []*Poll
	var err = q.GetResults(&polls)
	if err != nil {
		panic(err)
	}

	sort.Slice(polls, func(i, j int) bool {
		return polls[i].Week < polls[j].Week
	})

	// Load CFDB Poll data
	q = session.QueryCollection("CFBDWeeks")
	q = q.WhereEquals("Season", year)
	//q = q.OrderBy("week")

	var weeks []*CFBDWeek
	err = q.GetResults(&weeks)
	if err != nil {
		panic(err)
	}

	sort.Slice(weeks, func(i, j int) bool {
		return weeks[i].Week < weeks[j].Week
	})
	//fmt.Println(weeks)

	//var apPolls []CFBDPoll
	//
	//for _, week := range weeks {
	//	for _, poll := range week.Polls {
	//		if poll.Poll == "AP Top 25" {
	//			apPolls = append(apPolls, poll)
	//		}
	//	}
	//}
	//
	//fmt.Println(apPolls)

	// End Load CFDB Poll data

	// Get Teams
	q = session.QueryCollection("Teams")

	var teams []*Team
	err = q.GetResults(&teams)
	if err != nil {
		panic(err)
	}

	//json.NewEncoder(w).Encode(teams)

	var i, _ = strconv.Atoi(year)

	//scoreBoards := getScoreBoards(i)
	games := getGames(i)

	var season = Season{
		Year: i,
	}

	var xPosition = 0

	var weekPoll CFBDPoll

	for position, cfbdPoll := range weeks {
		xPosition += 250

		for _, poll := range cfbdPoll.Polls {
			if poll.Poll == "AP Top 25" {
				weekPoll = poll
			}
		}

		sort.Slice(weekPoll.Ranks, func(i, j int) bool {
			return weekPoll.Ranks[i].Rank < weekPoll.Ranks[j].Rank
		})

		var week = Week{
			Number:    cfbdPoll.Week,
			XPosition: xPosition,
		}

		var yPosition = 50

		for _, rank := range weekPoll.Ranks {

			team, err := getTeam(rank.School, teams)

			if err != nil {
				log.Println(err)
			} else {
				team.Position = position
				team.Cx = xPosition
				team.Cy = yPosition
				week.Teams = append(week.Teams, *team)
			}

			yPosition += 75

		}

		season.Weeks = append(season.Weeks, week)
		cfbdPoll = nil // Free memory??
	}

	addPaths(&season, games)
	//season.Paths = addPaths(season)
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
		"opponent": func(team Team, week int) Team {
			opp := getOpponent(team, week, games)
			other, err := getTeam(opp.Name, teams)
			if err != nil {
				return Team{}
			}
			return *other
		},
		"getResult": func(team Team, week int) string {
			game := getGame2(team, week, games)
			return game.Result()
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

	name = strings.TrimSpace(name)

	for _, team := range teams {
		if team.Name == name {
			return team, nil
		}
		for _, alias := range team.Names {
			if alias == name {
				return team, nil
			}
		}
	}
	log.Printf("getTeam: '%s' not found", name)
	var err = errors.New("'" + name + "' not found")
	return nil, err
}

func getImage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	image := vars["image"]
	//log.Print("image: ", image)
	data, _ := fs.ReadFile(os.DirFS("images"), image)

	w.Write(data)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/ap/{year}", getAPSeason)
	//router.HandleFunc("/cfbd/{year}/{week}", getCFBD)
	//router.HandleFunc("/load/{year}/{week}", loadGames)
	router.HandleFunc("/image/{image}", getImage)

	log.Printf("Running...")
	log.Fatal(http.ListenAndServe(":10000", router))
}
