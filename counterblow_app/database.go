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

func database_connect() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close() // do not close, keep it opened

	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Successfully connected!")
	return nil
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
func database_addRule(r RoutingRule) error {
	fmt.Println("Entering database_addRule function")
	stmt, err := db.Prepare("INSERT INTO rules (rule_type, rule_ipaddr, rule_subnetmask, rule_servers, rule_source, rule_dest) VALUES($1, $2, $3, $4, $5, $6) RETURNING rule_id;")

	if err != nil {
		panic(err.Error())
	}

	res, err := stmt.Exec(r.rule_type, r.rule_ipaddr, r.rule_subnetmask, r.rule_servers, r.rule_source, r.rule_dest)

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Successfully inserted with id %s\n", res)
	return nil
}

func database_removeRule(ruleid int) error {
	stmt, err := db.Prepare("DELETE FROM rules WHERE rule_id = $1;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(ruleid)

	if err != nil {
		return err
	}

	return nil
}

func database_loadRules() ([]RoutingRule, error) {
	rows, err := db.Query("SELECT rule_id, rule_type, rule_ipaddr, rule_subnetmask, rule_servers, rule_source, rule_dest FROM rules;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rulesSlice []RoutingRule

	// append works on nil slices.

	for rows.Next() {
		var rule RoutingRule

		if err := rows.Scan(&rule.rule_id, &rule.rule_type, &rule.rule_ipaddr, &rule.rule_subnetmask, &rule.rule_servers, &rule.rule_source, &rule.rule_dest); err != nil {
			return nil, err
		}
		rulesSlice = append(rulesSlice, rule)
		fmt.Printf("Loaded rule %v\n", rule)
	}
	if err := rows.Err(); err != nil {
		panic(err.Error())
	}
	return rulesSlice, nil
}
