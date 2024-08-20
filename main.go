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
			var competition persistence.Competition
			var homeTeam = event.Home_team
			var awayTeam = event.Away_team
			var h2hOdds []data.Outcome
			var spreadsOdds []data.Outcome
			bookmaker := event.Bookmakers[0] // Only using one bookmaker for now. Could use multiple later then average the results. Might get tricky when odds strattle +/- 100
			for _, market := range bookmaker.Markets {
				if market.Key == "h2h" {
					for _, outcome := range market.Outcomes {
						var h2h data.Outcome
						h2h.Name = outcome.Name
						h2h.Price = outcome.Price
						h2hOdds = append(h2hOdds, h2h)
					}
				} else if market.Key == "spreads" {
					for _, outcome := range market.Outcomes {
						var spreads data.Outcome
						spreads.Name = outcome.Name
						spreads.Price = outcome.Price
						spreads.Point = outcome.Point
						spreadsOdds = append(spreadsOdds, spreads)
					}
				}
			}

			// Calculate h2h odds

			// Build h2h offered bet

			// (Create competition and create h2h offered bet) or (update h2h offered bet)

			// Build spreads offered bet

			// Create spreads offered bet or update spreads offered bet

			fmt.Printf("Event team: %s, Event time: %s, Bet type: %s\n", event.Away_team, event.Commence_time, event.Bookmakers[0].Title)
		}
	}

}

// Currently incomplete
// func getAverageOdds(outcomes []data.Outcome, homeTeam string, awayTeam string) (homeOdds int, awayOdds int) {
// 	n := len(outcomes)
// 	homeOddsSum := 0.0
// 	awayOddsSum := 0.0
// 	for _, outcome := range outcomes {
// 		if outcome.Name == homeTeam {
// 			homeOddsSum += outcome.Price
// 		} else if outcome.Name == awayTeam {
// 			awayOddsSum += outcome.Price
// 		} else {
// 			fmt.Println(fmt.Errorf("\nERROR Calculating odds: Team name mismatch: %s, %s, %s\n", homeTeam, awayTeam, outcome.Name))
// 		}
// 	}
// 	homeOddsAvg := homeOddsSum / float64(n)
// 	awayOddsAvg := awayOddsSum / float64(n)

// }

// https://api.the-odds-api.com/v4/sports/americanfootball_nfl/odds/?apiKey={{apiKey}}&regions=us&markets=h2h,spreads&oddsFormat=american&bookmakers=betmgm,fanduel
