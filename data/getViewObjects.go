package data

import (
	"log"

	"github.com/asucaden/oddsService/persistence"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func GetAllCompetitionHeadlines(db *sqlx.DB) []CompetitionHeadline {
	competitions, err := persistence.AllCompetitions(db)
	if err != nil {
		log.Fatal(err)
	}

	var competitionHeadlines []CompetitionHeadline
	for _, competition := range competitions {
		competitionHeadlines = append(competitionHeadlines, CompetitionHeadline{Id: competition.CompetitionId, Title: competition.CompetitionName, Date: competition.EventDate})
	}

	return competitionHeadlines
}
