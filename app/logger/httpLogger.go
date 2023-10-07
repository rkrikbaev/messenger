package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"

    _ "github.com/mattn/go-sqlite3"

    sql_handler "app/logger/sqlHandler"
)

const (
    databaseFile = "data.db"
    host = "192.168.1.109"
)

type Data struct {
    Value     float64 `json:"Value"`
    Quality   int     `json:"Quality"`
    Timestamp string  `json:"Timestamp"`
}

func LoadDataFromHTTP() {
    // Initialize the SQLite database
    // db, err := sql.Open("sqlite3", databaseFile)
    // checkErr(err)
    // defer db.Close()

    // Fetch data from a web service
    url := fmt.Printf("https://%s:4444/tags", host) // Replace with your web service URL
    response, err := http.Get(url)
    checkErr(err)
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        log.Fatalf("HTTP GET request failed with status: %v", response.Status)
    }

    // Read the response body
    body, err := ioutil.ReadAll(response.Body)
    checkErr(err)

    // Parse the JSON response
    var jsonData map[string]Data
    err = json.Unmarshal(body, &jsonData)
    checkErr(err)

    // Insert data into the SQLite database
    for key, data := range jsonData {
        insertData := "INSERT INTO data (value, quality, timestamp) VALUES (?, ?, ?)"
        _, err = db.Exec(insertData, data.Value, data.Quality, data.Timestamp)
        if err != nil {
            log.Printf("Failed to insert data for key %s: %v", key, err)
        }
    }

    sql_handler.LogData(eventid, columns, values, "")

    fmt.Println("Data saved to SQLite database successfully!")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}