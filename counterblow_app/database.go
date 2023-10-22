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

func database_connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close() // do not close, keep it opened

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

}

func database_addHit(from string, to string) error {
	fmt.Println("Entering addHit function")
	stmt, err := db.Prepare("INSERT INTO hits (hit_from, hit_to) VALUES($1, $2) RETURNING hit_datetime;")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(from, to)

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Successfully inserted with datetime %s\n", res)
	return nil
}
func database_addRule(from string, to string) error {
	// todo
	return nil
}
func database_removeRule(from string, to string) error {
	// todo
	return nil
}

func database_loadRules() {
	rows, err := db.Query("SELECT rule_id, rule_type, rule_ipaddr, rule_subnetmask, rule_servers FROM rules;")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var rule RoutingRule

		if err := rows.Scan(&rule.rule_id, &rule.rule_type, &rule.rule_ipaddr, &rule.rule_subnetmask, &rule.rule_servers); err != nil {
			panic(err.Error())
		}
		fmt.Printf("Loaded rule %f", rule)
	}
	if err := rows.Err(); err != nil {
		panic(err.Error())
	}
}
