package main

import (
	"bytes"
	"database/sql"
	"io"
	"log"
	"path/filepath"

	// "database/sql"
	"encoding/csv"
	// "encoding/json"

	"fmt"
	// "io"
	"io/ioutil"
	// "log"
	"net/http"
	// "net/smtp"
	"os"
	"path"

	// "reflect"
	"strings"
	"time"

	// "github.com/gokalkan/gokalkan"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type EventType interface {
	processFile()
}

type Document struct {
	id int
	Document string
	EventDate string
	CreatedAt string
	MessageSent string
	State string
}

type EventTypeFlow struct {
	id string
	datetime string
	deviceTypeId string
	operationTypeId string
	deviceNameId string
	productTypeId string
	temperature string
	density string
	volume string
	pipelineId string
	massflowbegin string
	massflowend string
	mass string
}

type EventTypeLevel struct {
	id string
	datetime string
	deviceTypeId string
	operationTypeId string
	deviceNameId string
	productTypeId string
	mass string
	tankLevel string
	volume string
	temperature string
	density string
}

const (
	StateTable = "RESPONSE"
	XMLDocumentTable = "DOCUMENTS"
	CSVTimeFormat = "02.01.2006"
	EventTimeFormat = "2006-01-02T15:04:05.000-07:00"
	ip = "127.0.0.1"//"195.12.113.29"
	port = "8080"
	ReportTable = "REPORT"
	dbname = "app.db"
)

var (
	eventsMap map[string]interface{}
	events = map[string]string{}
	EventTime string
	EventRecordDate time.Time
	// filesSourcePath string
	FilesDirPath string
	extension string
	mailReceiver string
	db     *sql.DB
	allTabels = []string{"407001","407002","407003", "407004","407005","407006","407007","407008","407009"}
	serviceId string
	senderId string
	senderPassword string

	cwd string
	dbPath string
	certPath string
	certPassword string
	randomUUID string
)

func typeFlow(id string) *EventTypeFlow {
	return &EventTypeFlow{
		id: id,
		deviceTypeId:       "2",	// Default value
		operationTypeId: 	"3", 	// Default value
		deviceNameId:    	"1",	// Default value
		productTypeId:    	"1",	// Default value
		pipelineId: 		"3",
		}
	}

func typeLevel(id string) *EventTypeLevel {
	return &EventTypeLevel{
		id: id,
		deviceTypeId:       "1",	// Default value
		operationTypeId: 	"7", 	// Default value
		deviceNameId:    	"2",	// Default value
		productTypeId:    	"1",	// Default value
		}
	}


func main() {
	var err error
	fmt.Printf("Start")
	cwd, _ = os.Getwd()
	dbPath = path.Join(cwd, "/db/", dbname)
	
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	env := path.Join(cwd, "app.env")
	err = godotenv.Load(env)
	if err != nil {
		fmt.Printf("error parse env file: ", env)
	}
	certPath = os.Getenv("certPath")
	certPassword = os.Getenv("certPassword")
	serviceId = os.Getenv("serviceId")
	senderId = os.Getenv("senderId")
	senderPassword = os.Getenv("senderPassword")
	mailReceiver = os.Getenv("mailReceiver")

	eventsMap["407001"] = typeFlow("407001")
	eventsMap["407002"] = typeFlow("407002")
	eventsMap["407003"] = typeFlow("407003")
	eventsMap["407004"] = typeLevel("407004")
	eventsMap["407005"] = typeLevel("407005")
	eventsMap["407006"] = typeLevel("407006")
	eventsMap["407007"] = typeLevel("407007")
	eventsMap["407008"] = typeLevel("407008")
	eventsMap["407009"] = typeLevel("407009")

	run()
}

func run() {
	TodayDatetime := time.Now()
	for {
		var err error
		if TodayDatetime.After(EventRecordDate) {
		filesSourcePath := path.Join(cwd, "/csv_data/")
		err = LogDataFromFile(filesSourcePath)
		if err != nil {
			fmt.Printf("Error when call LogDataFromFile():", err)
			return
		}
		err = processEvent()
		if err != nil {
			fmt.Printf("Error when call CreateXML():", err)
			return
		}
		} else {
			break
		}
		time. Sleep(10 * time.Second)
	}
}
//------------------------Log CSV files
func LogDataFromFile(path string) error {	
	var err error	
	files, err := listFilesInDirectory(path, ".csv")
	if err != nil {
		return err
	}
	for _, filename := range files {
		fileSavePath := fmt.Sprintf("%s/saved/%s",path,filename)
		eventType := strings.Split(filename, "_")[0]
		sourceFile := fmt.Sprintf("%s/%s", path, filename)
		sourceData, err := ReadCSVFile(sourceFile)
		if err != nil {
			log.Fatalf("Error reading CSV file: %s", err)
		}
		err = processFile(eventType, sourceData)
		if err != nil {
			log.Fatalf("Error process CSV file: %s", err)
		}
		MoveFile(sourceFile, fileSavePath)
		}
	return nil
}

func processFile(eventType string, sourceData [][]string) error {
	object := eventsMap[eventType] 
	dt, _ := time.Parse(CSVTimeFormat, sourceData[1][0])
	sourceDatetime := ConverToEventDate(EventTimeFormat, dt)
	savedData, _ := getLastRecord(eventType)
	savedDatetime := savedData["datetime"]
	var e map[string]string
	if _, ok := object.(*EventTypeFlow); ok {
		e["datetime"] = sourceDatetime
		e["mass"] = sourceData[1][1]
		e["massflowbegin"] = sourceData[1][2]
		e["massflowend"] = sourceData[1][3]
		e["temperature"] = sourceData[1][4]
		e["density"] = sourceData[1][5]
		e["volume"] = sourceData[1][6]
	} else if _, ok := object.(*EventTypeLevel); ok {
		e["datetime"] = sourceDatetime
		e["mass"] = sourceData[1][1]
		e["tankLevel"] = sourceData[1][2]
		e["volume"] = sourceData[1][3]
		e["temperature"] = sourceData[1][4]
		e["density"] = sourceData[1][5]
	} else {
		// Не удалось определить тип события
		fmt.Printf("Не определен тип события")
	}
	if len(savedData) == 0 {
		_ = LogData(eventType, e)
	} else {
		
		if isAfter, err := date1_after_date2(sourceDatetime, savedDatetime); err != nil {
					log.Fatalf("Error comparing dates: %s", err)
				} else if isAfter {
					err = LogData(TableName, EventParams, LoggedData)
					if err != nil {
						fmt.Printf("Error when call logger():", err)
						return nil
					}
				} else {
					fmt.Printf("Parsed data must be greater than already stored in DB: %s\n", LastSavedEventTime)
				}
			}		
		
	// event.datetime = LastSaveData["datetime"]
	// ParsedEventTime := ConverToEventDate(EventTimeFormat, LoggedDataTime)
	// LoggedData[0] = ParsedEventTime
	// if len(LastSaveData) == 0 {
	// 	_ = LogData(TableName, EventParams, LoggedData)
	// } else {
	// 	if isAfter, err := date1_after_date2(ParsedEventTime, LastSavedEventTime); err != nil {
	// 		log.Fatalf("Error comparing dates: %s", err)
	// 	} else if isAfter {
	// 		err = LogData(TableName, EventParams, LoggedData)
	// 		// err = CreateXML(ParsedEventTime)
	// 		if err != nil {
	// 			event["Error when call logger():", err)
	// 			return nil
	// 		}
	// 	} else {
	// 		fmt.Printf("Parsed data must be greater than already stored at DB: %s\n", LastSavedEventTime)
	// 	}
	// } 

	return nil
}
//------------------------Prepare XML document
func processEvent() error {
	var err error
	var eventDate string
	randomstring := generateRandomString()
	lastCreatedXML, err := getLastRecord("DOCUMENTS")
	if err != nil {
		return err
	}	
	lastEventDate := lastCreatedXML["datetime"]
	if lastEventDate == "" {
		var dateInTable string
		dateInTable = ""
		for idx, table := range allTabels {
			data, err := getFirstRecord(table, nil)
			if err != nil {
				return err
			}
			if idx == 0 {
				dateInTable = data["datetime"]
			} else {
				if data["datetime"] == dateInTable {
					dateInTable = data["datetime"]
				} else {
					dateInTable = ""
					break
				}
			}
		}
		if dateInTable != "" {
			eventDate = dateInTable
		}
	} else {
		eventPrevDate, _ := time.Parse(EventTimeFormat, lastEventDate)
		eventNewDate := eventPrevDate.AddDate(0, 0, 1)
		eventDate = eventNewDate.Format(EventTimeFormat)
	}
	EventData, err := getEventData(eventDate)

	if err != nil {
		return err
	}
	if EventData == "" {
		return nil
	}
	dateTimeNowString := time.Now().Format(EventTimeFormat)
	var xmlBuffer bytes.Buffer
	xmlBuffer.WriteString(fmt.Sprintf(
		`<ns2:SendMessage xmlns:ns2="http://bip.bee.kz/SyncChannel/v10/Types"> <request> <requestInfo> <messageId>%s</messageId> <serviceId>%s</serviceId> <messageDate>%s</messageDate> <sender> <senderId>%s</senderId> <password>%s</password> </sender> </requestInfo> <requestData> <data xmlns:cs="http://message.persistence.interactive.nat" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="cs:Request">%s</data> </requestData> </request> </ns2:SendMessage>`,
		randomstring, serviceId, dateTimeNowString, senderId, senderPassword, EventData))
	xmlString := xmlBuffer.String()
	fmt.Printf(xmlString)
	err = LogData("DOCUMENTS", []string{"document", "datetime"}, []string{xmlString, eventDate})
	if err != nil {
		return err
	}
	content, err := SendMessage(xmlString)
	if err != nil {
		fmt.Printf("Error when call SendMessage():", err)
		return err
	}
	destPath := fmt.Sprintf("/app/xml_data/message_%s.xml",eventDate)
	SaveToFile(destPath, content)
	return nil
}

func getEventData(eventDate string) (string, error) {
	var eventDataArray []string
	for _, item := range objects {
		if itemMap, ok := item.([]interface{}); ok {
			for _, value:= range itemMap {
				if eventData, ok := value.(map[string]interface{}); ok {
					getEventData, err := generateXMLString(eventData, eventDate)
					if getEventData == "" {
						return "", nil
					} else {
					eventDataArray = append(eventDataArray, getEventData)
					}
					if err != nil {
						errMsg2 := fmt.Sprintf("Error generate XML with data): %s", err)
						fmt.Printf(errMsg2)
						return "", err
					}
				}
			}
		}
	}
	eventData := strings.Join(eventDataArray, " ")
	return eventData, nil
}

func generateXMLString(event map[string]interface{}, EventDate string) (string, error) {
	// get process data
	excludeEventItems := map[string]bool{
		"deviceTypeID":     true,
		"operationTypeID":  true,
		"productTypeID":    true,
		"pipelineID":       true,
		"deviceNameID":     true,
		"id":				true,
		// "sentMessageDate":	false,
		// "state":			false,
		// "createdAtDate":	false,
	}

	exludeDataItems := map[string]bool{
		"ID":				true,
		"createdAtDate":	true,
	}

	table, ok := event["id"].(string)
	if !ok {
		return "", fmt.Errorf("event ID is not a string")
	}
	filter := []string{fmt.Sprintf("datetime = '%s'", EventDate)}
	data, err := getFirstRecord(table, filter)
	if err != nil {
		return ``, err
	}
	if data == nil {
		return "", nil
	}
	EventRecordDate, _ = time.Parse(EventTimeFormat,data["datetime"])
	// create xml string
	var xmlBuffer bytes.Buffer
	xmlBuffer.WriteString("<events>")
	// add static configuration 
	for key, value := range event {
		if !excludeEventItems[key] {
			continue
		}
		xmlBuffer.WriteString(fmt.Sprintf("<%s>%s</%s>", key, fmt.Sprintf("%v", value), key))
	}
	// add process values
	for key, value := range data {
		if exludeDataItems[key] {
			continue
		} else {
			if value == "<nil>" {
				value = "0"
			}
			xmlBuffer.WriteString(fmt.Sprintf("<%s>%s</%s>", key, value, key))
		}
	}
	xmlBuffer.WriteString("</events>")
	return xmlBuffer.String(), nil
}

//------------------------Send SOAP message
func SendMessage(xmlstring string) (string, error) {
	state := "FAIL"

	// для теста
	// opts := gokalkan.OptsTest

	// для прода
	opts := gokalkan.OptsProd
	cli, err := gokalkan.NewClient(opts...)
	if err != nil {
		fmt.Print(fmt.Sprintf("ERROR, new kalkan client create error: %s", err))
		return "", err
	}
	defer cli.Close()

	filter := []string{"state=='STDBY'"}
	data, _ := getFirstRecord("DOCUMENTS", filter)
	if data == nil {
		return "", nil
	}
	xmlString := data["document"]
	if xmlString == "" {
		return "", nil
	}
	eventDate := data["datetime"]
	if eventDate == "" {
		return "", nil
	}

	//sign message
	randomUUID = generateRandomString()
	err = cli.LoadKeyStore(certPath, certPassword)
	if err != nil {
		return "", err
	}
	message, err := cli.SignWSSE(xmlString, fmt.Sprintf("id-%s", randomUUID))
	fmt.Printf(randomUUID)
	fmt.Printf(message)
	err = sendRequest(ip, port, message)
	if err != nil {
		fmt.Printf("Error creating request:", err)
		return "", err 
	}
	table := "DOCUMENTS"
	err = updateData(table, eventDate, state)
	if err != nil {
		fmt.Printf("Error update data after success send message", err)
		return "", err 
	}
return message, nil
}

func sendRequest(ip, port, message string) (error) {

	// Set the URL and HTTP method
	url := fmt.Sprintf("http://%s:%s/bip-sync-wss-gost/", ip, port)
	method := "POST"
	payload := strings.NewReader(message)

	// Create a request with the SOAP message
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Printf("Error creating request:", err)
		return err 
	}
	// Set headers for the SOAP request
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	// Make the HTTP request
	res, err := client.Do(req)
	if err != nil {
		return err 
	}
	defer res.Body.Close()
	
	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err 
	}
	responseString := string(responseBody)
	fmt.Printf("Response Body:", string(responseString))
	return nil
}

// Update status from STDBY to DONE
func updateData(table, eventDate, state string) (error) {
	var err error
	columns := []string{"state", "sentMessageDate"}
	values := []string{state, time.Now().Format(EventTimeFormat)}

	var columnValuePairs []string
	
	// for i := 0; i > len(columns); i++ {
	for i := range columns {
		columnValuePairs = append(columnValuePairs, fmt.Sprintf(`%s = '%s'`, columns[i], values[i]))
	}
	condition := fmt.Sprintf("datetime = '%s'", eventDate)
	// UPDATE 'DOCUMENTS' SET state='SUCCESS' WHERE datetime = '2023-08-27T13:00:00.000+00:00'
	query := fmt.Sprintf(`UPDATE '%s' SET %s WHERE (%s)`, table, strings.Join(columnValuePairs, ","), condition)
	err = putData(db, query, values)
	if err != nil {
		return err
	}
	return nil
}

//
func generateRandomString() string {
	uuid := uuid.New()
	return uuid.String()
}

// GetData fetches data from the specified table in the SQLite database and processes it.
func getLastRecord(table string) (map[string]string, error) {
	// // Construct the SQL query to fetch the last row from the specified table
	query := fmt.Sprintf("SELECT * FROM '%s' ORDER BY datetime DESC LIMIT 1", table)
	fmt.Printf(query)
	// Query the database
	_, data, err := getData(query)
	if err != nil {
		return nil, err
	}
	// event[data)
	return data, nil
}

func getFirstRecord(table string, filter []string) (map[string]string, error) {
	var query string
	if filter != nil {
		query = fmt.Sprintf(`SELECT * FROM '%s' WHERE %s ORDER BY datetime ASC LIMIT 1`, table, strings.Join(filter, ","))
	} else {
		query = fmt.Sprintf(`SELECT * FROM '%s' ORDER BY datetime ASC LIMIT 1`, table)
	}
	fmt.Printf(query)
	_, data, err := getData(query)
	if err != nil {
		return nil, err
	}
	fmt.Printf(data)
	return data, nil
}

// InsertData inserts data into the specified table in the SQLite database.
func LogData(table string, data map[string]string) error {
	var err error
	columns := EventParams
	columns = append(columns, "createdAtDate")
	values = append(values, time.Now().Format(EventTimeFormat))
	if table == "DOCUMENTS" {
		columns = append(columns, "state")
		values = append(values, "STDBY")
	}
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
func putData(db *sql.DB, data map[string]interface{}) error {
	for col := range data {
		columns := append(columns, col)
		values := append(values, data[col])
	}

	// Construct the SQL query for insertion
	valuePlaceholders := strings.Repeat("?, ", len(columns)-1) + "?" // Repeat the "?" for placeholders
	EventParamsString := strings.Join(columns, ",")
	query := fmt.Sprintf("INSERT OR IGNORE INTO '%s' (%s) VALUES (%s)", table, EventParamsString, valuePlaceholders)

	// d := make([]interface{}, len(values))
	// for i, v := range data {
	// 	d[i] = v
	// }
	tx, err := db.Begin() // Start a transaction
    if err != nil {
        return fmt.Errorf("error beginning transaction: %w", err)
    }
    
    defer tx.Rollback() // Rollback the transaction in case of errors
    
    _, err = tx.Exec(query, values...)
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

func SaveToFile(destPath, content string) {
	// Open the file for writing
	file, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Write data to the file
	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}	
	if err != nil {
		fmt.Printf("Error saving file:", err)
		return
	}
}

func ConverToEventDate(layout string, eventTime time.Time) (string) {
	return time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 13, 0, 0, 0, eventTime.Location()).Format(EventTimeFormat)
}
// func setRootDir() {
// 	// Получаем путь к папке, в которой находится исполняемый файл
// 	executablePath, err := os.Executable()
// 	if err != nil {
// 		event["Error:", err)
// 		return
// 	}

// 	// Извлекаем путь к папке из пути к исполняемому файлу
// 	rootPath := filepath.Dir(executablePath)

// 	// Меняем текущую рабочую директорию на корневую папку
// 	err = os.Chdir(rootPath)
// 	if err != nil {
// 		event["Error:", err)
// 		return
// 	}

// 	// Теперь текущая директория установлена в корневую папку
// 	event["Root path:", rootPath)
// }













// func sendMail(to string, subject string, body string) {
// 	// Set up authentication information.
// 	mailSender   := os.Getenv("mailSender")
// 	appPassword  := os.Getenv("appPassword")

// 	smtpHost := os.Getenv("smtpHost")
// 	smtpPort := os.Getenv("smtpPort")

// 	// Compose the email message
// 	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", mailSender, to, subject, body)

// 	// SMTP authentication setup
// 	auth := smtp.PlainAuth("", mailSender, appPassword, smtpHost)

// 	// Send the email using mail.ru SMTP
// 	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, mailSender, []string{to}, []byte(message))
// 	if err != nil {
// 		logging(fmt.Sprintf("Error sending email: %s",err))
// 		return
// 	}

// 	logging("Email sent successfully.")
// }

// func logging(text string) {
// 	// Get the current date and time
// 	now := time.Now()

// 	// Format the date as "YYYY-M-D"
// 	date := now.Format("2006-1-2")

// 	// Format the time as "HH:MM:SS"
// 	time := now.Format("15:04:05")

// 	// Create the log text
// 	logText := fmt.Sprintf("%s %s: %s \n", date, time, text)
// 	logFileName :="log.log"
// 	// check for file exist if not create one
// 	checkFileExistIfNotCreate(logFileName)

// 	// Append the log text to the log file
// 	file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
// 	if err != nil {
// 		event["Error opening file:", err)
// 	}
// 	defer file.Close()

// 	// Append text to the file
// 	_, err = io.WriteString(file, logText)
// 	if err != nil {
// 		event["Error writing to file:", err)
// 	}
// }

// func checkFileExistIfNotCreate(filePath string) {
// 	// check for file exist if not create one
// 	// Check if the file exists
// 	_, err := os.Stat(filePath)
// 	if err != nil {
// 		// File does not exist, create it
// 		file, err := os.Create(filePath)
// 		if err != nil {
// 			event["Error creating file:", err)
// 			return
// 		}
// 		defer file.Close()

// 		event["File created:", filePath)
// 	}
// }