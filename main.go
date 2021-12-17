package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	// connect to a database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=test_connect user=test_connect password=secretpassword")
	if err != nil {
		log.Fatalf("unable to connect %v\n", err)
	}
	defer conn.Close()

	log.Println("connected to database")

	// test connection
	err = conn.Ping()
	if err != nil {
		log.Fatal("cannot ping database")
	}

	log.Println("pinged database")

	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// insert a row into table
	query := `insert into users (first_name, last_name) values ($1, $2)`
	_, err = conn.Exec(query, "Jack", "Brown")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("inserted a row")

	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// update a row
	stmt := `update users set first_name = $1 where id = $2`
	_, err = conn.Exec(stmt, "Jack", 5)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("updated one or more rows")

	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	// get one row by id
	query = `select id, first_name, last_name from users where id = $1`
	var firstName, lastName string
	var id int

	row := conn.QueryRow(query, 1)

	err = row.Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("queryrow returns", id, firstName, lastName)

	// delete a row
	query = `delete from users where id = $1`
	_, err = conn.Exec(query, 6)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("deleted a row")

	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("select id, first_name, last_name from users")
	if err != nil {
		log.Println(err)
		return err
	}

	// close the connection if the query returns more than one row
	defer rows.Close()

	var firstName, lastName string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("record is", id, firstName, lastName)
	}

	// once we are done scanning rows let's make sure that there are no other errors
	if err = rows.Err(); err != nil {
		log.Fatal("error scanning rows")
	}

	fmt.Println("-------------------------")

	return nil
}
