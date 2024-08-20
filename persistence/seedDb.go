package persistence

import (
	"fmt"
	"log"
	"time"

	"github.com/asucaden/oddsService/auth"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// TODO improve performance using transactions or concurrency
func SeedDb(db *sqlx.DB) {

	// Clear all existing tables
	db.MustExec("DROP TABLE IF EXISTS users CASCADE")
	db.MustExec("DROP TABLE IF EXISTS competition CASCADE")
	db.MustExec("DROP TABLE IF EXISTS offered_bet CASCADE")
	db.MustExec("DROP TABLE IF EXISTS bet CASCADE")

	// Create tables
	schema := `
		CREATE TABLE IF NOT EXISTS users(
			user_id           serial PRIMARY KEY,
			username          varchar(63) NOT NULL,
			hash        	  varchar(63) NOT NULL,
			balance           int NOT NULL
		);

		CREATE TABLE IF NOT EXISTS competition(
			competition_id      text PRIMARY KEY,
			competition_name    varchar(255) NOT NULL,
			event_status        int NOT NULL,
			event_date          date NOT NULL
		);

		CREATE TABLE IF NOT EXISTS offered_bet(
			offered_bet_id      serial PRIMARY KEY,
			offered_bet_name    varchar(255) NOT NULL,
			outcome1            varchar(255) NOT NULL,
			outcome1odds        int NOT NULL,
			outcome2            varchar(255) NOT NULL,
			outcome2odds        int NOT NULL,
			point_spread		real NOT NULL,
			event_date          date NOT NULL,
			event_status        int NOT NULL,
			competition_id      text REFERENCES competition
		);

		CREATE TABLE IF NOT EXISTS bet(
			bet_id          serial PRIMARY KEY,
			bet_status      int NOT NULL,
			amount1         int NOT NULL,
			outcome1odds    int NOT NULL,
			user1agreed     boolean NOT NULL,
			user1id         int REFERENCES users,
			amount2         int NOT NULL,
			outcome2odds    int NOT NULL,
			user2agreed     boolean NOT NULL,
			user2id         int REFERENCES users,
			point_spread    real NOT NULL,
			customized		boolean NOT NULL,
			offered_bet_id  int REFERENCES offered_bet,
			settled         boolean NOT NULL
		);
	`
	db.MustExec(schema)

	// Populate tables
	// Users
	mustHashPass := func(password string) string {
		hash, err := auth.HashPassword(password)
		if err != nil {
			log.Fatal(err.Error())
		}
		return hash
	}
	mustAddUser := func(user *User) int {
		id, err := AddUser(db, user)
		if err != nil {
			log.Fatal(err.Error())
		}
		return id
	}
	uid0 := mustAddUser(&User{Username: "the house", Hash: mustHashPass("admin"), Balance: 1000000.0})
	uid1 := mustAddUser(&User{Username: "Caden M", Hash: mustHashPass("password"), Balance: 0.0})
	uid2 := mustAddUser(&User{Username: "Parker B", Hash: mustHashPass("password"), Balance: 0.0})
	uid3 := mustAddUser(&User{Username: "Ryan E", Hash: mustHashPass("password"), Balance: 0.0})
	uid4 := mustAddUser(&User{Username: "Alec V", Hash: mustHashPass("password"), Balance: 0.0})
	uid5 := mustAddUser(&User{Username: "Tanner B", Hash: mustHashPass("password"), Balance: 0.0})

	// Competitions
	mustAddCompetition := func(competition *Competition) string {
		id, err := AddCompetition(db, competition)
		if err != nil {
			log.Fatal(err.Error())
		}
		return id
	}
	mustParseTime := func(timeStr string) time.Time {
		time, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			log.Fatal(err.Error())
		}
		return time
	}

	cid0 := mustAddCompetition(&Competition{CompetitionId: "1", CompetitionName: "Nba Finals Game 5", EventStatus: 2, EventDate: mustParseTime("2024-06-17T00:00:00Z")})
	cid1 := mustAddCompetition(&Competition{CompetitionId: "2", CompetitionName: "Olympics 1500m Men's Final", EventStatus: 1, EventDate: mustParseTime("2024-08-06T00:00:00Z")})
	cid2 := mustAddCompetition(&Competition{CompetitionId: "3", CompetitionName: "USA Presidential Election", EventStatus: 1, EventDate: mustParseTime("2024-11-05T00:00:00Z")})

	// Offered Bets
	mustAddOfferedBet := func(ob *OfferedBet) int {
		id, err := AddOfferedBet(db, ob)
		if err != nil {
			log.Fatal(err.Error())
		}
		return id
	}
	obid0 := mustAddOfferedBet(&OfferedBet{OfferedBetName: "Finals Moneyline", Outcome1: "Mavs win", Outcome1Odds: 250, Outcome2: "Celtics win", Outcome2Odds: -250, EventDate: mustParseTime("2024-06-17T00:00:00Z"), PointSpread: 0, EventStatus: 2, CompetitionId: cid0})
	obid1 := mustAddOfferedBet(&OfferedBet{OfferedBetName: "Luka scores most points in series", Outcome1: "Yes", Outcome1Odds: 100, Outcome2: "No", Outcome2Odds: 100, EventDate: mustParseTime("2024-06-17T00:00:00Z"), PointSpread: 0, EventStatus: 1, CompetitionId: cid0})
	obid2 := mustAddOfferedBet(&OfferedBet{OfferedBetName: "Ingebrigtsen win", Outcome1: "Yes", Outcome1Odds: -120, Outcome2: "No", Outcome2Odds: 120, EventDate: mustParseTime("2024-08-06T00:00:00Z"), PointSpread: 0, EventStatus: 0, CompetitionId: cid1})
	obid3 := mustAddOfferedBet(&OfferedBet{OfferedBetName: "Election winner", Outcome1: "Biden", Outcome1Odds: 110, Outcome2: "Trump", Outcome2Odds: -110, EventDate: mustParseTime("2024-11-05T00:00:00Z"), PointSpread: 0, EventStatus: 0, CompetitionId: cid2})
	obid4 := mustAddOfferedBet(&OfferedBet{OfferedBetName: "Finals Point Spread", Outcome1: "Mavs win", Outcome1Odds: 100, Outcome2: "Celtics win", Outcome2Odds: -100, EventDate: mustParseTime("2024-06-17T00:00:00Z"), PointSpread: 3.5, EventStatus: 2, CompetitionId: cid0})

	// Bets
	mustAddBet := func(bet *Bet) {
		_, err := AddBet(db, bet)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	mustAddBet(&Bet{BetStatus: 1, Amount1: 10000, Outcome1Odds: 100, User1Agreed: false, User1Id: uid1, Amount2: 10000, Outcome2Odds: 100, User2Agreed: true, User2Id: uid0, Customized: false, PointSpread: 0, OfferedBetId: obid1, Settled: true})
	mustAddBet(&Bet{BetStatus: 2, Amount1: 1000, Outcome1Odds: 250, User1Agreed: true, User1Id: uid4, Amount2: 2500, Outcome2Odds: -250, User2Agreed: false, User2Id: uid2, Customized: false, PointSpread: 0, OfferedBetId: obid0, Settled: true})
	mustAddBet(&Bet{BetStatus: 2, Amount1: 2000, Outcome1Odds: 250, User1Agreed: true, User1Id: uid3, Amount2: 5000, Outcome2Odds: -250, User2Agreed: true, User2Id: uid2, Customized: false, PointSpread: 0, OfferedBetId: obid0, Settled: true})
	mustAddBet(&Bet{BetStatus: 0, Amount1: 12000, Outcome1Odds: -120, User1Agreed: true, User1Id: uid3, Amount2: 10000, Outcome2Odds: 120, User2Agreed: true, User2Id: uid1, Customized: false, PointSpread: 0, OfferedBetId: obid2, Settled: false})
	mustAddBet(&Bet{BetStatus: 0, Amount1: 1100, Outcome1Odds: -110, User1Agreed: true, User1Id: uid3, Amount2: 1000, Outcome2Odds: 110, User2Agreed: false, User2Id: uid5, Customized: false, PointSpread: 0, OfferedBetId: obid3, Settled: false})
	mustAddBet(&Bet{BetStatus: 2, Amount1: 1100, Outcome1Odds: -100, User1Agreed: true, User1Id: uid3, Amount2: 1100, Outcome2Odds: 100, User2Agreed: false, User2Id: uid5, Customized: false, PointSpread: 3.5, OfferedBetId: obid4, Settled: false})

	fmt.Println("DB seeded!")
}
