package dbsql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)


import (
	"log"
)

var (
	EventTimeFormat = "2006-01-02T15:04:05.000-07:00"
)

func Update(db *sql.DB, table, columns, values []string, condition, state string) error {

	var columnValuePairs []string
	
	for i := range columns {
		columnValuePairs = append(columnValuePairs, fmt.Sprintf(`%s = '%s'`, columns[i], values[i]))
	}
	// Construct the SQL query for updating the record
	query := fmt.Sprintf(`UPDATE '%s' SET %s WHERE (%s)`, table, strings.Join(columnValuePairs, ","), condition)
	err := _put(db, query, values)
	if err != nil {
		return err
	}
	return nil
}

// GetRecord fetches the last or first record from the specified table in the SQLite database based on the provided parameters.
func Select(db *sql.DB, table string, orderBy string, limit int, filter []string) (map[string]string, error) {
	var query string
	if filter != nil {
		query = fmt.Sprintf(`SELECT * FROM '%s' WHERE %s ORDER BY datetime %s LIMIT %d`, table, strings.Join(filter, ","), orderBy, limit)
	} else {
		query = fmt.Sprintf(`SELECT * FROM '%s' ORDER BY datetime %s LIMIT %d`, table, orderBy, limit)
	}
	fmt.Println(query)
	_, data, err := _get(db, query)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	return data, nil
}

// InsertData inserts data into the specified table in the SQLite database.
func Insert(db *sql.DB, table string, EventParams []string, values []string, state string) error {
	// Add the createdAtDate column to the EventParams slice
	columns := EventParams
	columns = append(columns, "createdAtDate")
	values = append(values, time.Now().Format(EventTimeFormat))

	if table == "DOCUMENTS" {
		columns = append(columns, "state")
		values = append(values, state)
	}
	// Construct the SQL query for insertion
	valuePlaceholders := strings.Repeat("?, ", len(columns)-1) + "?" // Repeat the "?" for placeholders
	EventParamsString := strings.Join(columns, ",")

	query := fmt.Sprintf("INSERT OR IGNORE INTO '%s' (%s) VALUES (%s)", table, EventParamsString, valuePlaceholders)

	fmt.Println(query)
	fmt.Println(values)

	err := _put(db, query, values)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// GetData fetches data from the specified table in the SQLite database and processes it.
// It returns true if the operation is successful.
func _get(db *sql.DB, query string) ([]string, map[string]string, error) {
    rows, err := db.Query(query)
    if err != nil {
        log.Fatalf("Error querying data: %s", err)
    }
	defer rows.Close()
	// Check if there is a row to scan
	if !rows.Next() {
		return nil, nil, nil
	}
	// Create a map to hold the scanned data
	scannedData := make(map[string]string)
	columns, _ := rows.Columns()
	// Scan the row into the scannedData map
	values := make([]interface{}, len(columns))
	value := make([]interface{}, len(columns))
	for i := range columns {
		values[i] = &value[i]
	}
	err = rows.Scan(values...)
	if err != nil {
		log.Printf("Error scanning row: %s", err)
		return nil, nil, err
	}
	// Populate the scanned data into the map
	for i, column := range columns {
		scannedData[column] = fmt.Sprint(*values[i].(*interface{}))
	}
	return columns, scannedData, nil
}

// putData executes a SQL query with the given values in a transaction.
// It takes a database connection, a query string, and a slice of values as parameters.
// The values are used to replace placeholders in the query.
// If an error occurs during the transaction, it will be rolled back.
// Returns an error if there was an issue with the transaction or query execution.
func _put(db *sql.DB, query string, values []string) error {

	data := make([]interface{}, len(values))

	for i, v := range values {
		data[i] = v
	}

	tx, err := db.Begin() // Start a transaction

	if err != nil {
        return fmt.Errorf("error beginning transaction: %w", err)
    }
    
    defer tx.Rollback() // Rollback the transaction in case of errors
    
    _, err = tx.Exec(query, data...)
    if err != nil {
        return fmt.Errorf("error executing query: %w", err)
    }
    
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("error committing transaction: %w", err)
    }
	return nil
}