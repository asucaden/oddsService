package persistence

import "time"

type User struct {
	UserId   int `db:"user_id"`
	Username string
	Balance  float64
	Hash     string
}

type Bet struct {
	BetId        int `db:"bet_id"`
	BetStatus    int `db:"bet_status"`
	Amount1      int
	Outcome1Odds int
	User1Agreed  bool
	User1Id      int
	Amount2      int
	Outcome2Odds int
	User2Agreed  bool
	User2Id      int
	PointSpread  float32 `db:"point_spread"`
	Customized   bool
	OfferedBetId int `db:"offered_bet_id"`
	Settled      bool
}

type OfferedBet struct {
	OfferedBetId   int    `db:"offered_bet_id"`
	OfferedBetName string `db:"offered_bet_name"`
	Outcome1       string
	Outcome1Odds   int
	Outcome2       string
	Outcome2Odds   int
	PointSpread    float32   `db:"point_spread"`
	EventDate      time.Time `db:"event_date"`
	EventStatus    int       `db:"event_status"`
	CompetitionId  string    `db:"competition_id"`
}

type Competition struct {
	CompetitionId   string    `db:"competition_id"`
	CompetitionName string    `db:"competition_name"`
	EventStatus     int       `db:"event_status"`
	EventDate       time.Time `db:"event_date"`
}

// JOINED TYPES
type BetAndUser struct {
	BetId        int `db:"bet_id"`
	BetStatus    int `db:"bet_status"`
	Amount1      int
	Outcome1Odds int
	User1Agreed  bool
	User1Id      int
	Amount2      int
	Outcome2Odds int
	User2Agreed  bool
	User2Id      int
	PointSpread  float32 `db:"point_spread"`
	OfferedBetId int     `db:"offered_bet_id"`
	Settled      bool

	Username1 string
	Balance1  float64

	Username2 string
	Balance2  float64

	ObName    string `db:"offered_bet_name"`
	Outcome1  string
	Outcome2  string
	EventDate time.Time `db:"event_date"`
}
