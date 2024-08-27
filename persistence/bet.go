package persistence

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// INSERT
func AddBet(q Querier, bet *Bet) (int, error) {
	var id int

	row := q.QueryRowx("INSERT INTO bet (bet_status, amount1, outcome1odds, user1agreed, user1id, amount2, outcome2odds, user2agreed, user2id, offered_bet_id, customized, settled, point_spread)"+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)"+
		"RETURNING bet_id",
		bet.BetStatus, bet.Amount1, bet.Outcome1Odds, bet.User1Agreed, bet.User1Id, bet.Amount2, bet.Outcome2Odds, bet.User2Agreed, bet.User2Id, bet.OfferedBetId, bet.Customized, bet.Settled, bet.PointSpread)

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("addBet: %v", err)
	}

	return id, nil
}

// UPDATE
func AgreeToBet(q Querier, betId int, userN int, agreed bool) (int, error) {
	var id int
	var userNagreed string
	var statusUpdate string

	switch userN {
	case 1:
		userNagreed = "user1agreed"
	case 2:
		userNagreed = "user2agreed"
	default:
		return 0, fmt.Errorf("AgreeToBet: Invalid userN argument")
	}

	if !agreed {
		statusUpdate = ", bet_status = 4 " // TODO formalize bet_status enumeration. using 4 to mean 'declined' here
	}

	row := q.QueryRowx("UPDATE bet "+
		"SET "+userNagreed+" = $1 "+statusUpdate+
		"WHERE bet_id = $2 "+
		"RETURNING bet_id",
		agreed, betId)

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("AgreeToBet: %v", err)
	}
	return id, nil
}

func UpdateBet(q Querier, bet *Bet) (int, error) {
	var id int

	row := q.QueryRowx("UPDATE bet "+
		"SET bet_status = $1, "+
		"amount1 = $2, outcome1odds = $3, user1agreed = $4, user1id = $5, "+
		"amount2 = $6, outcome2odds = $7, user2agreed = $8, user2id = $9, "+
		"offered_bet_id = $10, customized = $11, settled = $12, point_spread = $13 "+
		"WHERE bet_id = $14 "+
		"RETURNING bet_id",
		bet.BetStatus, bet.Amount1, bet.Outcome1Odds, bet.User1Agreed, bet.User1Id, bet.Amount2, bet.Outcome2Odds, bet.User2Agreed, bet.User2Id, bet.OfferedBetId, bet.Customized, bet.Settled, bet.PointSpread, bet.BetId)

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("UpdateBet: %v", err)
	}

	return id, nil
}

// func VoidBet(q Querier, bet *Bet) (int, error) {
// }

// func UpdateAmountBet(q Querier, bet *Bet) (int, error) {
// }

// SELECT

// many

func BetsByUser(q Querier, user_id int) ([]Bet, error) {
	var bets []Bet

	rows, err := q.Query("SELECT * FROM bet WHERE user1id = $1 OR user2id = $1", user_id)
	if err != nil {
		return nil, fmt.Errorf("\n\n\nERROR running the query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bet Bet
		if err := rows.Scan(&bet.BetId, &bet.BetStatus,
			&bet.Amount1, &bet.Outcome1Odds, &bet.User1Agreed, &bet.User1Id,
			&bet.Amount2, &bet.Outcome2Odds, &bet.User2Agreed, &bet.User2Id,
			&bet.OfferedBetId, &bet.Settled); err != nil {
			return nil, fmt.Errorf("bets by user: %v", err)
		}
		bets = append(bets, bet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("bets by user: %v", err)
	}
	return bets, nil
}

func BetsByOfferedBet(q Querier, offeredBetId int) ([]Bet, error) {
	var bets []Bet

	err := q.Select(&bets, "SELECT * FROM bet WHERE offered_bet_id = $1", offeredBetId)
	if err != nil {
		return nil, fmt.Errorf("BetsByOfferedBet %v", err)
	}
	return bets, nil
}

// joins
func BetAndUsersByBetId(q Querier, betId int) (*BetAndUser, error) {
	var bet BetAndUser
	err := q.Get(&bet, "SELECT bet_id, bet_status, amount1, point_spread, bet.outcome1odds, user1agreed, "+
		"amount2, bet.outcome2odds, user2agreed, settled, offered_bet_id, "+
		"user1id, u1.username AS username1, u1.balance AS balance1, "+
		"user2id, u2.username AS username2, u2.balance AS balance2 "+
		"FROM bet "+
		"INNER JOIN users u1 ON bet.user1id = u1.user_id "+
		"INNER JOIN users u2 ON bet.user2id = u2.user_id "+
		"WHERE bet.bet_id = $1", betId)

	if err != nil {
		return nil, fmt.Errorf("BetAndUsersByBetId(): %v", err)
	}
	return &bet, nil
}

func BetsAndObsAndUsersByUserId(q Querier, userId int) ([]BetAndUser, error) {
	var bets []BetAndUser
	err := q.Select(&bets, "SELECT bet_id, bet_status, amount1, bet.point_spread AS point_spread, bet.outcome1odds, user1agreed, "+
		"amount2, bet.outcome2odds, user2agreed, settled, "+
		"user1id, u1.username AS username1, u1.balance AS balance1, "+
		"user2id, u2.username AS username2, u2.balance AS balance2, "+
		"ob.offered_bet_name, ob.outcome1, ob.outcome2, ob.event_date "+
		"FROM bet "+
		"INNER JOIN users u1 ON bet.user1id = u1.user_id "+
		"INNER JOIN users u2 ON bet.user2id = u2.user_id "+
		"INNER JOIN offered_bet ob ON bet.offered_bet_id = ob.offered_bet_id "+
		"WHERE bet.user1id = $1 OR bet.user2id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("BetsAndObsUsersByUserId(): %v", err)
	}
	return bets, nil
}

// one

func OneBet(q Querier, betId int) (Bet, error) {
	var bet Bet

	err := q.Get(&bet, "SELECT * FROM bet WHERE bet_id = $1", betId)
	if err != nil {
		return bet, fmt.Errorf("OneOfferedBet %d: %v", betId, err)
	}
	return bet, nil
}
