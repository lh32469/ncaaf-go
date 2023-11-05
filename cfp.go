package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func getCFPseason(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	year := vars["year"]

	log.Print("Year: ", year)

	var session = openSession()
	defer session.Close()

	// Load College Football Playoff Rankings
	q := session.QueryCollection("CFBDWeeks")
	q = q.WhereEquals("season", year)
	q = q.WhereEquals("polls.Poll", "Playoff Committee Rankings")

	var weeks []*CFBDWeek
	err := q.GetResults(&weeks)
	if err != nil {
		panic(err)
	}

	// Sort the Weeks
	sort.Slice(weeks, func(i, j int) bool {
		return weeks[i].Week < weeks[j].Week
	})

	// Get Teams
	q = session.QueryCollection("Teams")

	var teams []*Team
	err = q.GetResults(&teams)
	if err != nil {
		panic(err)
	}

	//json.NewEncoder(w).Encode(teams)

	var y, _ = strconv.Atoi(year)

	games := getGames(y)

	var season = Season{
		Year: y,
	}

	var xPosition = 0

	var weekPoll CFBDPoll

	for _, cfbdPoll := range weeks {
		xPosition += 250

		//log.Printf("Week %d\n", strconv.Itoa(cfbdPoll.Week))

		for _, poll := range cfbdPoll.Polls {
			log.Printf("Poll: %s\n", poll.Poll)
			if poll.Poll == "Playoff Committee Rankings" {
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
				team.Position = rank.Rank
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
		"CFP-Season.tmpl",
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
			game := getGame(team, week, games)
			return game.Result()
		},
		// Get Ranking of opponent in current Teams list
		"getRank": func(team Team, Teams []Team) int {
			for i, entry := range Teams {
				if entry.Name == team.Name {
					//log.Println(entry.Name + " " + strconv.Itoa(i))
					return i
				}
			}
			return 100 // Not Ranked
		},
		// Get record of current team based on this season's games
		"getRecord": func(team Team, week int) string {
			wins := 0
			losses := 0
			for _, game := range games {
				if game.Week > week-1 {
					continue
				}
				if game.AwayTeam == team.Name {
					if game.AwayPoints > game.HomePoints {
						wins++
					} else {
						losses++
					}
				}
				if game.HomeTeam == team.Name {
					if game.HomePoints > game.AwayPoints {
						wins++
					} else {
						losses++
					}
				}
			}
			return team.Name + " (" + strconv.Itoa(team.Position) + ")\n" +
				strconv.Itoa(wins) + " - " + strconv.Itoa(losses)
		},
	}

	season.Title = "NCAAF CFP " + year

	t := template.
		Must(template.New("CFP-Season.tmpl").
			Funcs(funcMap).
			ParseFiles(paths...))

	err = t.Execute(w, season)
	if err != nil {
		panic(err)
	}

}
