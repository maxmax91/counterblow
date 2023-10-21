package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// as per docker-compose.yml
const (
	host     = "localhost"
	port     = 5438
	user     = "counterblow_user"
	password = "postgres123!?"
	dbname   = "counterblow_db"
)

var db *sql.DB
var err error

func connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

}

func addHit() error {
	fmt.Println("test")
	stmt, err := db.Prepare("INSERT INTO hits (hit_from, hit_to) VALUES(?,?) RETURNING hit_datetime")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec("test", "test@mail.com")

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)

	fmt.Println("Successfully inserted!")
	return nil
}
