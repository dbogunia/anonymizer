package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TableDescription struct {
	Field   sql.NullString `json:"field"`
	Type    sql.NullString `json:"type"`
	Null    sql.NullString `json:"null"`
	Key     sql.NullString `json:"key"`
	Default sql.NullString `json:"default"`
	Extra   sql.NullString `json:"extra"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("0123456789")

func RandNumberRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numberRunes[rand.Intn(len(numberRunes))]
	}
	return string(b)
}
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func anonymize(tableName string, columnName string, db sql.DB) {

	log.Println("Working on " + tableName + "." + columnName)

	// create table for values
	var values []string

	// Get distinct values from column
	results, err := db.Query("SELECT DISTINCT(" + columnName + ") FROM " + tableName)
	if err != nil {
		panic(err.Error())
	}

	// load values into slice
	// omit nulls

	for results.Next() {
		// we use nullable string
		var nullStr sql.NullString
		results.Scan((&nullStr))
		// add to values if string is not null (string is valid)
		if len(nullStr.String) > 0 {
			values = append(values, nullStr.String)
		}
	}

	// anonymize values
	for i := 0; i < len(values); i++ {
		var newValue = ""
		// lets build random string of the same length as original value
		if isNumber(values[i]) {
			newValue = RandNumberRunes(len(values[i]))
		} else {
			newValue = RandStringRunes(len(values[i]))
		}

		// crerate UPDATE statement
		updStr := "UPDATE " +
			tableName +
			" SET " +
			columnName +
			" = '" +
			newValue +
			"' WHERE " +
			columnName +
			" = '" +
			values[i] +
			"'"

		log.Println("Anonymizing: " + values[i] + " to: " + newValue)
		// and execute it
		update, err := db.Query(updStr)
		// in case of error we panic ;-)
		if err != nil {
			panic(err.Error())
		}

		// closing update object
		update.Close()

	}
}

func isNumber(str string) bool {
	if _, err := strconv.Atoi(str); err == nil {
		return true
	}
	return false
}

func loadTablesFromFile() []string {
	file, err := os.Open("./tables.txt")
	var tables []string

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		tables = append(tables, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return tables
}

func main() {

	log.Println("Starting anonymizer")
	tables := loadTablesFromFile()

	log.Println("Opening database connection")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/hybris")

	// If there is an error opening the connection, we panic ;-)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	for _, table := range tables {

		log.Println("Processing column " + table)

		tableName := strings.Split(table, ".")[0]
		columnName := strings.Split(table, ".")[1]

		anonymize(tableName, columnName, *db)

	}

}
