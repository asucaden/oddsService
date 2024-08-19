package data

import "time"

type Event struct {
	Id            string      `json:"id"`
	Sport_key     string      `json:"sport_key"`
	Sport_title   string      `json:"sport_title"`
	Commence_time time.Time   `json:"commence_time"`
	Home_team     string      `json:"home_team"`
	Away_team     string      `json:"away_team"`
	Bookmakers    []Bookmaker `json:"bookmakers"`
}

type Bookmaker struct {
	Key         string    `json:"key"`
	Title       string    `json:"title"`
	Last_update time.Time `json:"last_update"`
	Markets     []Market  `json:"markets"`
}

type Market struct {
	Key         string    `json:"key"`
	Last_update time.Time `json:"last_update"`
	Outcomes    []Outcome `json:"outcomes"`
}

type Outcome struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Point       float64 `json:"point"`
	Description string  `json:"description"`
}
