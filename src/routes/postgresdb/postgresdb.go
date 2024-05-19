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
func Select(db *sql.DB, dbname string, schema string, table string, columns []string, sorted string, order string, limit int, filter []string) (map[string]string, error) {
    fmt.Println("Call postgresdb.Select()")
	var query string
    // columns = strings.Join(columns, ", ")
    if filter != nil {
        query = fmt.Sprintf(`SELECT %s FROM %s."%s"."%s" WHERE %s ORDER BY %s %s LIMIT %d`, strings.Join(columns, ", "), dbname, schema, table, strings.Join(filter, ","), sorted, order, limit)
    } else {
        query = fmt.Sprintf(`SELECT %s FROM %s."%s"."%s" ORDER BY %s %s LIMIT %d`, strings.Join(columns, ", "), dbname, schema, table, sorted, order, limit)
    }
    _, data, err := _get(db, query)
	if err != nil {
		fmt.Println("Error when call Select():", err)
		return nil, err
	}
	fmt.Println(query)
	fmt.Println("Query executed successfully")
	fmt.Println(data)
    return data, nil
}

// InsertData inserts data into the specified table in the PostgreSQL database.
func Insert(db *sql.DB, dbname string, schema string, table string, columns []string, values []string) error {
	fmt.Println("Call postgresdb.Insert()")

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
	fmt.Println("Call postgresdb._get()")
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

// ListSchemas lists all schemas in the PostgreSQL database.
func ListSchemas(db *sql.DB) ([]string, error) {
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

func ListTables(db *sql.DB, dbname string, schema string) ([]string, error) {
	var tables []string

	query := fmt.Sprintf(`SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'`, schema)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func FindDiff(db *sql.DB, tables []string) ([]string, error) {
	fmt.Println("Call postgresdb.FindDiff()")
	var query string
	var result []string

	// Создание строки с запросом
	for i, table := range tables {
		if i > 0 {
			query += " INTERSECT "
		}
		query += fmt.Sprintf(`SELECT datetime FROM logger."%s"`, table)
	}

	query = fmt.Sprintf(`SELECT all_dates.datetime FROM ( %s ) AS all_dates
		LEFT JOIN messenger."DOCUMENTS" ON all_dates.datetime = messenger."DOCUMENTS".event_dt
		WHERE messenger."DOCUMENTS".event_dt IS NULL`, query)

	// fmt.Print(query)
	
	// Выполнение запроса
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Обработка результатов
	for rows.Next() {
		var datetime string
		if err := rows.Scan(&datetime); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Find to work with datetime:", datetime)
		result = append(result, datetime)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return result, nil
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
