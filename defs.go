package main

type Poll struct {
	Poll      string   `json:"poll"`
	Year      string   `json:"year"`
	Week      int      `json:"week"`
	TeamNames []string `json:"teams"`
}

type Team struct {
	Name  string
	Image string
	Cx    int
	Cy    int
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
