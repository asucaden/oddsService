package persistence

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectDb() *sqlx.DB {
	user := os.Getenv("POSTGRES_USER")
	host := os.Getenv("POSTGRES_HOST")
	password := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("POSTGRES_PORT")
	drivername := "pgx"
	database := "bets"
	sslmode := "disable"

	db, err := sqlx.Connect(drivername, fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=%s", user, password, host, port, database, sslmode))
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to DB!")

	return db
}
