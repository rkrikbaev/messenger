//------------------------Log CSV files
package main

import (
	"fmt"
	"app/logger/sqlHandler"
)

const (
	csv_data = "/csv_data/"
)

func LogDataFromFile(events []map[string]interface{}) error {
	path := path.Join(cwd, csv_data)
	today := time.Now()
	i := 0
	for i < 7 {
		newdt := today.AddDate(0,0,i-7).Format(CSVTimeFormat)
		fmt.Println(newdt)
		for _, event := range events {
			eventid := event["id"].(string)
			filename := fmt.Sprintf("%s_%s.csv", eventid, newdt)
			pathFile := fmt.Sprintf("%s/%s", path, filename)
			data, err := ReadCSVFile(pathFile)
			if err != nil {
				fmt.Printf("Can not read file: %s", pathFile)
				continue
			}
			var columns []string
			if findString(Type1Events, eventid) {
				columns = Type1Params
			} else if findString(Type2Events, eventid)  { 
				columns = Type2Params
			}
			dtString :=data[0][0]
			values := data[0]
			dt, _ := time.Parse(CSVTimeFormat, dtString)
			values[0] = ConverToEventDate(EventTimeFormat, dt)
			err = LogData(eventid, columns, values, "")
			if err != nil {
				continue
			}
			pathTo := fmt.Sprintf("%s/saved/%s",path,filename)
			MoveFile(pathFile, pathTo)
		}
		i++	
	}
	return nil
}

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

	// Skip the header row
	if len(records) > 0 {
		rows = records[1:]
	}
	return rows, nil
}