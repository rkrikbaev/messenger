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

// Type1Events = []string{"407001","407002","407003"}
// Type2Events = []string{"407004","407005","407006","407007","407008","407009"}
// Type1Params = []string{"datetime","density","massflowbegin","massflowend","mass"}
// Type2Params = []string{"datetime","density","volume","temperature","tankLevel","mass"}

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

	"encoding/json"

	"fmt"
	"net/http"
	"os"
	"path"

	// "reflect"
	"strings"
	"time"

	"github.com/gokalkan/gokalkan"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"main/modules/filelogger"
	"main/modules/httplogger"
	"main/modules/utils"
)

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
						"tankLevel": 0,
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
						"tankLevel": 0,
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
						"tankLevel": 0,
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
						"tankLevel": 0,
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
						"tankLevel": 0,
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
						"tankLevel": 0,
						"mass": 0
					}
				]`
		)

var (
	EventTime string
	EventRecordDate time.Time
	events []map[string]interface{}
	FilesDirPath string
	// extension string
	// mailReceiver string
	db     *sql.DB
	Events Event
	
	serviceId string
	senderId string
	senderPassword string

	cwd = ""
	dbPath string
	certPath string
	certPassword string
	randomUUID string
	location *time.Location
	logOver = "file"
	servers = []string{"localhost"}
)

func main() {
	var err error
	fmt.Println("Start")
	cwd, _ = os.Getwd()
	// cwd = "/Users/rustamkrikbayev/projects/parser/app"
	dbPath = path.Join(cwd, "./app/db/", dbname)
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	env := path.Join(cwd, "./app/app.env")
	err = godotenv.Load(env)
	if err != nil {
		fmt.Println("error parse env file: ", env)
	}
	certPath = os.Getenv("certPath")
	certPassword = os.Getenv("certPassword")
	serviceId = os.Getenv("serviceId")
	senderId = os.Getenv("senderId")
	senderPassword = os.Getenv("senderPassword")
	// mailReceiver = os.Getenv("mailReceiver")

	logOver = os.Getenv("logOver")
	servers = strings.Split(os.Getenv("httpServers"),",")

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
		if logOver == "file" {
			fmt.Println("LogFiles")
			err = LogFiles()
			if err != nil {
				fmt.Println("Error when log Data:", err)
			}
		}

		if logOver == "http" {
			fmt.Println("LogFiles")
			for _, server:= range servers {
				err = logHttp(server)
				if err != nil {
					fmt.Println("Error when log Data:", err)
				}
			}
		}

		fmt.Print("processEvent")
		err = processEvent()
		if err != nil {
			fmt.Println("Error when processing Event:", err)
		}		
		time.Sleep(600 * time.Second)
	}
}

// -----------------------Log data from Http server--------------

func logHttp(address string) error {
	
	url := fmt.Sprintf("http://%s:8765/tags",address)

	jsonData, err := httpclient.FetchData(url)
	if err != nil {
		log.Fatalf("Error fetching data: %s", err)
	}

	data, err := utils.ParseFields(jsonData)
	if err != nil {
		log.Fatalf("Error parsing selected fields: %s", err)
	}

    replacements := map[string]string{
        "ayraq_dev1.mass": "400703.mass",
        // Добавьте дополнительные замены здесь
    }

    replacedData := utils.ReplaceKeys(data, replacements)
	var values []string
	var columns []string
	var eventid string
    for key, value := range replacedData {
        log.Printf("%s: %f\n", key, value)
		eventid = strings.Split(key, ".")[0]
		column := strings.Split(key, ".")[1]
		values = append(values, fmt.Sprintf("%f", value))
		columns = append(columns, column)
    }

	fmt.Println(eventid)
	fmt.Println(columns)
	fmt.Println(values)

	err = InsertDataIntoDB(eventid, columns, values, "")
	if err != nil {
		fmt.Println("Error when insert into db: ", values)
	}

	return nil
}

//------------------------Log data from CSV files

func LogFiles() error {
    path := path.Join(cwd, "./app/csv_data/")
	validPrefixes := []string{"407001","407002","407003", "407004","407005","407006","407007","407008","407009"}

	files, err := os.ReadDir(path)
    if err != nil {
        return err
    }

    for _, file := range files {
        fileName := file.Name()

        if !file.IsDir() && strings.HasSuffix(fileName, ".csv") {
            eventid, isValid := utils.GetValidPrefix(fileName, validPrefixes)
            if isValid {
                pathFile := filepath.Join(path, fileName)

				//parse file to list
				columns, values, err := filelogger.ParseAndPrepareData(pathFile)
				if err != nil {
					fmt.Println("Error when parse file: ", fileName)
				}
				fmt.Println(eventid)
				fmt.Println(columns)
				fmt.Println(values)
				err = InsertDataIntoDB(eventid, columns, values, "")
				if err != nil {
					continue
				}
				pathTo := fmt.Sprintf("%s/saved/%s",path,fileName)
				fmt.Println(pathTo)
				utils.MoveFile(pathFile, pathTo)
			}
		}
	}
	return err
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
	randomstring := utils.GenerateRandomString()
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
		err = InsertDataIntoDB("DOCUMENTS", []string{"document", "datetime"}, []string{xmlString, eventDate}, "STDBY")
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
		fmt.Printf("ERROR, new kalkan client create error: %s", err)
		return err
	}
	defer cli.Close()

	//sign message
	randomUUID = utils.GenerateRandomString()
	err = cli.LoadKeyStore(certPath, certPassword)
	if err != nil {
		fmt.Printf("ERROR, new cli.LoadKeyStore error: %s", err)
		return err
	}

	message, err := cli.SignWSSE(xmlstring, fmt.Sprintf("id-%s", randomUUID))
	if err != nil {
		fmt.Printf("ERROR, new cli.SignWSSE sign error: %s", err)
		return err
	}
	// message := "uuuuuuuuuuuuuuu"
	
	fmt.Println(randomUUID)
	// fmt.Println(message)

	destPath := fmt.Sprintf("%s/xml_data/message_%s.xml",cwd,eventDate)
	utils.SaveToFile(destPath, message)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
	fmt.Println(err)
	return err
	}
	fmt.Println(string(body))
	return nil
}

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
func InsertDataIntoDB(table string, EventParams []string, values []string, state string) error {
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
