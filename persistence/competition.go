package persistence

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func AddCompetition(q Querier, competition *Competition) (string, error) {
	var id string

	row := q.QueryRowx("INSERT INTO competition (competition_id, competition_name, event_status, event_date)"+
		"VALUES ($1, $2, $3, $4)"+
		"RETURNING competition_id",
		competition.CompetitionId, competition.CompetitionName, competition.EventStatus, competition.EventDate)

	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("addCompetition: %v", err)
	}

	return id, nil
}

// UPDATE

// one
func UpdateCompetition(q Querier, competition *Competition) (string, error) {
	var id string

	row := q.QueryRowx("UPDATE competition "+
		"SET competition_name = $1, "+
		"event_status = $2, "+
		"event_date = $3 "+
		"WHERE competition_id = $4 "+
		"RETURNING competition_id",
		competition.CompetitionName, competition.EventStatus, competition.EventDate, competition.CompetitionId)

	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("UpdateCompetition: %v", err)
	}

	return id, nil
}

// SELECT

// many
func AllCompetitions(q Querier) ([]Competition, error) {
	var competitions []Competition
	rows, err := q.Query("SELECT * FROM competition")
	if err != nil {
		return nil, fmt.Errorf("ERROR running the query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var competition Competition
		if err := rows.Scan(&competition.CompetitionId, &competition.CompetitionName, &competition.EventStatus, &competition.EventDate); err != nil {
			return nil, fmt.Errorf("allCompetitions: %v", err)
		}
		competitions = append(competitions, competition)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("allCompetitions: %v", err)
	}
	return competitions, nil
}

// one

func OneCompetition(q Querier, competitionId string) (*Competition, error) {
	var competition Competition

	err := q.Get(&competition, "SELECT * FROM competition WHERE competition_id = $1", competitionId)
	if err != nil {
		return nil, fmt.Errorf("OneCompetition %s: %v", competitionId, err)
	}
	return &competition, nil
}
