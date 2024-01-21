package main

import (
	"bytes"
	"crypto/tls"
	// "context"
	"database/sql"
	"io"
	"log"
	"path/filepath"
	"unicode"

	// "database/sql"
	"encoding/csv"
	"encoding/json"

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

	"github.com/gokalkan/gokalkan"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

//СИКН Амангельды    407001
//СИКН Айракты   407002
//СИКН Жаркум   407003
//
//Резервуар Амангельды - 1   407004
//Резервуар Амангельды - 2   407005
//
//Резервуар Айракты - 1   407006
//Резервуар Айракты - 2   407007
//
//Резервуар Жаркум - 1   407008
//Резервуар Жаркум - 2   407009

type Event struct {
	Id int
	Document string
	EventDate string
	CreatedAt string
	MessageSent string
	State string
}

const (
	StateTable = "RESPONSE"
	XMLDocumentTable = "DOCUMENTS"
	CSVTimeFormat = "02.01.2006"
	EventTimeFormat = "2006-01-02T15:04:05.000-07:00"
	ip = "195.12.113.29"
	port = "80"
	ReportTable = "REPORT"
	dbname = "app.db"
	jsonData = `[
					{
						"id": "407001",
						"deviceTypeId": 2,
						"operationTypeId": 3,
						"deviceNameId": 1,
						"productTypeId": 1,
						"pipelineId": 3,
						"temperature": 0,
						"density": 0,
						"volume": 0,
						"massflowbegin": 0,
						"massflowend": 0,
						"mass": 0
					},
					{
						"id": "407004",
						"deviceTypeId": 1,
						"operationTypeId": 7,
						"deviceNameId": 2,
						"productTypeId": 1,
						"temperature": 0,
						"density": 0,
						"volume": 0,
						"tenkLevel": 0,
						"mass": 0
					},
					{
						"id": "407005",
						"deviceTypeId": 1,
						"operationTypeId": 7,
						"deviceNameId": 2,
						"productTypeId": 1,
						"temperature": 0,
						"density": 0,
						"volume": 0,
						"tenkLevel": 0,
						"mass": 0
					},
					{
						"id": "407002",
						"deviceTypeId": 2,
						"operationTypeId": 3,
						"deviceNameId": 1,
						"productTypeId": 1,
						"pipelineId": 3,
						"temperature": 0,
						"density": 0,
						"volume": 0,
						"massflowbegin": 0,
						"massflowend": 0,
						"mass": 0
					},
					{
						"id": "407006",
						"deviceTypeId": 1,
						"operationTypeId": 7,
						"deviceNameId": 2,
						"productTypeId": 1,
						"temperature": 0,
						"density": 0,
						"volume": 0,
						"tenkLevel": 0,
						"mass": 0
					},
					{
						"id": "407007",
						"deviceTypeId": 1,
						"operationTypeId": 7,
						"deviceNameId": 2,
						"productTypeId": 1,
						"temperature": 0,
						"density": 0,
						"volume": 0,
						"tenkLevel": 0,
						"mass": 0
					},
					{
						"id": "407003",
						"deviceTypeId": 2,
						"operationTypeId": 3,
						"deviceNameId": 1,
						"productTypeId": 1,
						"pipelineId": 3,
						"temperature": 0,
						"density": 0,
						"volume": 0,
						"massflowbegin": 0,
						"massflowend": 0,
						"mass": 0
					},
					{
						"id": "407008",
						"deviceTypeId": 1,
						"operationTypeId": 7,
						"deviceNameId": 2,
						"productTypeId": 1,
						"temperature": 0,
						"density": 0,
						"volume": 0,
						"tenkLevel": 0,
						"mass": 0
					},
					{
						"id": "407009",
						"deviceTypeId": 1,
						"operationTypeId": 7,
						"deviceNameId": 2,
						"productTypeId": 1,
						"temperature": 0,
						"density": 0,
						"volume": 0,
						"tenkLevel": 0,
						"mass": 0
					}
				]`
		)

var (
	EventTime string
	EventRecordDate time.Time
	events []map[string]interface{}
	FilesDirPath string
	extension string
	mailReceiver string
	db     *sql.DB
	Events Event
	// allevents = []string{"407001","407002","407003", "407004","407005","407006","407007","407008","407009"}
	serviceId string
	senderId string
	senderPassword string
	Type1Events = []string{"407001","407002","407003"}
	Type2Events = []string{"407004","407005","407006","407007","407008","407009"}
	Type1Params = []string{"datetime","density","massflowbegin","massflowend","mass"}
	Type2Params = []string{"datetime","density","volume","temperature","tankLevel","mass"}
	cwd = ""
	dbPath string
	certPath string
	certPassword string
	randomUUID string
	location *time.Location
)

func main() {
	var err error
	fmt.Println("Start")
	cwd, _ = os.Getwd()
	// cwd = "/Users/rustamkrikbayev/projects/parser/app"
	dbPath = path.Join(cwd, "/db/", dbname)
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	env := path.Join(cwd, "app.env")
	err = godotenv.Load(env)
	if err != nil {
		fmt.Println("error parse env file: ", env)
	}
	certPath = os.Getenv("certPath")
	certPassword = os.Getenv("certPassword")
	serviceId = os.Getenv("serviceId")
	senderId = os.Getenv("senderId")
	senderPassword = os.Getenv("senderPassword")
	mailReceiver = os.Getenv("mailReceiver")
	err = json.Unmarshal([]byte(jsonData), &events)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	location, err = time.LoadLocation("Asia/Almaty")
	if err != nil {
		fmt.Println("Ошибка загрузки часового пояса:", err)
		return
	}
	run()
}

func run() {
	var err error
	for {
		
		fmt.Println("LogDataFromFile")
		err = LogDataFromFile(events)
		if err != nil {
			fmt.Println("Error when log Data:", err)
		}
		fmt.Print("processEvent")
		err = processEvent()
		if err != nil {
			fmt.Println("Error when processing Event:", err)
		}		
		time.Sleep(600 * time.Second)
	}
}

//------------------------Log CSV files
func LogDataFromFile(events []map[string]interface{}) error {
	path := path.Join(cwd, "/csv_data/")
	today := time.Now()
	i := 0
	for i < 14 {
		newdt := today.AddDate(0,0,i-14).Format(CSVTimeFormat)
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
//------------------------Prepare XML docement
func processEvent() error {
	var err error
	// var eventDate string
	// Try to find latest recordset with state STDBY
	// And sebd it after 18:00
	eventRecordSet, _ := getLastRecord("DOCUMENTS", []string{"state='STDBY'"})
	if len(eventRecordSet) > 0 {
		date1 := time.Now().In(location)
		date2 := time.Date(date1.Year(), date1.Month(), date1.Day(), 18,0,0,0, date1.Location())
		if date1.After(date2) {
			DocumentXML := eventRecordSet["document"]
			eventDate := eventRecordSet["datetime"]
			err = SendMessage(DocumentXML, eventDate)
			if err != nil {
				fmt.Println("Error when call SendMessage():", err)
			}
		} else {
			fmt.Println("Waiting...")
		}
	} else {
		// Try to create XML document after 13:00 wyen process data come
		fmt.Println("Collecting data...")
		eventRecordSet, _ = getLastRecord("DOCUMENTS", nil)
		if len(eventRecordSet) > 0 {
			existEventDate := eventRecordSet["datetime"]
			err = createDocumentXML(existEventDate)
			if err != nil {
				fmt.Println("Error when call createDocumentXML():", err)
			}
		}
	}
	return nil
}

func createDocumentXML(existEventDate string) (error) {
	randomstring := generateRandomString()
	var err error
	exDate, err := time.Parse(EventTimeFormat, existEventDate)
	if err != nil {
		return nil
	}
	todayDate := time.Now()
	for todayDate.After(exDate) {
		newDate := exDate.AddDate(0, 0, 1)
		exDate = newDate
		eventDate := newDate.Format(EventTimeFormat)
		fmt.Print("Try to collect data for new XML document\n")
		fmt.Print(eventDate)
		
		eventData, err := getEventData(eventDate)
		if stringIsEmpty(eventData) {
			continue
		}
		if err != nil {
			continue
		}
		var xmlBuffer bytes.Buffer
		xmlBuffer.WriteString(fmt.Sprintf(
			`<ns2:SendMessage xmlns:ns2="http://bip.bee.kz/SyncChannel/v10/Types"> <request> <requestInfo> <messageId>%s</messageId> <serviceId>%s</serviceId> <messageDate>%s</messageDate> <sender> <senderId>%s</senderId> <password>%s</password> </sender> </requestInfo> <requestData> <data xmlns:cs="http://message.persistence.interactive.nat" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="cs:Request">%s</data> </requestData> </request> </ns2:SendMessage>`,
			randomstring, serviceId, todayDate.Format(EventTimeFormat), senderId, senderPassword, eventData))
		xmlString := xmlBuffer.String()
		err = LogData("DOCUMENTS", []string{"document", "datetime"}, []string{xmlString, eventDate}, "STDBY")
		if err != nil {
			fmt.Print(err)
		}
		
		time. Sleep(1 * time.Second)
	}
	return nil
}

func stringIsEmpty(String string) bool {
	//
	isEmpty := true
	for _, char := range String {
		if !unicode.IsSpace(char) {
			isEmpty = false
			break
		}
	}
	if isEmpty {
		fmt.Println("Empty string")
	}
	return isEmpty
}

func getEventData(eventDate string) (string, error) {
	var eventDataArray []string
	for _, event:= range events {
		getEventData, err := generateXMLString(event, eventDate)
		if err != nil {
			errMsg2 := fmt.Sprintf("Error generate XML with data): %s", err)
			fmt.Println(errMsg2)
			return "", err
		}
		if stringIsEmpty(getEventData) {
			fmt.Printf("Empty string for: %s", eventDate)
			return "", nil
		} else {
			eventDataArray = append(eventDataArray, getEventData)
		}
	}
	eventData := strings.Join(eventDataArray, " ")
	return eventData, nil
}

func generateXMLString(event map[string]interface{}, EventDate string) (string, error) {
	// get process data
	excludeEventItems := map[string]bool{
		"deviceTypeId":     true,
		"operationTypeId":  true,
		"productTypeId":    true,
		"pipelineId":       true,
		"deviceNameId":     true,
	}

	exludeDataItems := map[string]bool{
		"ID":				true,
		"createdAtDate":	true,
		"datetime":			true,
	}

	eventId, ok := event["id"].(string)
	if !ok {
		return "", fmt.Errorf("event ID is not a string")
	}
	filter := []string{fmt.Sprintf("datetime = '%s'", EventDate)}
	data, err := getFirstRecord(eventId, filter)
	if err != nil {
		return ``, err
	}
	if data == nil {
		return "", nil
	}
	datetime := data["datetime"]
	EventRecordDate, _ = time.Parse(EventTimeFormat,datetime)
	// create xml string
	var xmlBuffer bytes.Buffer
	xmlBuffer.WriteString("<events>")
	xmlBuffer.WriteString(fmt.Sprintf("<id>%s</id>", eventId))
	xmlBuffer.WriteString(fmt.Sprintf("<datetime>%s</datetime>", datetime))
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
func SendMessage(xmlstring, eventDate string) (error) {
	var err error
	state := "FAIL"
	table := "DOCUMENTS"
	condition := fmt.Sprintf("datetime = '%s'", eventDate)
	// для теста
	// opts := gokalkan.OptsTest

	opts := gokalkan.OptsProd
	cli, err := gokalkan.NewClient(opts...)
	if err != nil {
		fmt.Print(fmt.Sprintf("ERROR, new kalkan client create error: %s", err))
		return err
	}
	defer cli.Close()

	//sign message
	randomUUID = generateRandomString()
	err = cli.LoadKeyStore(certPath, certPassword)
	if err != nil {
		return err
	}
	message, err := cli.SignWSSE(xmlstring, fmt.Sprintf("id-%s", randomUUID))
	// message := "uuuuuuuuuuuuuuu"
	fmt.Println(randomUUID)
	fmt.Println(message)

	destPath := fmt.Sprintf("%s/xml_data/message_%s.xml",cwd,eventDate)
	SaveToFile(destPath, message)

	// timeout := 10 * time.Second // Пример: таймаут 10 секунд
	// err = sendRequest(ip, port, message, timeout)
	err = sendRequest(message)
	if err != nil {
		fmt.Println("Error creating request:", err)
		err = updateState(table, condition, state)
		return err 
	}
	state = "SUCCESS"
	err = updateState(table, condition, state)
	if err != nil {
		fmt.Println("Error update data after success send message", err)
		return err 
	}
return nil
}

func sendRequest(data string) error {

	url := "https://195.12.113.29/bip-sync-wss-gost/"
	method := "POST"

	payload := strings.NewReader(data)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client {Transport: tr }
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
	fmt.Println(err)
	return err
	}
	req.Header.Add("Content-Type", "text/xml")

	res, err := client.Do(req)
	if err != nil {
	fmt.Println(err)
	return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
	fmt.Println(err)
	return err
	}
	fmt.Println(string(body))
	return nil
}

// func sendRequest(ip, port, message string, timeout time.Duration) ([]byte, error) {
//     url := fmt.Sprintf("https://%s/bip-sync-wss-gost/", ip)
//     method := "POST"
//     payload := strings.NewReader(message)

//     client := &http.Client{}
//     ctx, cancel := context.WithTimeout(context.Background(), timeout)
//     defer cancel() // Освобождаем ресурсы контекста после выполнения функции
    
//     req, err := http.NewRequestWithContext(ctx, method, url, payload)
//     if err != nil {
//         fmt.Println("Error creating request:", err)
//         return nil, err
//     }
//     req.Header.Set("Content-Type", "text/xml; charset=utf-8")

//     res, err := client.Do(req)
//     if err != nil {
//         fmt.Println("Error making HTTP request:", err)
//         return nil, err
//     }
//     defer res.Body.Close()

//     if res.StatusCode != http.StatusOK {
//         fmt.Printf("Received non-success status code: %d\n", res.StatusCode)
//         return nil, fmt.Errorf("non-success status code: %d", res.StatusCode)
//     }

//     responseBody, err := ioutil.ReadAll(res.Body)
//     if err != nil {
//         fmt.Println("Error reading response body:", err)
//         return nil, err
//     }

//     fmt.Println("Response Body:", string(responseBody))
//     return responseBody, nil
// }



// func sendRequest(ip, port, message string) (error) {

// 	// Set the URL and HTTP method
// 	url := fmt.Sprintf("http://%s:%s/bip-sync-wss-gost/", ip, port)
// 	method := "POST"
// 	payload := strings.NewReader(message)

// 	// Create a request with the SOAP message
// 	client := &http.Client {
// 	}
// 	req, err := http.NewRequest(method, url, payload)
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 		return err 
// 	}
// 	// Set headers for the SOAP request
// 	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
// 	// Make the HTTP request
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return err 
// 	}
// 	defer res.Body.Close()
	
// 	responseBody, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return err 
// 	}
// 	responseString := string(responseBody)
// 	fmt.Println("Response Body:", string(responseString))
// 	return nil
// }


// Update status from STDBY to DONE
func updateState(table, condition, state string) (error) {
	var err error
	columns := []string{"state", "sentMessageDate"}
	values := []string{state, time.Now().Format(EventTimeFormat)}
	var columnValuePairs []string
	for i := range columns {
		columnValuePairs = append(columnValuePairs, fmt.Sprintf(`%s = '%s'`, columns[i], values[i]))
	}
	// condition := fmt.Sprintf("datetime = '%s'", eventDate)
	// UPDATE 'DOCUMENTS' SET state='SUCCESS' WHERE datetime = '2023-08-27T13:00:00.000+00:00'
	query := fmt.Sprintf(`UPDATE '%s' SET %s WHERE (%s)`, table, strings.Join(columnValuePairs, ","), condition)
	err = putData(db, query, values)
	if err != nil {
		return err
	}
	return nil
}

// //----------------------------- // ------------------------------------------ //
func generateRandomString() string {
	uuid := uuid.New()
	return uuid.String()
}

// --------------- work with database ---------------------
// GetData fetches data from the specified table in the SQLite database and processes it.
func getLastRecord(table string, filter []string) (map[string]string, error) {
	var query string
	// Construct the SQL query to fetch the last row from the specified table
	if filter != nil {
		query = fmt.Sprintf(`SELECT * FROM '%s' WHERE %s ORDER BY datetime DESC LIMIT 1`, table, strings.Join(filter, ","))
	} else {
		query = fmt.Sprintf("SELECT * FROM '%s' ORDER BY datetime DESC LIMIT 1", table)
	}
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
func LogData(table string, EventParams []string, values []string, state string) error {
	var err error
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

	// Skip the header row
	if len(records) > 0 {
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

	// Parse the date strings into time.Time events
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

func MoveFile(pathFile, pathTo string) error {    
	inputFile, err := os.Open(pathFile)
    if err != nil {
        return fmt.Errorf("Couldn't open source file: %s", err)
    }
    outputFile, err := os.Create(pathTo)
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
    err = os.Remove(pathFile)
    if err != nil {
        return fmt.Errorf("Failed removing original file: %s", err)
    }
    return nil
}

func ConverToEventDate(layout string, eventTime time.Time) (string) {
	return time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 13, 0, 0, 0, eventTime.Location()).Format(EventTimeFormat)
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
		fmt.Println("Error saving file:", err)
		return
	}
}
