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

	// Create event object (from json for testing, or from API call when working for real)
	var events []data.Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Now you can work with the events
	var obs []persistence.OfferedBet
	var competitions []persistence.Competition
	for _, event := range events {
		if len(event.Bookmakers) == 0 || len(event.Bookmakers[0].Markets) == 0 {
			continue
		}
		var competition persistence.Competition

		competition.CompetitionId = event.Id
		competition.CompetitionName = fmt.Sprintf("%s: %s @ %s", event.Sport_title, event.Away_team, event.Home_team)
		competition.EventDate = event.Commence_time
		// Only using one bookmaker for now. Could use multiple later then average the results. May be tricky when odds strattle +/- 100
		bookmaker := event.Bookmakers[0]
		for _, market := range bookmaker.Markets {
			ob, err := buildOfferedBet(&competition, &market)
			if err != nil {
				fmt.Printf(err.Error())
				continue
			}
			obs = append(obs, *ob)
		}
		competitions = append(competitions, competition)
	}

	// Need to post the results (either Inserts or updates)
	// You have a competiton and a slice of offered bets
	// First, handle the competitions. You should check if a competition with that ID exists. If yes, UPDATE it. Else, INSERT a new competition

	tx, err := db.Beginx()
	if err != nil {
		// TODO: Add error handling. Retry logic and/or writing the competitions & offered bets to a file for later retrying.
		panic(err)
	}

	for _, competition := range competitions {
		persistence.AddCompetition(tx, competition)
	}

}

func calcAvgOdds(odds1 int, odds2 int) int {
	if odds1 < 0 && odds2 < 0 {
		return 100
	}
	if odds1 < 0 {
		odds1 = -odds1
	}
	if odds2 < 0 {
		odds2 = -odds2
	}
	return (odds1 + odds2) / 2
}

// Utility functions
func calcAvgSpread(spread1 float64, spread2 float64) float32 {
	return float32(spread1)
}

func buildOfferedBet(competition *persistence.Competition, market *data.Market) (*persistence.OfferedBet, error) {
	if len(market.Outcomes) != 2 {
		return nil, fmt.Errorf("buildOfferedBet: Market doesn't have 2 outcomes", market.Key)
	}
	outcome1 := market.Outcomes[0]
	outcome2 := market.Outcomes[1]

	var spread float32
	switch market.Key {
	case data.Spreads.String():
		spread = calcAvgSpread(outcome1.Point, outcome2.Point)
	case data.H2h.String():
		spread = 0.0
	default:
		return nil, fmt.Errorf("buildOfferedBet: bet type %w not supported", market.Key)
	}

	odds := calcAvgOdds(int(outcome1.Price), int(outcome2.Price))
	odds1 := odds
	odds2 := odds
	if outcome1.Price < 0 {
		odds1 = -odds
	} else {
		odds2 = -odds
	}

	ob := persistence.OfferedBet{
		OfferedBetName: competition.CompetitionName + " - Point Spread",
		Outcome1:       outcome1.Name,
		Outcome1Odds:   odds1,
		Outcome2:       outcome2.Name,
		Outcome2Odds:   odds2,
		PointSpread:    spread,
		EventDate:      competition.EventDate,
		EventStatus:    0,
		CompetitionId:  competition.CompetitionId,
	}

	return &ob, nil
}

// https://api.the-odds-api.com/v4/sports/americanfootball_nfl/odds/?apiKey={{apiKey}}&regions=us&markets=h2h,spreads&oddsFormat=american&bookmakers=betmgm,fanduel
