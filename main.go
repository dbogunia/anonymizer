package main

import (
	"bufio"
	"database/sql"
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
	selectQuery := "SELECT DISTINCT(" + columnName + ") FROM " + tableName
	results, err := db.Query(selectQuery)
	// in case of error we don't panic but log query to log output along with error
	if err != nil {
		logQueryError(selectQuery, err)
		return
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
		var newValue = anonymizeString(values[i])

		/* not needed anynmore
		// lets build random string of the same length as original value
		if isNumber(values[i]) {
			newValue = RandNumberRunes(len(values[i]))
		} else {
			newValue = RandStringRunes(len(values[i]))
		}
		*/
		// crerate UPDATE statement
		updateQuery := "UPDATE " +
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
		update, err := db.Query(updateQuery)

		// in case of error we don't panic but log query to log output along with error
		if err != nil {
			logQueryError(updateQuery, err)
		}

		// closing update object
		update.Close()
	}
}

func logQueryError(query string, err error) {
	log.Println("---------------------")
	log.Println("Error running query: ")
	log.Println(query)
	log.Println("Error message: " + err.Error())
}

func isNumber(str string) bool {
	if _, err := strconv.Atoi(str); err == nil {
		return true
	}
	return false
}

func anonymizeString(str string) string {
	retString := ""
	// exception for phone muber
	if strings.HasPrefix(str, "+") {
		retString = "+"
		str = str[1:len(str)]
	}
	// end of exceptions
	elements := strings.Split(str, " ")
	for _, element := range elements {
		if isNumber(element) {
			retString = strings.Join([]string{retString, RandNumberRunes(len(element))}, " ")
		} else {
			retString = strings.Join([]string{retString, RandStringRunes(len(element))}, " ")
		}
	}
	return retString
}

func loadTablesFromFile() []string {
	file, err := os.Open("./tables.txt")
	var tables []string
	if err != nil {
		// if we can't open file tables.txt we log error and exit
		log.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tables = append(tables, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		// if we can't read the tables.txt file we log error and exit
		log.Fatal(err)
	}

	return tables
}

func main() {
	if len(os.Args) < 2 {
		log.Println("Please provide connection string")
	} else {
		connString := os.Args[1]
		log.Println("Starting anonymizer v1.0.1")
		tables := loadTablesFromFile()
		log.Println("Opening database connection")
		db, err := sql.Open("mysql", connString)

		// If there is an error opening the connection, we log error and exit ;-)
		if err != nil {
			log.Fatal(err.Error())
		}

		defer db.Close()
		for _, table := range tables {
			log.Println("Processing column " + table)
			tableName := strings.Split(table, ".")[0]
			columnName := strings.Split(table, ".")[1]
			anonymize(tableName, columnName, *db)
		}
	}
}
