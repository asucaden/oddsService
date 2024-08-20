package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Experiments(db *sql.DB) {
	users, err := AllUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	// Print out each user
	for _, u := range users {
		fmt.Println(u)
	}

	// Print out all the bets for user X
	bets, err := BetsByUser(db, 2)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println()
	for _, b := range bets {
		fmt.Println(b)
	}
}
