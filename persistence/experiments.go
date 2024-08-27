package persistence

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Experiments(q Querier) {
	users, err := AllUsers(q)
	if err != nil {
		log.Fatal(err)
	}

	// Print out each user
	for _, u := range users {
		fmt.Println(u)
	}

	// Print out all the bets for user X
	bets, err := BetsByUser(q, 2)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println()
	for _, b := range bets {
		fmt.Println(b)
	}
}
