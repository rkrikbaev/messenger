package main

import (
	"bytes"
	"database/sql"
	"log"

	// "database/sql"
	"encoding/json"

	// "encoding/xml"
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
	ip = "127.0.0.1"
	port = "8080"
	ReportTable = "REPORT"
	dbname = "app/db/app.db"
)

var (
	// certPath string
	// certPassword string
	serviceId	string
	senderId	string
	senderPassword	string
	// mailReceiver string
	
	EventTime string
	EventRecordDate time.Time
	
	// documentId string
	objects map[string]interface{}
	FilesDirPath string
	extension string

	db     *sql.DB

	Events Event
	allTabels = []string{"407001","407002","407003", "407004","407005","407006","407007","407008","407009"}
)

func main() {
	// Read enviroiments
	fmt.Println("Start")
	//
	var err error
	db, err = sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Close the database connection at the end of the program

	// Load the environment variables from .env file
	cwd, _ := os.Getwd()
	envFileDirPath := path.Join(cwd, "app.env")
	// fmt.Println("parse env file: ", err)
	// fmt.Println("parse env file: ", envFileDirPath)
	err = godotenv.Load(envFileDirPath)
	if err != nil {
		fmt.Println("error parse env file: ", envFileDirPath)
	}
	serviceId = os.Getenv("serviceId")
	senderId = os.Getenv("senderId")
	senderPassword = os.Getenv("senderPassword")
	// mailReceiver = os.Getenv("mailReceiver")

	FilesDirPath = os.Getenv("DIR")
	extension = os.Getenv("EXT")
	// dbname = os.Getenv("DB_NAME")
	json_data, err := ioutil.ReadFile("app/src/data.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	err = json.Unmarshal(json_data, &objects)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	run()
}

func run() {
	TodayDatetime := time.Now()
	for {
		var err error
		if TodayDatetime.After(EventRecordDate) {
		
		err = LogDataFromFile(objects)
		if err != nil {
			fmt.Println("Error when call LogDataFromFile():", err)
			return
		}
		err = CreateXML()
		if err != nil {
			fmt.Println("Error when call CreateXML():", err)
			return
		}
		err = SendMessage()
		if err != nil {
			fmt.Println("Error when call SendMessage():", err)
			return
		}
		} else {
			break
		}
		time. Sleep(10 * time.Second)
	}
}
//------------------------Log CSV files


//------------------------Prepare XML docement
func CreateXML() error {
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
	xmlStringBody := xmlBuffer.String()
	err = LogData("DOCUMENTS", []string{"document", "datetime"}, []string{xmlStringBody, eventDate})
	if err != nil {
		return err
	}	
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
						fmt.Println(errMsg2)
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
func SendMessage() error {
	state := "FAIL"
	// для теста
	// opts := gokalkan.OptsTest

	// для прода
	opts := gokalkan.OptsProd
	
	cli, err := gokalkan.NewClient(opts...)
	if err != nil {
		fmt.Print(fmt.Sprintf("ERROR, new kalkan client create error: %s", err))
		// errMsg := fmt.Sprintf("ERROR, new kalkan client create error: %s", err)
		// logging(errMsg)
		//sendMail( mailReceiver, ERROR_WITH_KALKAN_MAIL_SUBJECT, errMsg)
		return err
	}
	// Обязательно закрывайте клиент, иначе приведет к утечкам ресурсов
	defer cli.Close()

	// Construct the SQL query with placeholders
	filter := []string{"state=='STDBY'"}
	data, _ := getFirstRecord("DOCUMENTS", filter)
	if data == nil {
		return nil
	}
	xmlString := data["document"]
	if xmlString == "" {
		return nil
	}
	eventDate := data["datetime"]
	if eventDate == "" {
		return nil
	}

	//sign message
	randomUUID := generateRandomString()
	err = cli.LoadKeyStore(certPath, certPassword)
	if err != nil {
		return err
	}
	message, err := cli.SignWSSE(xmlString, fmt.Sprintf("id-%s", randomUUID))

	status, err := sendRequest(ip, port, message)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err 
	}
	if status == "200 OK" {
		state = "SUCCESS"
	}
	table := "DOCUMENTS"
	err = updateData(table, eventDate, state)
	if err != nil {
		fmt.Println("Error update data after success send message", err)
		return err 
	}
	destPath := fmt.Sprintf("app/xml_data/message_%s.xml",eventDate)
	content := message
	SaveToFile(destPath, content)
return nil
}

func sendRequest(ip, port, message string) (string, error) {

	// Set the URL and HTTP method
	url := fmt.Sprintf("http://%s:%s/bip-sync-wss-gost/", ip, port)
	method := "POST"
	payload := strings.NewReader(message)

	// Create a request with the SOAP message
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err 
	}
	// Set headers for the SOAP request
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	// Make the HTTP request
	res, err := client.Do(req)
	if err != nil {
		return res.Status, err 
	}
	defer res.Body.Close()
	
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return res.Status, err 
	}
	// fmt.Println("Response Body:", res.Status)
	return res.Status, nil
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

// //----------------------------- // ------------------------------------------ //

func generateRandomString() string {
	uuid := uuid.New()
	return uuid.String()
}

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
// 		fmt.Println("Error opening file:", err)
// 	}
// 	defer file.Close()

// 	// Append text to the file
// 	_, err = io.WriteString(file, logText)
// 	if err != nil {
// 		fmt.Println("Error writing to file:", err)
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
// 			fmt.Println("Error creating file:", err)
// 			return
// 		}
// 		defer file.Close()

// 		fmt.Println("File created:", filePath)
// 	}
// }