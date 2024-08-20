package persistence

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// INSERT
func AddCompetition(db *sqlx.DB, competition *Competition) (string, error) {
	var id string

	row := db.QueryRow("INSERT INTO competition (competition_id, competition_name, event_status, event_date)"+
		"VALUES ($1, $2, $3, $4)"+
		"RETURNING competition_id",
		competition.CompetitionId, competition.CompetitionName, competition.EventStatus, competition.EventDate)

	if err := row.Scan(&id); err != nil {
		return "", fmt.Errorf("addCompetition: %v", err)
	}

	return id, nil
}

// SELECT

// many
func AllCompetitions(db *sqlx.DB) ([]Competition, error) {
	var competitions []Competition
	rows, err := db.Query("SELECT * FROM competition")
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

func OneCompetition(db *sqlx.DB, competitionId string) (Competition, error) {
	var competition Competition

	err := db.Get(&competition, "SELECT * FROM competition WHERE competition_id = $1", competitionId)
	if err != nil {
		return competition, fmt.Errorf("OneCompetition %s: %v", competitionId, err)
	}
	return competition, nil
}
