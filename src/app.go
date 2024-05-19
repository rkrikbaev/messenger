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
	"database/sql"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	// "reflect"
	"strconv"
	"unicode"

	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gokalkan/gokalkan"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"

	"postgresdb"
	"utils"

	"github.com/fatih/structs"
)

type Background struct {
	XMLName string `xml:"background"`
	Items   []interface{}
}

// Event struct represents the structure of each device data
type Event struct {
	ID              string `yaml:"id"`
	DeviceTypeID    int    `yaml:"deviceTypeId"`
	OperationTypeID int    `yaml:"operationTypeId"`
	DeviceNameID    int    `yaml:"deviceNameId"`
	ProductTypeID   int    `yaml:"productTypeId"`
	PipelineID      int    `yaml:"pipelineId,omitempty"`
	Parameters      []string	`yaml:"parameters"`
}

// Devices struct represents the structure of the data file
type Devices struct {
	Data []Event
}

type Document struct {
	ID int
	Document string
	EventDate string
	Created string
	MessageSent string
	State string
}

const (
	recordsTable = "DOCUMENTS"
	EventTimeFormat = "2006-01-02T15:04:05.000-07:00"
		)

var (
	EventTime string
	EventRecordDate time.Time
	events []map[string]interface{}
	FilesDirPath string
	db     *sql.DB
	DB_NAME = "postgres"
	url string
	Events Event
	devices Devices
	serviceId string
	senderId string
	senderPassword string

	cwd = ""
	timesleep = 600

	certPath string
	certPassword string
	randomUUID string
	location *time.Location

	hourSend = 17
)

func main() {
	
	var err error
	fmt.Println("Start")

	env := path.Join(cwd, "./.env")
	err = godotenv.Load(env)
	
	if err != nil {
		fmt.Println("error parse env file: ", env)
	}

	// Read YAML file
	SETTINGS_FILE := os.Getenv("SETTINGS")
	yamlFile, err := ioutil.ReadFile(SETTINGS_FILE)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Unmarshal YAML data into a Devices struct
	err = yaml.Unmarshal(yamlFile, &devices)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	cwd, _ = os.Getwd()

	dbPath := os.Getenv("db_path")
	fmt.Println(dbPath)
	
	DB_HOST:=os.Getenv("DB_HOST")
	DB_PORT:=os.Getenv("DB_PORT")
	DB_USER:=os.Getenv("DB_USER")
	DB_PASSWORD:=os.Getenv("DB_PASSWORD") 
	DB_NAME:=os.Getenv("DB_NAME")

	
	db, err = postgresdb.Connect(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	if err != nil {
		log.Fatal(err)
	}
	// set cert path
	certPath = os.Getenv("certPath")
	// set cert password
	certPassword = os.Getenv("certPassword")
	serviceId = os.Getenv("serviceId")
	senderId = os.Getenv("senderId")
	senderPassword = os.Getenv("senderPassword")
	
	// set url
	url = os.Getenv("isun_url")

	// set time to sleep
	s := os.Getenv("timesleep")
    timesleep, err = strconv.Atoi(s)
    if err != nil {
        panic(err)
    }

	// set hour to send message
	h := os.Getenv("hourSendMessage")
    timesleep, err = strconv.Atoi(h)
    if err != nil {
        panic(err)
    }

	// set location
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
		time.Sleep(time.Duration(timesleep) * time.Second)
	}
}		

//------------------------Prepare XML docement
func processEvent() error {
	var err error
	// var eventDate string
	// Try to find latest recordset with state STDBY
	// And sebd it after 17:00

	// get last recordset with state STDBY
	filter := []string{"state='STDBY'"}
	eventRecordSet, _ := postgresdb.Select(db, DB_NAME, "messenger", recordsTable, []string{"document, datetime"}, "DESC", 1, filter)

	if len(eventRecordSet) > 0 {
		date1 := time.Now().In(location)
		date2 := time.Date(date1.Year(), date1.Month(), date1.Day(), hourSend,0,0,0, date1.Location())
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
		err = createDocumentXML()
		if err != nil {
			fmt.Println("Error when call createDocumentXML():", err)
		}
	}
	return nil
}

//------------------------Create XML docement
// Create XML document for each day after last recordset
func createDocumentXML() (error) {

	var eventRecordSet map[string]string
	var existEventDate string
	var err error
	var randomstring string

	randomstring = utils.GenerateRandomString()

	// Try to create XML document after 13:00 wyen process data come
	// Collect data for each day after last recordset
	fmt.Println("Collecting data...")
	eventRecordSet, _ = postgresdb.Select(db, DB_NAME, "messenger", recordsTable, []string{"datetime"}, "ASC", 1, nil)


	if len(eventRecordSet) > 0 {
		existEventDate = eventRecordSet["datetime"]
	} else {
		return fmt.Errorf("No data in the database")
	}

	exDate, err := time.Parse(EventTimeFormat, existEventDate)
	if err != nil {
		return err
	}

	var todayDate time.Time
	todayDate = time.Now()
	
	for todayDate.After(exDate) {
		newDate := exDate.AddDate(0, 0, 1)
		exDate = newDate
		eventDate := newDate.Format(EventTimeFormat)
		fmt.Print("Try to collect data for new XML document\n")
		fmt.Print(eventDate)
		
		eventData, err := getEventData(eventDate, devices)
		if stringIsEmpty(eventData) {
			continue
		}
		if err != nil {
			fmt.Println("Error when call getEventData():", err)
			continue
		}
		var xmlBuffer bytes.Buffer
		xmlBuffer.WriteString(fmt.Sprintf(
			`<ns2:SendMessage xmlns:ns2="http://bip.bee.kz/SyncChannel/v10/Types"> <request> <requestInfo> <messageId>%s</messageId> <serviceId>%s</serviceId> <messageDate>%s</messageDate> <sender> <senderId>%s</senderId> <password>%s</password> </sender> </requestInfo> <requestData> <data xmlns:cs="http://message.persistence.interactive.nat" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="cs:Request">%s</data> </requestData> </request> </ns2:SendMessage>`,
			randomstring, serviceId, todayDate.Format(EventTimeFormat), senderId, senderPassword, eventData))
		xmlString := xmlBuffer.String()
		columns := []string{"document", "datetime", "collected"}
		values := []string{xmlString, eventDate, time.Now().Format(EventTimeFormat)}
		if recordsTable == "DOCUMENTS" {
			columns = append(columns, "state")
			values = append(values, "STDBY")
		}

		err = postgresdb.Insert(db, DB_NAME, "messenger", recordsTable, columns, values)
		if err != nil {
			fmt.Println("Error when call Insert():", err)
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

func getEventData(eventDate string, objects Devices) (string, error) {

	var eventDataArray []string
	
	for _, object := range objects.Data {

		eventXMLString, err := generateXMLString(eventDate, object)
		
		if err != nil {
			errMsg2 := fmt.Sprintf("Error generate XML with data): %s", err)
			fmt.Println(errMsg2)
			return "", err
		}

		if stringIsEmpty(eventXMLString) {
			fmt.Printf("Empty string for: %s", eventDate)
			return "", nil
		} else {
			eventDataArray = append(eventDataArray, eventXMLString)
		}
	}
	eventData := strings.Join(eventDataArray, " ")
	return eventData, nil
}

// generateXMLString generates an XML string based on the provided EventDate and object.
// It returns the generated XML string and an error, if any.
// The EventDate parameter specifies the date of the event.
// The object parameter contains the event details.
// The generated XML string represents the event data in XML format.
// If an error occurs during the generation process, the error is returned.
func generateXMLString(EventDate string, object Event) (string, error) {

	s := &Event{
		ID: object.ID,
		DeviceTypeID: object.DeviceTypeID, 
		OperationTypeID: object.OperationTypeID, 
		DeviceNameID: object.DeviceNameID, 
		ProductTypeID: object.ProductTypeID, 
		PipelineID: object.PipelineID}

	event := structs.Map(s)

	//get first record with event ID
	filter := []string{fmt.Sprintf("datetime = '%s'", EventDate)}
	columns := object.Parameters
	data, err := postgresdb.Select(db, DB_NAME, "logger", object.ID, columns, "ASC", 1, filter)
	if err != nil {
		return ``, err
	}
	if data == nil {
		return "", nil
	}

	// event datetime
	datetime := data["datetime"]
	EventRecordDate, _ = time.Parse(EventTimeFormat,datetime)
	
	for k, v := range data {
    event[k] = v
	}

	bk := Background{ Items: []interface{}{ map[string]interface{}{"events": &event } } }

	// encode xml
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "\t")

	err = enc.Encode(&bk)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	
	return "", nil
}

//------------------------Send SOAP message
func SendMessage(xmlstring, eventDate string) (error) {
	var err error
	var state string

	state = "FAIL"

	// create new kalkan client
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
	
	fmt.Println(randomUUID)

	destPath := fmt.Sprintf("%s/xml_data/message_%s.xml",cwd,eventDate)
	utils.SaveToFile(destPath, message)
	
	// send message
	err = sendRequest(message)

	if err != nil {
		fmt.Println("Error when call sendRequest():", err)
	} else {
		state = "SUCCESS"
	}

	// update state
	columns := []string{"state", "sent"}
	values := []string{state, time.Now().Format(EventTimeFormat)}
	filter := fmt.Sprintf("datetime = '%s'", eventDate)

	err = postgresdb.Update(db, DB_NAME, "messenger", recordsTable, columns, values, filter, state)
	if err != nil {
		fmt.Println("Error update data after success send message", err)
		return err
	}

return nil
}

func sendRequest(data string) error {

	payload := strings.NewReader(data)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client {Transport: tr }
	req, err := http.NewRequest("POST", url, payload)

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