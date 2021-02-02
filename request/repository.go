package request

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

const (
	createTableRequest = `
create table if not exists REQUEST (
	ID integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	PATH varchar(255),
	METHOD varchar(15),
	REQUEST varchar(4000),
	RESPONSE varchar(4000),
	RESPONSE_CODE int
);`

	createTableHeader = `
create table if not exists HEADERS (
	ID integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	NAME varchar(50),
	VALUE varchar(50),
	REQUEST_ID integer not null,
	HEADER_TYPE varchar(10) 
);
`
	insertRequest = `
insert into REQUEST (
	PATH
	, METHOD
	, REQUEST
	, RESPONSE
	, RESPONSE_CODE
) values (
	?
	, ?
	, ?
	, ?
	, ?
)
`
	insertHeader = `
insert into HEADERS (
	NAME
	, VALUE
	, REQUEST_ID
	, HEADER_TYPE
) values (
	?
	, ?
	, ?
	, ?
)
`
)

func initDB() *sql.DB {
	if db == nil {
		openDB()
	} else {
		if err := db.Ping(); err != nil {
			openDB()
		}
	}

	return db
}

func openDB() *sql.DB {
	var err error
	if db, err = sql.Open("sqlite3", "./mock.db"); err == nil {
		log.Println("Create request table...")
		statement, err := db.Prepare(createTableRequest) // Prepare SQL Statement
		if err != nil {
			log.Fatal(err.Error())
		}
		result, err := statement.Exec() // Execute SQL Statements
		if err != nil {
			fmt.Printf("Failed to insert record\n%s", err.Error())
		}
		log.Println(result.RowsAffected())
		statement, err = db.Prepare(createTableHeader) // Prepare SQL Statement
		if err != nil {
			log.Fatal(err.Error())
		}
		result, err = statement.Exec() // Execute SQL Statements
		if err != nil {
			fmt.Printf("Failed to insert record\n%s", err.Error())
		}
		log.Println(result.RowsAffected())
		log.Println("request table created")
	} else {
		log.Fatalf("Failed to open database\n%s", err.Error())
	}

	return db
}

func Persist(r *Record) *Record {
	db := initDB()
	if result, err := db.Exec(insertRequest, r.Request.Path, r.Request.Method, r.Request.Body, r.Response.Body, r.Response.Code); err != nil {
		fmt.Println("Failed to insert request to db")
		log.Fatal(err.Error())
	} else {
		fmt.Println(result)
		reqId, _ := result.LastInsertId()
		for k, v := range r.Request.Headers {
			_, _ = db.Exec(
				insertHeader,
				k,     // NAME
				v,     // VALUE
				reqId, // REQUEST_ID
				"IN",  // HEADER_TYPE
			)
		}
		for k, v := range r.Response.Headers {
			_, _ = db.Exec(
				insertHeader,
				k,     // NAME
				v,     // VALUE
				reqId, // REQUEST_ID
				"OUT", // HEADER_TYPE
			)
		}

	}
	return r
}
