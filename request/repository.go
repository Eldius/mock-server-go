package request

import (
	"context"
	"database/sql"
	"encoding/json"
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
	REQ_ID varchar(50),
	REQ_DATE timestamp,
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
	, REQ_ID
	, REQ_DATE
) values (
	?
	, ?
	, ?
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

	selectRequests       = `SELECT ID, REQ_ID, REQ_DATE, PATH, METHOD, REQUEST, RESPONSE, RESPONSE_CODE FROM REQUEST`
	selectRequestHeaders = `
-- SQLite
SELECT
    NAME,
    VALUE
FROM
    HEADERS
WHERE
    REQUEST_ID = ?
    AND HEADER_TYPE = ?
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
			log.Printf("Failed to insert record\n%s", err.Error())
		}
		log.Println(result.RowsAffected())
		statement, err = db.Prepare(createTableHeader) // Prepare SQL Statement
		if err != nil {
			log.Fatal(err.Error())
		}
		result, err = statement.Exec() // Execute SQL Statements
		if err != nil {
			log.Printf("Failed to insert record\n%s", err.Error())
		}
		log.Println(result.RowsAffected())
		log.Println("request table created")
	} else {
		log.Fatalf("Failed to open database\n%s", err.Error())
	}

	return db
}

func Persist(r *Record) {
	debug(r)
	db := initDB()
	if result, err := db.Exec(insertRequest, r.Request.Path, r.Request.Method, r.Request.Body, r.Response.Body, r.Response.Code, r.ReqID, r.RequestDate); err != nil {
		log.Println("Failed to insert request to db")
		log.Fatal(err.Error())
	} else {
		log.Println(result)
		reqId, _ := result.LastInsertId()
		for k, v_ := range r.Request.Headers {
			for _, v := range v_ {

				if _, err = db.ExecContext(
					context.Background(),
					insertHeader,
					k,     // NAME
					v,     // VALUE
					reqId, // REQUEST_ID
					"IN",  // HEADER_TYPE
				); err != nil {
					log.Printf("Failed to insert request headers\n%s\n", err)
				}
			}
		}
		for k, v_ := range r.Response.Headers {
			for _, v := range v_ {
				if _, err = db.ExecContext(
					context.Background(),
					insertHeader,
					k,     // NAME
					v,     // VALUE
					reqId, // REQUEST_ID
					"OUT", // HEADER_TYPE
				); err != nil {
					log.Printf("Failed to insert request headers\n%s\n", err)
				}
			}
		}

	}
}

func GetRequests() []Record {
	records := make([]Record, 0)
	db := initDB()
	if row, err := db.Query(selectRequests); err != nil {
		log.Println("Failed to list requests")
		log.Printf("Failed to query requests\n%s\n", err.Error())
	} else {
		defer row.Close()
		for row.Next() { // Iterate and fetch the records from result cursor
			record := Record{
				Request:  RequestRecord{},
				Response: ResponseRecord{},
			}

			_ = row.Scan(
				&record.ID,             // ID
				&record.ReqID,          // REQ_ID
				&record.RequestDate,    // REQ_DATE
				&record.Request.Path,   // PATH
				&record.Request.Method, // METHOD
				&record.Request.Body,   // REQUEST
				&record.Response.Body,  // RESPONSE
				&record.Response.Code,  // RESPONSE_CODE
			)
			if reqHeadersRow, err := db.Query(selectRequestHeaders, record.ID, "IN"); err == nil {
				var reqHeaders Headers = make(Headers)
				for reqHeadersRow.Next() { // Iterate and fetch the records from result cursor
					var key, value string
					_ = reqHeadersRow.Scan(&key, &value)
					reqHeaders[key] = append(reqHeaders[key], value)
				}
				record.Request.Headers = reqHeaders
			} else {
				log.Printf("Failed to fetch request headers\n%s\n", err)
			}
			if resHeadersRow, err := db.Query(selectRequestHeaders, record.ID, "OUT"); err == nil {
				var resHeaders Headers = make(Headers)
				for resHeadersRow.Next() { // Iterate and fetch the records from result cursor
					var key, value string
					_ = resHeadersRow.Scan(&key, &value)
					resHeaders[key] = append(resHeaders[key], value)
				}
				record.Response.Headers = resHeaders
			} else {
				log.Printf("Failed to fetch response headers\n%s\n", err)
			}
			records = append(records, record)
		}
	}

	return records
}

func debug(obj interface{}) {
	dbg, _ := json.Marshal(obj)
	fmt.Println(string(dbg))
}
