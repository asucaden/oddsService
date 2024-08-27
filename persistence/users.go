package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// INSERT
func AddUser(q Querier, user *User) (int, error) {
	var id int

	row := q.QueryRowx("INSERT INTO users (username, balance, hash)"+
		"VALUES ($1, $2, $3)"+
		"RETURNING user_id",
		user.Username, user.Balance, user.Hash)

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("addBet: %v", err)
	}

	return id, nil
}

// SELECT

// many
func AllUsers(q Querier) ([]User, error) {
	var users []User
	rows, err := q.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("ERROR running the query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserId, &user.Username, &user.Balance); err != nil {
			return nil, fmt.Errorf("allUsers: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("allUsers: %v", err)
	}
	return users, nil
}

// one

func OneUser(q Querier, userId int) (User, error) {
	var user User

	err := q.Get(&user, "SELECT * FROM users WHERE user_id = $1", userId)
	if err != nil {
		return user, fmt.Errorf("OneCompetition %d: %v", userId, err)
	}
	return user, nil
}

// case insensitive
func OneUserByName(q Querier, username string) (User, error) {
	var user User

	err := q.Get(&user, "SELECT * FROM users WHERE username ILIKE $1", username)
	if err != nil {
		return user, fmt.Errorf("OneUserByName %s: %v", username, err)
	}
	return user, nil
}

// one partial

func OneUsername(q Querier, userId int) (string, error) {
	var username string

	err := q.Get(&username, "SELECT username FROM users WHERE user_id = $1", userId)
	if err != nil {
		return username, fmt.Errorf("OneCompetition %s: %v", username, err)
	}
	return username, nil
}

// case insensitive
func OneUserIdByName(q Querier, username string) (int, error) {
	var userId int

	err := q.Get(&userId, "SELECT user_id FROM users WHERE username ILIKE $1", username)
	if err != nil {
		return userId, fmt.Errorf("OneCompetition %s: %v", username, err)
	}
	return userId, nil
}

// case insensitive
func OneHashByName(q Querier, username string) (string, error) {
	var hash string

	err := q.Get(&hash, "SELECT hash FROM users WHERE username ILIKE $1", username)
	if err != nil {
		return hash, fmt.Errorf("OneHashByName %s: %v", username, err)
	}
	return hash, nil
}

// counts
func CountUserWinLoss(q Querier, userId int) (int, int, error) {
	var betsWon int
	var betsLost int

	err := q.Get(&betsWon, "SELECT COUNT(bet_id) FROM bet "+
		"WHERE (bet_status = 1 AND user1id = $1) "+
		"OR (bet_status = 2 AND user2id = $1)", userId)
	if err != nil {
		return 0, 0, fmt.Errorf("UserWinLoss: %s", err.Error())
	}

	err = q.Get(&betsLost, "SELECT COUNT(bet_id) FROM bet "+
		"WHERE (bet_status = 1 AND user2id = $1) "+
		"OR (bet_status = 2 AND user1id = $1)", userId)
	if err != nil {
		return 0, 0, fmt.Errorf("UserWinLoss: %s", err.Error())
	}

	return betsWon, betsLost, nil

}

// check

// case insensitive
// true means username can be used, false means it cannot be used
func CheckUsername(q Querier, username string) (bool, error) {
	row := q.QueryRowx("SELECT hash FROM users WHERE username ILIKE $1", username)
	if err := row.Scan(); err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}
	return false, nil

}
