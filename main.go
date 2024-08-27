package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/asucaden/oddsService/data"
	"github.com/asucaden/oddsService/persistence"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

const (
	useApi   = false
	fileName = "response.json"
)

func main() {
	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to db
	db := persistence.ConnectDb()
	defer db.Close()

	// Get data
	events := MustGetEventData("api") // "file" to use file or "api" to use the odds api

	// Process events
	obs, competitions := ProcessEvents(events)
	if err = UpdateObs(db, obs); err != nil {
		fmt.Println("Error updating offered bets: " + err.Error())
	}

	if err = UpdateCompetitions(db, competitions); err != nil {
		fmt.Println("Error updating competitions: " + err.Error())
	}
}

// Need to post the results (either Inserts or updates)
// You have a competiton and a slice of offered bets
// First, handle the competitions. You should check if a competition with that ID exists. If yes, UPDATE it. Else, INSERT a new competition

// Next, update all the offered bets.

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

func RequestNflOdds(storeResults bool) ([]byte, error) {
	oddsApiKey := os.Getenv("ODDS_API_KEY")
	requestUrl := fmt.Sprintf("https://api.the-odds-api.com/v4/sports/americanfootball_nfl/odds/?apiKey=%s&regions=us&markets=h2h,spreads&oddsFormat=american&bookmakers=betmgm,fanduel", oddsApiKey)
	resp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if storeResults {
		if err = os.WriteFile("response.json", body, 0777); err != nil {
			return nil, err
		}
	}

	return body, nil

}

func calcAvgSpread(spread1 float64, _ float64) float32 {
	return float32(spread1)
}

func buildOfferedBet(competition *persistence.Competition, market *data.Market) (*persistence.OfferedBet, error) {
	if len(market.Outcomes) != 2 {
		return nil, fmt.Errorf("buildOfferedBet: Market doesn't have 2 outcomes")
	}
	outcome1 := market.Outcomes[0]
	outcome2 := market.Outcomes[1]

	var spread float32
	switch market.Key {
	case "spreads":
		spread = calcAvgSpread(outcome1.Point, outcome2.Point)
	case "h2h":
		spread = 0.0
	default:
		return nil, fmt.Errorf("buildOfferedBet: bet type %s not supported", market.Key)
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

func MustGetEventData(strat string) []data.Event {
	var body []byte
	var err error
	switch strat {
	case "api":
		fmt.Println("Calling the API")
		body, err = RequestNflOdds(false)
	case "file":
		fmt.Println("Reading existing file")
		body, err = os.ReadFile(fileName)
	}
	if err != nil {
		panic(err)
	}

	// Convert to event
	var events []data.Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		panic(err)
	}
	return events
}

func ProcessEvents(events []data.Event) ([]persistence.OfferedBet, []persistence.Competition) {
	var competitions []persistence.Competition
	var obs []persistence.OfferedBet
	for _, event := range events {
		if len(event.Bookmakers) == 0 || len(event.Bookmakers[0].Markets) == 0 {
			continue
		}
		var competition persistence.Competition
		competition.CompetitionId = event.Id
		competition.CompetitionName = fmt.Sprintf("%s: %s @ %s", event.Sport_title, event.Away_team, event.Home_team)
		competition.EventDate = event.Commence_time

		bookmaker := event.Bookmakers[0] // Consider using multiple bookmakers later. May be tricky when odds stratle +/- 100
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
	return obs, competitions
}

func UpdateCompetitions(db *sqlx.DB, competitions []persistence.Competition) error {
	tx, err := db.Beginx()
	if err != nil {
		// TODO: Add error handling. Retry logic and/or writing the competitions & offered bets to a file for later retrying.
		return err
	}
	competitionsNum := len(competitions)
	addedCompetitions := 0
	updatedCompetitions := 0
	failedCompetitions := 0
	for _, competition := range competitions {
		queriedCompetition, _ := persistence.OneCompetition(tx, competition.CompetitionId)
		if queriedCompetition != nil { // If a competition was found and returned
			_, err = persistence.UpdateCompetition(tx, &competition)
			if err != nil {
				fmt.Println("Error updating competition: " + err.Error())
				failedCompetitions++
				continue
			}
			updatedCompetitions++
			continue
		}
		_, err = persistence.AddCompetition(tx, &competition)
		if err != nil {
			fmt.Println("Error updating competition: " + err.Error())
			failedCompetitions++
			continue
		}
		addedCompetitions++
	}
	if err = tx.Commit(); err != nil {
		fmt.Println("Error running competition transaction: " + err.Error())
	}
	fmt.Printf("%d competitions found\n", competitionsNum)
	fmt.Printf("%d existing competitions updated\n", updatedCompetitions)
	fmt.Printf("%d new competitions found\n", addedCompetitions)
	fmt.Printf("%d competitions failed to be added/updated\n", failedCompetitions)
	return nil
}

func UpdateObs(db *sqlx.DB, obs []persistence.OfferedBet) error {
	tx, err := db.Beginx()
	if err != nil {
		// TODO: Add error handling. Retry logic and/or writing the competitions & offered bets to a file for later retrying.
		return err
	}
	obNum := len(obs)
	addedObs := 0
	updatedObs := 0
	failedObs := 0
	for _, ob := range obs {
		queriedOb, _ := persistence.OneOfferedBet(tx, ob.OfferedBetId)
		if queriedOb != nil { // If an ob was found and returned
			_, err = persistence.UpdateOfferedBet(tx, &ob)
			if err != nil {
				fmt.Println("Error updating offered bet: " + err.Error())
				failedObs++
				continue
			}
			updatedObs++
			continue
		}
		_, err = persistence.AddOfferedBet(tx, &ob)
		if err != nil {
			fmt.Println("Error updating offered bet: " + err.Error())
			failedObs++
			continue
		}
		addedObs++
	}
	if err = tx.Commit(); err != nil {
		fmt.Println("Error running competition transaction: " + err.Error())
	}
	fmt.Printf("%d offered bets found\n", obNum)
	fmt.Printf("%d existing offered bets updated\n", updatedObs)
	fmt.Printf("%d new offered bets found\n", addedObs)
	fmt.Printf("%d offered bets failed to be added/updated\n", failedObs)
	return nil
}
