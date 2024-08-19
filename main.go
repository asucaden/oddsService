package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/asucaden/oddsService/data"
	"github.com/asucaden/oddsService/persistence"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to db
	db := persistence.ConnectDb()
	defer db.Close()

	// oddsApiKey := os.Getenv("ODDS_API_KEY")
	// requestUrl := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/americanfootball_nfl/odds/?apiKey=%s&regions=us&markets=h2h,spreads&oddsFormat=american&bookmakers=betmgm,fanduel", oddsApiKey)
	// resp, err := http.Get(requestUrl)
	// if err != nil {
	// fmt.Println("Error: " + err.Error())
	// }
	// defer resp.Body.Close()
	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(body))
	// _ = os.WriteFile("response.json", body, 0777)

	fileName := "response.json"
	body, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var events []data.Event

	// Unmarshal the JSON data into the events slice
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Now you can work with the events
	for _, event := range events {
		if len(event.Bookmakers) > 0 && len(event.Bookmakers[0].Markets) > 0 {
			var h2hOdds int[]
			var outrightOdds int[]
			for _, market := range event.B

			fmt.Printf("Event team: %s, Event time: %s, Bet type: %s\n", event.Away_team, event.Commence_time, event.Bookmakers[0].Title)
		}
	}

}

// https://api.the-odds-api.com/v4/sports/americanfootball_nfl/odds/?apiKey={{apiKey}}&regions=us&markets=h2h,spreads&oddsFormat=american&bookmakers=betmgm,fanduel
