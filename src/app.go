package main

// //СИКН Амангельды    407001
// //СИКН Айракты   407002
// //СИКН Жаркум   407003
// //
// //Резервуар Амангельды - 1   407004
// //Резервуар Амангельды - 2   407005
// //
// //Резервуар Айракты - 1   407006
// //Резервуар Айракты - 2   407007
// //
// //Резервуар Жаркум - 1   407008
// //Резервуар Жаркум - 2   407009


import (
	// "bytes"
	"crypto/tls"
	"database/sql"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
	

	"github.com/gokalkan/gokalkan"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"

	"postgresdb"
	"utils"
)

// Event struct represents the structure of each device data
type Device struct {
	ID              int `yaml:"id"`
	DeviceTypeID    int    `yaml:"deviceTypeId"`
	OperationTypeID int    `yaml:"operationTypeId"`
	DeviceNameID    int    `yaml:"deviceNameId"`
	ProductTypeID   int    `yaml:"productTypeId"`
	PipelineID      int    `yaml:"pipelineId,omitempty"`
	Parameters      []string	`yaml:"parameters"`
}

type Devices struct {
    Settings []Device
}

type SendMessage struct {
	XMLName xml.Name `xml:"ns2:SendMessage"`
	XMLns string `xml:"xmlns:ns2,attr"`
	Request Request `xml:"request"`
}

type Request struct {
	RequestInfo   RequestInfo `xml:"requestInfo"`
	RequestData   RequestData `xml:"requestData"`
}

type RequestInfo struct {
	MessageID   string `xml:"messageId"`
	ServiceID   string `xml:"serviceId"`
	MessageDate string `xml:"messageDate"`
	Sender      Sender `xml:"sender"`
}

type Sender struct {
	SenderID string `xml:"senderId"`
	Password string `xml:"password"`
}

type RequestData struct {
	Data Data `xml:"data"`
}

type Data struct {
	XMLName xml.Name `xml:"data"`
	XMLNsCs string   `xml:"xmlns:cs,attr"`
	XMLNsXsi string  `xml:"xmlns:xsi,attr"`
	XsiType string   `xml:"xsi:type,attr"`
	Events []interface{} 	`xml:"events"`
}

type EventType1 struct {
	ID              int     `xml:"id"`
	DateTime        string  `xml:"datetime"`
	DeviceTypeID    int     `xml:"deviceTypeId"`
	OperationTypeID int     `xml:"operationTypeId"`
	DeviceNameID    int     `xml:"deviceNameId"`
	ProductTypeID   int     `xml:"productTypeId"`
	PipelineID      int     `xml:"pipelineId"`
	MassFlowBegin   float64 `xml:"massflowbegin"`
	MassFlowEnd     float64 `xml:"massflowend"`
	Temperature     float64 `xml:"temperature"`
	Density         float64 `xml:"density"`
	Mass            float64 `xml:"mass"`
	Volume          float64 `xml:"volume"`
}

type EventType2 struct {
	ID              int     `xml:"id"`
	DateTime        string  `xml:"datetime"`
	DeviceTypeID    int     `xml:"deviceTypeId"`
	OperationTypeID int     `xml:"operationTypeId"`
	DeviceNameID    int     `xml:"deviceNameId"`
	ProductTypeID   int     `xml:"productTypeId"`
	Temperature     float64 `xml:"temperature"`
	Density         float64 `xml:"density"`
	Mass            float64 `xml:"mass"`
	Volume          float64 `xml:"volume"`
	TankLevel       float64 `xml:"tankLevel"`
}

const (
	recordsTable = "DOCUMENTS"
	EventTimeFormat = "2006-01-02T15:04:05.000-07:00"
		)

var (
	EventTime string
	EventRecordDate time.Time
	FilesDirPath string
	db     *sql.DB
	DB_NAME = ""
	url string
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

	hour = 17
)

func main() {
	
	var err error
	fmt.Println("Start")

	cwd, _ = os.Getwd()
	env := path.Join(cwd, "./.env")
	err = godotenv.Load(env)
	
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Read YAML file
	CONFIG_FILE := os.Getenv("CONFIG_FILE")
	fmt.Println(CONFIG_FILE)

	yamlFile, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Unmarshal YAML data into a Devices struct
	err = yaml.Unmarshal(yamlFile, &devices)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}
	fmt.Println("Read YAML file: ", devices.Settings[0].ID)
	
	DB_HOST:=os.Getenv("DB_HOST")
	DB_PORT:=os.Getenv("DB_PORT")
	DB_USER:=os.Getenv("DB_USER")
	DB_PASSWORD:=os.Getenv("DB_PASSWORD") 
	DB_NAME=os.Getenv("DB_NAME")

	if DB_NAME == "" {
		fmt.Println("DB_NAME is empty")
		return
	}

	// connect to db
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	fmt.Print(dsn)

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Open database error: %v\n", err)
	}

	schemas, err := postgresdb.ListSchemas(db)
	if err != nil {
		fmt.Println("Schemas missing Error:", err)
	} else {
		fmt.Println("Schemas:", schemas)
	}

	defer db.Close()

	// set cert path
	certPath = os.Getenv("certPath")
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
	h := os.Getenv("hour")
    hour, err = strconv.Atoi(h)
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

	fmt.Println("Start run()")
	var err error
	
	// run loop
	for {
		processEvent()
		if err != nil {
			fmt.Println("Error when processing Event:", err)
		}		
		time.Sleep(time.Duration(timesleep) * time.Second)
	}
}		

//------------------------Prepare XML docement
func processEvent() {

	fmt.Println("Start processEvent()")
	// var err error

	date1 := time.Now().In(location)
	fmt.Println("Current time:", date1)

	if date1.Hour() > hour {
		
		// get last recordset with state STDBY
		record, _ := postgresdb.Select(db, DB_NAME, "messenger", recordsTable, []string{"date", "document"}, "date", "DESC", 1, []string{"state='STDBY'"})

		// send all available XML documents after 17:00
		if len(record) > 0 {
			err := SendDoc(record["document"], record["date"])
			if err != nil {
				fmt.Println("Error when call SendMessage():", err)
			}
		} else {
			fmt.Println("No data for send")
		}

	} else {

		fmt.Println("Waiting...")

		tables, _ := postgresdb.ListTables(db, DB_NAME,"logger")
		fmt.Println("Tables:", tables)

		dt, err := postgresdb.FindDiff(db, tables)
		if err != nil {
			fmt.Println("Error when call FindDiff():", err)
		}
		if len(dt) == 0 {
			fmt.Println("No data for create XML")
		}
		for _, d := range dt {
			err = createDocument(d)
			if err != nil {
				fmt.Println("Error when call createDocumentXML():", err)
			}
		}
	}
}


//------------------------Create XML docement
// Create XML document for each day after last recordset
func createDocument(event_dt string) (error) {

	fmt.Println("Create XML document for date: ", event_dt)
	randomstring := utils.GenerateRandomString()

	events, _ := getProccessData(event_dt, devices)
	if len(events) == 0 {
		fmt.Println("No data for this date")
		return nil
	}
	fmt.Println("Data of all events:", events)

	todayDate := time.Now().In(location).Format(EventTimeFormat)

	// Create the struct to represent the XML

			// <ns2:SendMessage xmlns:ns2="http://bip.bee.kz/SyncChannel/v10/Types"> 
			// 	<request> 
			// 		<requestInfo> 
			// 			<messageId>04e0ee99-407a-4b1e-a702-90230f9d7f8f</messageId> 
			// 			<serviceId>ISUN_Service2</serviceId> 
			// 			<messageDate>2024-01-18T18:13:25.583+06:00</messageDate> 
			// 			<sender> <senderId>amangeldygas</senderId> 
			// 			<password>Amangeldy2023</password> 
			// 			</sender> 
			// 		</requestInfo> 
			// 		<requestData> 
			// 			<data xmlns:cs="http://message.persistence.interactive.nat" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="cs:Request">
			// 				<events>
			// 					<id>407001</id>
			// 					<datetime>2024-01-11T13:00:00.000+00:00</datetime>
			// 					<operationTypeId>3</operationTypeId>
			// 					<productTypeId>1</productTypeId>
			// 					<pipelineId>3</pipelineId>
			// 					<deviceTypeId>2</deviceTypeId>
			// 					<deviceNameId>1</deviceNameId>
			// 					<density>716.1492</density>
			// 					<volume>0</volume>
			// 					<mass>80.46875</mass>
			// 					<temperature>0</temperature>
			// 					<massflowbegin>132856.1</massflowbegin>
			// 					<massflowend>132936.6</massflowend>
			// 				</events> 
			// 				<events>
			// 					<id>407004</id>
			// 					<datetime>2024-01-11T13:00:00.000+00:00</datetime>
			// 					<productTypeId>1</productTypeId>
			// 					<deviceTypeId>1</deviceTypeId>
			// 					<operationTypeId>7</operationTypeId>
			// 					<deviceNameId>2</deviceNameId>
			// 					<temperature>0.669922</temperature>
			// 					<density>709.065</density>
			// 					<volume>106.7144</volume>
			// 					<tankLevel>124.9048</tankLevel>
			// 					<mass>75.66756</mass>
			// 				</events> 
			// 				<events><id>407005</id><datetime>2024-01-11T13:00:00.000+00:00</datetime><deviceTypeId>1</deviceTypeId><operationTypeId>7</operationTypeId><deviceNameId>2</deviceNameId><productTypeId>1</productTypeId><temperature>0.242188</temperature><density>705.322</density><volume>52.64053</volume><tankLevel>61.60059</tankLevel><mass>37.12854</mass></events> 
			// 				<events><id>407002</id><datetime>2024-01-11T13:00:00.000+00:00</datetime><deviceTypeId>2</deviceTypeId><deviceNameId>1</deviceNameId><productTypeId>1</productTypeId><pipelineId>3</pipelineId><operationTypeId>3</operationTypeId><volume>0</volume><mass>0</mass><temperature>0</temperature><density>0</density><massflowbegin>0</massflowbegin><massflowend>0</massflowend></events> 
			// 				<events><id>407006</id><datetime>2024-01-11T13:00:00.000+00:00</datetime><operationTypeId>7</operationTypeId><productTypeId>1</productTypeId><deviceTypeId>1</deviceTypeId><deviceNameId>2</deviceNameId><mass>2.792</mass><temperature>0.5976563</temperature><density>742.4</density><volume>3.766</volume><tankLevel>55.875</tankLevel></events> 
			// 				<events><id>407007</id><datetime>2024-01-11T13:00:00.000+00:00</datetime><deviceTypeId>1</deviceTypeId><deviceNameId>2</deviceNameId><productTypeId>1</productTypeId><operationTypeId>7</operationTypeId><tankLevel>170.4625</tankLevel><mass>13.424</mass><temperature>3</temperature><density>732</density><volume>18.192</volume></events> 
			// 				<events><id>407003</id><datetime>2024-01-11T13:00:00.000+00:00</datetime><deviceNameId>1</deviceNameId><pipelineId>3</pipelineId><deviceTypeId>2</deviceTypeId><operationTypeId>3</operationTypeId><productTypeId>1</productTypeId><temperature>0</temperature><density>763.0197</density><massflowbegin>1457.847</massflowbegin><massflowend>1457.847</massflowend><mass>0.0001220703</mass><volume>0</volume></events> 
			// 				<events><id>407008</id><datetime>2024-01-11T13:00:00.000+00:00</datetime><deviceNameId>2</deviceNameId><productTypeId>1</productTypeId><deviceTypeId>1</deviceTypeId><operationTypeId>7</operationTypeId><temperature>6.021851</temperature><density>732</density><volume>77.66721</volume><tankLevel>233.6</tankLevel><mass>54.528</mass></events> 
			// 				<events><id>407009</id><datetime>2024-01-11T13:00:00.000+00:00</datetime><operationTypeId>7</operationTypeId><deviceTypeId>1</deviceTypeId><deviceNameId>2</deviceNameId><productTypeId>1</productTypeId><temperature>5.786926</temperature><density>758.4</density><volume>95.8464</volume><tankLevel>307.6</tankLevel><mass>72.704</mass></events>
			// 			</data> 
			// 		</requestData> 
			// 	</request> 
			// </ns2:SendMessage>
			
	ifaceEvents := make([]interface{}, 0)
	ifaceEvents = append(ifaceEvents, events...)

	b := SendMessage{		
						XMLName: xml.Name{Local: "ns2:SendMessage"},
						XMLns: "http://bip.bee.kz/SyncChannel/v10/Types",
						Request: Request{
							RequestInfo: RequestInfo{
								MessageID:   randomstring,
								ServiceID:   serviceId,
								MessageDate: todayDate,
								Sender: Sender{
									SenderID: senderId,
									Password: senderPassword,
								},
							},
							RequestData: RequestData{
								Data: Data{
									XMLNsCs: "http://message.persistence.interactive.nat",
									XMLNsXsi: "http://www.w3.org/2001/XMLSchema-instance",
									XsiType: "cs:Request",
									Events: ifaceEvents,
								},
							},
						},
					}


	// Marshal the struct into XML
	XMLSendMessage, err := xml.MarshalIndent(b, "", "    ")
	if err != nil {
		fmt.Printf("Error marshalling XML: %v", err)
		return err
	}

	// Print the XML
	fmt.Print("XMLSendMessage: ")
	fmt.Println(string(XMLSendMessage))

	// Write the XML to a DB
	columns := []string{"document", "date", "created", "state"}
	values := []string{ string(XMLSendMessage), event_dt, time.Now().Format(EventTimeFormat), "STDBY" }

	err = postgresdb.Insert(db, DB_NAME, "messenger", recordsTable, columns, values)
	if err != nil {
		fmt.Println("Error when call Insert():", err)
		return err
	}
	time.Sleep(1 * time.Second)
	return nil
}
// Process data for each device
func getProccessData(event_date string, objects Devices) ([]interface{}, error) {

	fmt.Println("Get event data on: ", event_date)

	var arr []interface{}
	
	for _, object := range objects.Settings {

		fmt.Println("Generate XML for device: ", object.ID)
		data, err := getEventData(event_date, object)

		if err != nil {	
			fmt.Println("Error when call getEventData():", err)
			return 	nil, err
		}
		arr = append(arr, data)

	}
	return arr, nil
}
// Get data for each device
func getEventData(EventDate string, object Device) (interface{}, error) {

	//get first record with event ID
	filter := []string{fmt.Sprintf("datetime = '%s'", EventDate)}
	columns := object.Parameters
	table := fmt.Sprintf("%d", object.ID)
	data, err := postgresdb.Select(db, DB_NAME, "logger", table, columns, "datetime",  "ASC", 1, filter)
	if err != nil{
		fmt.Println("Error when call Select():", err)
		return nil, err
	}
	if data == nil {
		fmt.Println("No data for this date")
		return nil, nil
	}

	fmt.Println(object.OperationTypeID)
	var event interface{}
	if object.OperationTypeID == 3 {
		event = EventType1{
			ID: object.ID,
			DateTime: EventDate,
			DeviceTypeID: object.DeviceTypeID,
			OperationTypeID: object.OperationTypeID,
			DeviceNameID: object.DeviceNameID,
			ProductTypeID: object.ProductTypeID,
			PipelineID: object.PipelineID,
			Temperature: parseFloat(data["temperature"]),
			Density: parseFloat(data["density"]),
			Volume: parseFloat(data["volume"]),
			MassFlowBegin: parseFloat(data["massflowbegin"]),
			MassFlowEnd: parseFloat(data["massflow end"]),
			Mass: parseFloat(data["mass"])}
		return event, nil
	}

	if object.OperationTypeID == 7 {
		event = EventType2{
			ID: object.ID,
			DateTime: EventDate,
			DeviceTypeID: object.DeviceTypeID,
			OperationTypeID: object.OperationTypeID,
			DeviceNameID: object.DeviceNameID,
			ProductTypeID: object.ProductTypeID,
			Temperature: parseFloat(data["temperature"]),
			Density: parseFloat(data["density"]),
			Volume: parseFloat(data["volume"]),
			Mass: parseFloat(data["mass"]),
			TankLevel: parseFloat(data["tankLevel"])}
		return event, nil
	}
	return nil, nil
}

//------------------------Send SOAP message
func SendDoc(xmlstring, eventDate string) (error) {
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

	doc, err := cli.SignWSSE(xmlstring, fmt.Sprintf("id-%s", randomUUID))
	if err != nil {
		fmt.Printf("ERROR, new cli.SignWSSE sign error: %s", err)
		return err
	}
	
	fmt.Println(randomUUID)

	destPath := fmt.Sprintf("%s/xml_data/message_%s.xml",cwd,eventDate)
	utils.SaveToFile(destPath, doc)
	
	// send message
	err = sendRequest(doc)

	if err != nil {
		fmt.Println("Error when call sendRequest():", err)
	} else {
		state = "SUCCESS"
	}

	// update state
	columns := []string{"state", "sent"}
	values := []string{state, time.Now().Format(EventTimeFormat)}
	filter := fmt.Sprintf("date = '%s'", eventDate)

	err = postgresdb.Update(db, DB_NAME, "messenger", recordsTable, columns, values, filter, state)
	if err != nil {
		fmt.Println("Error update data after success send message", err)
		return err
	}

return nil
}
// sendRequest sends a request to the server
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

// parseFloat converts a string to a float64
func parseFloat(s string) float64 {
	if s == "" {
		return 0.0
	}
	if s == "<nil>" {
		return 0.0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println("Error parsing float:", err)
	}
	return f
}

