package data

import (
	"time"

	"github.com/asucaden/goBet/persistence"
)

type OfferedBetView struct {
	OfferedBet persistence.OfferedBet
	BetCount   int
	BetSum     int
	Username   string
}

type PlaceBetView struct {
	OfferedBet persistence.OfferedBet
	Username   string
}

type CompetitionHeadline struct {
	Id    string
	Title string
	Date  time.Time
}

type CompetitionHeadlineView struct {
	ActiveHeadlines []CompetitionHeadline
	PastHeadlines   []CompetitionHeadline
	Username        string
}

type CompetitionView struct {
	Competition persistence.Competition
	OfferedBets []persistence.OfferedBet
	Username    string
}

type CheckUserView struct {
	Username string
	Exists   bool
}

type BetAgreementView struct {
	Success      bool
	Agreed       bool
	ErrorMessage string
}

type PreviousBet struct {
	BetId            int
	UserAmount       int
	UserOutcome      string
	UserOutcomeOdds  int
	OtherAmount      int
	OtherUsername    string
	OtherOutcome     string
	OtherOutcomeOdds int
	TotalAmount      int
}

type CounterofferView struct {
	Username    string
	OfferedBet  persistence.OfferedBet
	PreviousBet PreviousBet
}

type UserView struct {
	Username   string
	UserId     int
	Balance    float64
	BetsWon    int
	BetsLost   int
	ActiveBets []UserBetView
	PastBets   []UserBetView
}

type ProfileView struct {
	Username     string
	UserId       int
	Balance      float64
	BetsWon      int
	BetsLost     int
	ActiveBets   []UserBetView
	PastBets     []UserBetView
	IncomingBets []UserBetView
	OutgoingBets []UserBetView
}

type UserBetView struct {
	BetId            int
	BetStatus        int
	BetName          string
	UserAmount       int
	UserOutcome      string
	UserOutcomeOdds  int
	UserAgreed       bool
	UserName         string
	UserWon          bool
	UserPointSpread  float32
	OtherAmount      int
	OtherOutcome     string
	OtherOutcomeOdds int
	OtherAgreed      bool
	OtherName        string
	OtherPointSpread float32
	TotalAmount      int
	EventDate        time.Time
}
