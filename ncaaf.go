package main

import (
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (p *Team) GetRecord(teamName string, weekNum int) string {
	return "This is record for " + teamName + ", Week: " +
		strconv.Itoa(weekNum)
}

func getAPSeason(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	year := vars["year"]

	if len(year) == 0 {
		var now = time.Now()
		var yr, _ = now.ISOWeek()
		year = strconv.Itoa(yr)
	}

	log.Print("Year: ", year)

	var session = openSession()
	defer session.Close()

	//q := session.QueryCollection("Polls")
	//q = q.WhereEquals("year", year)
	////q = q.OrderBy("week")

	//var polls []*Poll
	//var err = q.GetResults(&polls)
	//if err != nil {
	//	panic(err)
	//}

	// Load Regular Season Rankings
	q := session.QueryCollection("CFBDWeeks")
	q = q.WhereEquals("season", year)
	q = q.WhereEquals("seasonType", "regular")

	var weeks []*CFBDWeek
	err := q.GetResults(&weeks)
	if err != nil {
		panic(err)
	}

	// Load Post Season (Final) Rankings
	q = session.QueryCollection("CFBDWeeks")
	q = q.WhereEquals("season", year)
	q = q.WhereEquals("seasonType", "postseason")

	var post []*CFBDWeek
	err = q.GetResults(&post)
	if err != nil {
		panic(err)
	}

	for _, final := range post {
		final.Week = len(weeks) + 1
		weeks = append(weeks, final)
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

	router.HandleFunc("/", getAPSeason)

	router.HandleFunc("/ap/{year}", getAPSeason)
	router.HandleFunc("/AP/{year}", getAPSeason)
	//router.HandleFunc("/rankings/{year}/{week}/{type}", getRankings)
	//router.HandleFunc("/load/{year}/{week}/{type}", loadGames)
	router.HandleFunc("/image/{image}", getImage)

	s := gocron.NewScheduler(time.UTC)

	s.Cron("0 */2 * 8,9,10,11,12 SUN,MON").Do(func() {
		token := os.Getenv("CFDB_TOKEN")
		var now = time.Now()
		var year, week = now.ISOWeek()
		week = week - 34
		log.Printf("Loading CFB Data for Week %d/%d\n", year, week)
		getRankingsForWeek(year, week, token)
		loadGamesForWeek(year, week, token)
	})

	s.StartAsync()

	port := "10000"
	log.Printf("Running at port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
