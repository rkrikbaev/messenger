package main

import (
	// "database/sql"
	"database/sql"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	// "strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Type1Events = []string{"407001","407002","407003"}
	Type2Events = []string{"407004","407005","407006","407007","407008","407009"}
	Type1Params = []string{"datetime","density","massflowbegin","massflowend","mass"}
	Type2Params = []string{"datetime","density","volume","temperature","tankLevel","mass"}
)

//добавить создание запись в таблице REPORT
func LogDataFromFile(objects map[string]interface{}) error {
	for objectName := range objects {
		fullDirPath := path.Join(FilesDirPath, objectName)
		files, err := listFilesInDirectory(fullDirPath, extension)
		if err != nil {
			log.Fatalf("Error listing files in directory: %s", err)
		}
		for _, filename := range files {
			destPath := fmt.Sprintf("%s/saved/%s",fullDirPath,filename)
			TableName := strings.Split(filename, "_")[0]
			sourcePath := fmt.Sprintf("%s/%s", fullDirPath, filename)
			// fmt.Println("Found path to file:", fullFilePath)
			data, err := ReadCSVFile(sourcePath)
			if err != nil {
				log.Fatalf("Error reading CSV file: %s", err)
			}
			var EventParams []string
			if findString(Type1Events, TableName) {
				EventParams = Type1Params
			} else if findString(Type2Events, TableName)  { 
				EventParams = Type2Params
			}
			for _, LoggedData := range data {
				LastSaveData, _ := getLastRecord(TableName)
				LoggedDataTime, _ := time.Parse(CSVTimeFormat, LoggedData[0])
				LastSavedEventTime := LastSaveData["datetime"]
				ParsedEventTime := ConverToEventDate(EventTimeFormat, LoggedDataTime)
				LoggedData[0] = ParsedEventTime
				if len(LastSaveData) == 0 {
					_ = LogData(TableName, EventParams, LoggedData)
				} else {
					if isAfter, err := date1_after_date2(ParsedEventTime, LastSavedEventTime); err != nil {
						log.Fatalf("Error comparing dates: %s", err)
					} else if isAfter {
						err = LogData(TableName, EventParams, LoggedData)
						// err = CreateXML(ParsedEventTime)
						if err != nil {
							fmt.Println("Error when call logger():", err)
							return nil
						}
					} else {
						fmt.Printf("Parsed data must be greater than already stored at DB: %s\n", LastSavedEventTime)
					}
				} 
			}
			MoveFile(sourcePath, destPath)
		}
	}
	return nil
}
// --------------- work with database ---------------------
// GetData fetches data from the specified table in the SQLite database and processes it.
func getLastRecord(table string) (map[string]string, error) {
	// // Construct the SQL query to fetch the last row from the specified table
	query := fmt.Sprintf("SELECT * FROM '%s' ORDER BY datetime DESC LIMIT 1", table)
	fmt.Println(query)
	// Query the database
	_, data, err := getData(query)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	return data, nil
}

func getFirstRecord(table string, filter []string) (map[string]string, error) {
	var query string
	if filter != nil {
		query = fmt.Sprintf(`SELECT * FROM '%s' WHERE %s ORDER BY datetime ASC LIMIT 1`, table, strings.Join(filter, ","))
	} else {
		query = fmt.Sprintf(`SELECT * FROM '%s' ORDER BY datetime ASC LIMIT 1`, table)
	}
	fmt.Println(query)
	_, data, err := getData(query)
	if err != nil {
		return nil, err
	}
	fmt.Println(data)
	return data, nil
}

// InsertData inserts data into the specified table in the SQLite database.
func LogData(table string, EventParams []string, values []string) error {
	var err error
	columns := EventParams
	columns = append(columns, "createdAtDate")
	values = append(values, time.Now().Format(EventTimeFormat))
	if table == "DOCUMENTS" {
		columns = append(columns, "state")
		values = append(values, "STDBY")
	}
	// Construct the SQL query for insertion
	valuePlaceholders := strings.Repeat("?, ", len(columns)-1) + "?" // Repeat the "?" for placeholders
	EventParamsString := strings.Join(columns, ",")
	query := fmt.Sprintf("INSERT OR IGNORE INTO '%s' (%s) VALUES (%s)", table, EventParamsString, valuePlaceholders)
	fmt.Println(query)
    err = putData(db, query, values)
    if err != nil {
        log.Fatal(err)
    }
	return nil
}

// -------------- find & parse csv file -------------------
// ReadCSVFile reads the data from the provided CSV file and returns a list of Data.
func ReadCSVFile(fileName string) ([][]string, error) {
	// Open the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV records: %w", err)
	}

	var rows [][]string
	// var dataList []map[string]string

	// Skip the header row
	if len(records) > 0 {
		// columns := records[0]
		rows = records[1:]
	}
	return rows, nil
}

func findString(arr []string, target string) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}

// List files in filder with specified extention
func listFilesInDirectory(fullPath, extension string) ([]string, error) {
	
	// Check if the directory exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist: %s\n", fullPath)
		return nil, err
		}
	// List files
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}
	// Check if any file exist
	if len(files) == 0 {
		fmt.Printf("No files found in directory: %s\n", fullPath)
		return nil, nil
	}

	var filenames []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == extension {
			filenames = append(filenames, file.Name())
		}
	}

	return filenames, nil
}

// Compare and return True if date1 after date2 else return False
func date1_after_date2(dateStr1, dateStr2 string) (bool, error) {
	// Layout represents the format of the date string
	// layout := "02/01/2006"

	// Parse the date strings into time.Time objects
	date1, err := time.Parse(EventTimeFormat, dateStr1)
	if err != nil {
		return false, fmt.Errorf("error parsing date1: %w", err)
	}

	date2, err := time.Parse(EventTimeFormat, dateStr2)
	if err != nil {
		return false, fmt.Errorf("error parsing date2: %w", err)
	}

	// Compare the dates
	return date1.After(date2), nil
}

// GetData fetches data from the specified table in the SQLite database and processes it.
// It returns true if the operation is successful.
func getData(query string) ([]string, map[string]string, error) {

    rows, err := db.Query(query)
	// log.Printf(query)
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
	// If rows.Next() returns true, there's data to scan
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

// InsertData inserts data into the specified table in the SQLite database.
func putData(db *sql.DB, query string, values []string) error {
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

func MoveFile(sourcePath, destPath string) error {
    inputFile, err := os.Open(sourcePath)
    if err != nil {
        return fmt.Errorf("Couldn't open source file: %s", err)
    }
    outputFile, err := os.Create(destPath)
    if err != nil {
        inputFile.Close()
        return fmt.Errorf("Couldn't open dest file: %s", err)
    }
    defer outputFile.Close()
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {
        return fmt.Errorf("Writing to output file failed: %s", err)
    }
    // The copy was successful, so now delete the original file
    err = os.Remove(sourcePath)
    if err != nil {
        return fmt.Errorf("Failed removing original file: %s", err)
    }
    return nil
}

func ConverToEventDate(layout string, eventTime time.Time) (string) {
	return time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 13, 0, 0, 0, eventTime.Location()).Format(EventTimeFormat)
}

func SaveToFile(destPath, content string) {
	// filename = strings.Split(filename, ".")[0]
	// Marshal the struct back to formatted XML
	formattedXML, err := xml.MarshalIndent(content, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Save the formatted XML to a file
	err = ioutil.WriteFile(destPath, formattedXML, 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
}


// func zipString(lists ...[]string) func() []string {
//     zip := make([]string, len(lists))
//     i := 0
//     return func() []string {
//         for j := range lists {
//             if i >= len(lists[j]) {
//                 return nil
//             }
//             zip[j] = lists[j][i]
//         }
//         i++
//         return zip
//     }
// }

// // Helper function to parse a float value from a string and handle any errors.
// func parseFloat(str string) float32 {
// 	val, err := strconv.ParseFloat(strings.TrimSpace(str), 32)
// 	if err != nil {
// 		log.Printf("Error parsing float value %q: %s", str, err)
// 		return 0.0
// 	}
// 	return float32(val)
// }

// func parseInt(str string) int {
// 	str = "42"
// 	// Convert the string to an integer using Atoi
// 	num, err := strconv.Atoi(str)
// 	if err != nil {
// 		fmt.Println("Error converting string to integer:", err)
// 		return 9
// 	}
// 	return num
// }

// Common function
