package main

import (
	"encoding/json"
)

type Poll struct {
	Poll      string   `json:"poll"`
	Year      string   `json:"year"`
	Week      int      `json:"week"`
	TeamNames []string `json:"teams"`
}

type Team struct {
	Name     string
	Names    []string
	Image    string
	Position int
	Cx       int
	Cy       int
}

func (t Team) String() string {

	out, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(out)
}

type Week struct {
	Number    int
	XPosition int
	Teams     []Team
}

type Season struct {
	Title string
	Year  int
	Weeks []Week
	Paths []Path
}

type Path struct {
	D           string
	Stroke      string
	StrokeWidth string
}
