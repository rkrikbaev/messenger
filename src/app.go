package main

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

import (
	"bytes"
	"crypto/tls"

	// "context"
	"database/sql"
	"io"
	"log"
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

	"utils"
	"dbsql"
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
	// dbname = "db3.db"
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

	certPath string
	certPassword string
	randomUUID string
	location *time.Location
)

func main() {
	
	var err error
	fmt.Println("Start")

	env := path.Join(cwd, "./.env")
	err = godotenv.Load(env)
	
	if err != nil {
		fmt.Println("error parse env file: ", env)
	}
	cwd, _ = os.Getwd()

	dbPath := os.Getenv("db_path")
	fmt.Println(dbPath)

	db, err = sql.Open("sqlite3", dbPath)
	
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	certPath = os.Getenv("certPath")
	certPassword = os.Getenv("certPassword")
	serviceId = os.Getenv("serviceId")
	senderId = os.Getenv("senderId")
	senderPassword = os.Getenv("senderPassword")
	// mailReceiver = os.Getenv("mailReceiver")

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

		fmt.Print("processEvent")
		err = processEvent()
		if err != nil {
			fmt.Println("Error when processing Event:", err)
		}		
		time.Sleep(600 * time.Second)
	}
}		

//------------------------Prepare XML docement
func processEvent() error {
	var err error
	// var eventDate string
	// Try to find latest recordset with state STDBY
	// And sebd it after 17:00
	eventRecordSet, _ := dbsql.Select(db, "DOCUMENTS", []string{"state='STDBY'"})

	if len(eventRecordSet) > 0 {
		date1 := time.Now().In(location)
		date2 := time.Date(date1.Year(), date1.Month(), date1.Day(), 17,0,0,0, date1.Location())
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
		eventRecordSet, _ = dbsql.Select(db, "DOCUMENTS", "ASC", 1, nil)
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
		err = dbsql.Insert("DOCUMENTS", []string{"document", "datetime"}, []string{xmlString, eventDate}, "STDBY")
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

	data, err := dbsql.Select(eventId, filter)
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

	columns := []string{"state", "sentMessageDate"}
	values := []string{state, time.Now().Format(EventTimeFormat)}

	err = sendRequest(message)
	if err != nil {
		fmt.Println("Error creating request:", err)
		err = dbsql.Update(table, columns, values, condition, state)
		return err 
	}
	state = "SUCCESS"
	err = dbsql.Update(table, condition, state)
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