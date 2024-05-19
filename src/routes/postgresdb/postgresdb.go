package postgresdb

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

var (
	EventTimeFormat = "2006-01-02T15:04:05.000-07:00"
)

func Update(db *sql.DB, dbname string, schema string, table string, columns, values []string, filter, state string) error {

	var columnValuePairs []string

	for i := range columns {
		columnValuePairs = append(columnValuePairs, fmt.Sprintf(`%s = '%s'`, columns[i], values[i]))
	}
	// Construct the SQL query for updating the record
	query := fmt.Sprintf(`UPDATE %s."%s"."%s" SET %s WHERE (%s)`, dbname, schema, table, strings.Join(columnValuePairs, ","), filter)
	err := _put(db, query, values)
	if err != nil {
		return err
	}
	return nil
}

// GetRecord fetches the last or first record from the specified table in the PostgreSQL database based on the provided parameters.
func Select(db *sql.DB, dbname string, schema string, table string, columns []string, orderBy string, limit int, filter []string) (map[string]string, error) {
    var query string
    columnsStr := strings.Join(columns, ", ")
    if filter != nil {
        query = fmt.Sprintf(`SELECT %s FROM %s."%s"."%s" WHERE %s ORDER BY datetime %s LIMIT %d`, columnsStr, dbname, schema, table, strings.Join(filter, ","), orderBy, limit)
    } else {
        query = fmt.Sprintf(`SELECT %s FROM %s."%s"."%s" ORDER BY datetime %s LIMIT %d`, columnsStr, dbname, schema, table, orderBy, limit)
    }
    fmt.Println(query)
    _, data, err := _get(db, query)
    if err != nil {
        return nil, err
    }
    fmt.Println(data)
    return data, nil
}

// InsertData inserts data into the specified table in the PostgreSQL database.
func Insert(db *sql.DB, dbname string, schema string, table string, columns []string, values []string) error {
	// Add the createdAtDate column to the EventParams slice
	// columns := EventParams
	// columns = append(columns, "collected")
	// collected := time.Now().Format(EventTimeFormat)
	// fmt.Println(collected)
	// values = append(values, collected)

	// if table == "DOCUMENTS" {
	// 	columns = append(columns, "state")
	// 	values = append(values, state)
	// }
	// Construct the SQL query for insertion
	placeholders := make([]string, len(columns))
	for i := range columns {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	valuePlaceholders := strings.Join(placeholders, ", ")

	columnsNames := strings.Join(columns, ",")

	// query := fmt.Sprintf(`INSERT INTO %s."%s" (%s) VALUES (%s)`, schema, table, EventParamsString, valuePlaceholders)
	query := fmt.Sprintf(`INSERT INTO %s."%s"."%s" (%s) VALUES (%s) ON CONFLICT DO NOTHING`, dbname, schema, table, columnsNames, valuePlaceholders)

	err := _put(db, query, values)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// GetData fetches data from the specified table in the PostgreSQL database and processes it.
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

	fmt.Println(query)
	fmt.Println(values)

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

// ListSchemas lists all schemas in the PostgreSQL database.
func listSchemas(db *sql.DB) ([]string, error) {
	var schemas []string

	rows, err := db.Query(`SELECT schema_name FROM information_schema.schemata`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return nil, err
		}
		schemas = append(schemas, schemaName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schemas, nil
}

func Connect(dbname string, user string, password string, host string, port string) (*sql.DB, error) {
	connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable", dbname, user, password, host, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to the database")

	return db, nil
}


// func InitDB() (*sql.DB, error) {
// 	// Connect to the PostgreSQL database
// 	db, err := Connect(DB_NAME, DB_USER, DB_PASSWORD, DB_HOST, DB_PORT)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check if the database exists
// 	dbExists := checkDatabaseExists(DB_NAME, DB_USER)

// 	// If the database doesn't exist, create it
// 	if !dbExists {
// 		fmt.Printf("Database %s does not exist, creating.\n", DB_NAME)
// 		createDatabase(DB_NAME, DB_USER)
// 		executeSQLFile(DB_NAME, DB_USER, DB_FILE_PATH)
// 	} else {
// 		fmt.Printf("Database %s already exists, skipping creation.\n", DB_NAME)
// 	}

// 	return db, nil
// }

// if err != nil {
// // Check if the database exists
// dbExists := checkDatabaseExists(dsn, DB_NAME, DB_USER)

// // If the database doesn't exist, create it
// if !dbExists {
// 	fmt.Printf("Database %s does not exist, creating.\n", DB_NAME)
// 	createDatabase(DB_NAME, DB_USER)
// 	executeSQLFile(DB_NAME, DB_USER, DB_FILE_PATH)
// } else {
// 	fmt.Printf("Database %s already exists, skipping creation.\n", DB_NAME)
// }
// }

// if err != nil {
// log.Fatal(err)
// }

// defer db.Close()
