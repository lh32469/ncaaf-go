package main

import (
	"errors"
	"math"
	"strconv"
)

func addPaths(season *Season, games []*CFBDGame) []Path {

	var paths []Path
	var previousWeek = Week{
		Number:    0,
		XPosition: 0,
		Teams:     []Team{},
	}

	for _, week := range season.Weeks {

		var current = week.Teams
		for _, team := range current {

			prev, err := findTeam(team.Name, previousWeek)
			if err != nil {
				continue
			}

			curr, err := findTeam(team.Name, week)

			var startX = prev.Cx + 70
			var startY = prev.Cy + 35
			var endX = curr.Cx
			var endY = curr.Cy + 35

			var cpath = "M " +
				strconv.Itoa(startX) + " " +
				strconv.Itoa(startY) + " C " +
				strconv.Itoa(startX+50) + " " +
				strconv.Itoa(startY) + " " +
				strconv.Itoa(endX-50) + " " +
				strconv.Itoa(endY) + " " +
				strconv.Itoa(endX) + " " +
				strconv.Itoa(endY)

			game := getGame(team, previousWeek.Number, games)

			// Default for bye weeks
			stroke := "black"

			if game.ID != -1 {
				// If not bye week
				if game.IsWinner(team) {
					stroke = "green"
				} else {
					stroke = "red"
				}
			}

			diff := math.Abs(float64(game.HomePoints - game.AwayPoints))
			m := int(math.Ceil(diff / 10))
			if m == 0 {
				m = 1
			}
			width := strconv.Itoa(m)

			paths = append(paths, Path{
				D:           cpath,
				Stroke:      stroke,
				StrokeWidth: width,
			})
		}

		previousWeek = week
	}

	season.Paths = paths
	return paths
}

func findTeam(name string, week Week) (Team, error) {
	for i := range week.Teams {
		if week.Teams[i].Name == name {
			return week.Teams[i], nil
		}
	}
	var err = errors.New(name + " not found")
	return Team{}, err
}
