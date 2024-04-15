package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	basedd := "./bdd.sqlite3"
	theDB, err := sql.Open("sqlite3", basedd)

	if err != nil {
		log.Fatal(err)
	}
	rows, _ := theDB.Query("SELECT pseudo FROM USER")

	var pseudo string

	rows.Scan(pseudo)
	fmt.Println(pseudo)

	defer theDB.Close()

}
