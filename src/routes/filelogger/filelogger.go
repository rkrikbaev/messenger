package filelogger

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

const(
	EventTimeFormat = "2006-01-02T15:04:05.000-07:00"
	CSVTimeFormat = "02.01.2006"
)

// ReadCSVFile reads the data from the provided CSV file and returns a list of Data.
// readCSVFile reads a CSV file and returns the headers and records as slices.
// It takes the path of the CSV file as input and returns the headers, records, and any error encountered.
func readCSVFile(pathFile string) ([]string, [][]string, error) {
	// Open the CSV file
	file, err := os.Open(pathFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)
	reader.Comma = ','

	headers, err := reader.Read() // Читаем только заголовки
	if err != nil {
		return nil, nil, fmt.Errorf("unable to read the CSV file headers: %v", err)
	}

	// Read the CSV records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("error reading CSV records: %w", err)
	}
	return headers, records, err
}

func ParseAndPrepareData(pathFile string) ([]string, []string, error) {

	// Чтение и обработка cvs файла
	headers, records, err := readCSVFile(pathFile)
	if err != nil {
		fmt.Printf("Cannot read file: %s\n", pathFile)
	}

	dtValue := records[0][0]
	
	dt, err := time.Parse(CSVTimeFormat, dtValue)
	if err != nil {
		fmt.Printf("Error parsing datetime in file: %s\n", err)
	}

	values := records[0]
	values[0] = converToEventDate(EventTimeFormat, dt)

	return headers, values, err
}	

func converToEventDate(layout string, eventTime time.Time) (string) {
	return time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 13, 0, 0, 0, eventTime.Location()).Format(EventTimeFormat)
}