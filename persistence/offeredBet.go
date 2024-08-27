package persistence

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// INSERT

func AddOfferedBet(q Querier, ob *OfferedBet) (int, error) {
	var id int

	row := q.QueryRowx("INSERT INTO offered_bet (offered_bet_name, outcome1, outcome1odds, outcome2, outcome2odds, event_date, event_status, point_spread, competition_id)"+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"+
		"RETURNING offered_bet_id",
		ob.OfferedBetName, ob.Outcome1, ob.Outcome1Odds, ob.Outcome2, ob.Outcome2Odds, ob.EventDate, ob.EventStatus, ob.PointSpread, ob.CompetitionId)

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("addOfferedBet: %v", err)
	}

	return id, nil
}

func OfferedBetsByCompetiton(q Querier, competitionId string) ([]OfferedBet, error) {
	var obs []OfferedBet

	err := q.Select(&obs, "SELECT * FROM offered_bet WHERE competition_id = $1", competitionId)
	if err != nil {
		return nil, fmt.Errorf("OfferedBetsByCompetition: %v", err)
	}
	return obs, nil
}

// SELECT

// one
func OneOfferedBet(q Querier, offeredBetId int) (OfferedBet, error) {
	var ob OfferedBet

	err := q.Get(&ob, "SELECT * FROM offered_bet WHERE offered_bet_id = $1", offeredBetId)
	if err != nil {
		return ob, fmt.Errorf("OneOfferedBet %d: %v", offeredBetId, err)
	}
	return ob, nil
}

// partial

func OfferedBetNameById(q Querier, offeredBetId int) (string, error) {
	var obName string

	err := q.Select(&obName, "SELECT offered_bet_name FROM offered_bet WHERE offered_bet_id = $1", offeredBetId)
	if err != nil {
		return "", fmt.Errorf("OfferedBetNameById: %v", err)
	}
	return obName, nil
}
