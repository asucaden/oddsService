package persistence

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// INSERT
func AddUser(db *sqlx.DB, user *User) (int, error) {
	var id int

	row := db.QueryRow("INSERT INTO users (username, balance, hash)"+
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
func AllUsers(db *sql.DB) ([]User, error) {
	var users []User
	rows, err := db.Query("SELECT * FROM users")
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

func OneUser(db *sqlx.DB, userId int) (User, error) {
	var user User

	err := db.Get(&user, "SELECT * FROM users WHERE user_id = $1", userId)
	if err != nil {
		return user, fmt.Errorf("OneCompetition %d: %v", userId, err)
	}
	return user, nil
}

// case insensitive
func OneUserByName(db *sqlx.DB, username string) (User, error) {
	var user User

	err := db.Get(&user, "SELECT * FROM users WHERE username ILIKE $1", username)
	if err != nil {
		return user, fmt.Errorf("OneUserByName %s: %v", username, err)
	}
	return user, nil
}

// one partial

func OneUsername(db *sqlx.DB, userId int) (string, error) {
	var username string

	err := db.Get(&username, "SELECT username FROM users WHERE user_id = $1", userId)
	if err != nil {
		return username, fmt.Errorf("OneCompetition %s: %v", username, err)
	}
	return username, nil
}

// case insensitive
func OneUserIdByName(db *sqlx.DB, username string) (int, error) {
	var userId int

	err := db.Get(&userId, "SELECT user_id FROM users WHERE username ILIKE $1", username)
	if err != nil {
		return userId, fmt.Errorf("OneCompetition %s: %v", username, err)
	}
	return userId, nil
}

// case insensitive
func OneHashByName(db *sqlx.DB, username string) (string, error) {
	var hash string

	err := db.Get(&hash, "SELECT hash FROM users WHERE username ILIKE $1", username)
	if err != nil {
		return hash, fmt.Errorf("OneHashByName %s: %v", username, err)
	}
	return hash, nil
}

// counts
func CountUserWinLoss(db *sqlx.DB, userId int) (int, int, error) {
	var betsWon int
	var betsLost int

	err := db.Get(&betsWon, "SELECT COUNT(bet_id) FROM bet "+
		"WHERE (bet_status = 1 AND user1id = $1) "+
		"OR (bet_status = 2 AND user2id = $1)", userId)
	if err != nil {
		return 0, 0, fmt.Errorf("UserWinLoss: %s", err.Error())
	}

	err = db.Get(&betsLost, "SELECT COUNT(bet_id) FROM bet "+
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
func CheckUsername(db *sqlx.DB, username string) (bool, error) {
	row := db.QueryRow("SELECT hash FROM users WHERE username ILIKE $1", username)
	if err := row.Scan(); err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}
	return false, nil

}
