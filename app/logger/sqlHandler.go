package main

import (
	"fmt"
)

func InsertData(table string, columns []string, values []string, state string) error {

	columns = append(columns, "createdAtDate")
	values = append(values, time.Now().Format(EventTimeFormat))

	// Construct the SQL query for insertion
	placeholders := strings.Repeat("?, ", len(columns)-1) + "?" // Repeat the "?" for placeholders
	params := strings.Join(columns, ",")
	query := fmt.Sprintf("INSERT OR IGNORE INTO '%s' (%s) VALUES (%s)", table, params, placeholders)
	fmt.Println(query)
    err = execQuery(db, query, values)
	checkErr(err)
	return nil
}

// InsertData inserts data into the specified table in the SQLite database.
func execQuery(db *sql.DB, query string, values []string) error {
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

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}