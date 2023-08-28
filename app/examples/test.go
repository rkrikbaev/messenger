package main

import (
	"fmt"
	"encoding/csv"
	"os"
	"time"
	"net/smtp"
	//"strings"
	"io"
	//"io/ioutil"
	"strconv"
	"log"
	//"bytes"
	//"net/http"

	"github.com/gokalkan/gokalkan"
	"github.com/google/uuid"
)


//Расходомер Амангельды    407001
//Расходомер Айракты   407002
//Расходомер Жаркум   407003
//
//Резервуар Амангельды - 1   407004
//Резервуар Амангельды - 2   407005
//
//Резервуар Айракты - 1   407006
//Резервуар Айракты - 2   407007
//
//Резервуар Жаркум - 1   407008
//Резервуар Жаркум - 2   407009

type RashodomerEvent struct {
	ID                string
	DateTime          string
	DeviceTypeID      string
	OperationTypeID   string
	DeviceNameID      string
	ProductTypeID     string
	Temperature       string
	Density           string
	Volume            string
	PipelineID        string
	MassFlowBegin     string
	MassFlowEnd       string
	Mass              string
}

type ReservuarEvent struct {
	ID                string
	DateTime          string
	DeviceTypeID      string
	OperationTypeID   string
	DeviceNameID      string
	ProductTypeID     string
	Density           string
	TankLevel         string
	Temperature       string
	Volume            string
	Mass              string
}


func NewRashodomerEvent(id string) *RashodomerEvent {
	return &RashodomerEvent{
		ID: id,
		DeviceTypeID:       "2",// Default value
		OperationTypeID: "3", // Default value
		DeviceNameID:    "1",// Default value
		ProductTypeID:    "1",// Default value
		PipelineID: "3",
	}
}

func NewReservuarEvent(id string) *ReservuarEvent {
	return &ReservuarEvent{
		ID: id,
		DeviceTypeID:       "1",// Default value
		OperationTypeID: "7", // Default value
		DeviceNameID:    "2",// Default value
		ProductTypeID:    "1",// Default value
	}
}


var (
	certPath = "test_cert/GOSTKNCA.p12" // путь к хранилищу

	certPassword = "Aa1234" // пароль
	// P.S. никогда не храните пароли в коде

	rashod407001 = NewRashodomerEvent("407001");
	rashod407002 = NewRashodomerEvent("407002");
	rashod407003 = NewRashodomerEvent("407003");

	rezer407004 = NewReservuarEvent("407004");
	rezer407005 = NewReservuarEvent("407005");
	rezer407006 = NewReservuarEvent("407006");
	rezer407007 = NewReservuarEvent("407007");
	rezer407008 = NewReservuarEvent("407008");
	rezer407009 = NewReservuarEvent("407009");

	updatedAmangeldi = false
	updatedAirakty = false
	updatedJarkum = false

	//pathToFileAmangeldi = "Amangeldy/"
	//pathToFileAirakty = "Ayraqty/"
	//pathToFileZharkum = "Zharkum/"
	pathToFileAmangeldi = "/app/csv_data/amangeldy/"
	pathToFileAirakty = "/app/csv_data/ayraqty/"
	pathToFileZharkum = "/app/csv_data/zharkum/"

	timeFormatForEvent = "2006-01-02T15:04:05.000-07:00"
	serviceId = "ISUN_Service2"
	senderId = "amangeldygas"
	senderPassword = "Amangeldy2023"

	mailSender   = "amangeldi.isun@mail.ru"
	appPassword  = "h9UcQwTcdK75DVdpyEHY"
	mailReceiver = "ae.aubakirov@gmail.com"

	smtpHost = "smtp.mail.ru"
	smtpPort = 587

	logFileName = "log.txt"

	DATE_STR_FILE = os.Getenv("DATE") //"13.07.2024"
	
)



func main() {
	res := runMain()
	fmt.Println(res)
}

func validateMassResult() string {
	if rashod407001.MassFlowBegin == "" || rashod407001.MassFlowEnd == "" ||
		rashod407002.MassFlowBegin == "" || rashod407002.MassFlowEnd == "" ||
		rashod407003.MassFlowBegin == "" || rashod407003.MassFlowEnd == "" {
		return "ERROR: one of the mass values is empty"
	}

	massBegin1, err := strconv.ParseFloat(rashod407001.MassFlowBegin, 64)
	if err != nil {
		return "ERROR: rashodomer1 massBegin is not a number"
	}
	massEnd1, err := strconv.ParseFloat(rashod407001.MassFlowEnd, 64)
	if err != nil {
		return "ERROR: rashodomer1 massEnd is not a number"
	}
	mass1, err := strconv.ParseFloat(rashod407001.Mass, 64)
	if err != nil {
		return "ERROR: rashodomer1 mass is not a number"
	}
	massBegin2, err := strconv.ParseFloat(rashod407002.MassFlowBegin, 64)
	if err != nil {
		return "ERROR: rashodomer2 massBegin is not a number"
	}
	massEnd2, err := strconv.ParseFloat(rashod407002.MassFlowEnd, 64)
	if err != nil {
		return "ERROR: rashodomer2 massEnd is not a number"
	}
	mass2, err := strconv.ParseFloat(rashod407002.Mass, 64)
	if err != nil {
		return "ERROR: rashodomer2 mass is not a number"
	}
	massBegin3, err := strconv.ParseFloat(rashod407003.MassFlowBegin, 64)
	if err != nil {
		return "ERROR: rashodomer3 massBegin is not a number"
	}
	massEnd3, err := strconv.ParseFloat(rashod407003.MassFlowEnd, 64)
	if err != nil {
		return "ERROR: rashodomer3 massEnd is not a number"
	}
	mass3, err := strconv.ParseFloat(rashod407003.Mass, 64)
	if err != nil {
		return "ERROR: rashodomer3 mass is not a number"
	}

	if (massEnd1 - massBegin1) != mass1 || (massEnd2 - massBegin2) != mass2 || (massEnd3 - massBegin3) != mass3 {
		return "ERROR: mass values are not equal"
	}
	return "OK"
}

//----------------------------- // ------------------------------------------ //

func runMain() string {

	// 1. read file of results sent
	// if file is empty, then send all data and create file with last date
	// if file is not empty, then send data from last date to current date and update file with current date
	// else send all data and create file with last date
	// 2. send data to server
	// if failed then send notification to email
	// 3. update file with last date
	// 4. repeat after 1 hour

	// для теста
	opts := gokalkan.OptsTest

	// для прода
	// opts := gokalkan.OptsProd

	cli, err := gokalkan.NewClient(opts...)
	if err != nil {
		errMsg := fmt.Sprintf("ERROR, new kalkan client create error: %s", err)
		logging(errMsg)
		//sendMail( mailReceiver, ERROR_WITH_KALKAN_MAIL_SUBJECT, errMsg)
		return errMsg
	}
	// Обязательно закрывайте клиент, иначе приведет к утечкам ресурсов
	defer cli.Close()

	xmlToSign := prepareXmlToSign()


	fmt.Println("prepareXmlToSign: %s", xmlToSign) // need to delete

	// Подгружаем хранилище с паролем
	err = cli.LoadKeyStore(certPath, certPassword)
	if err != nil {
		errMsg1 := fmt.Sprintf("ERROR, load key store error: %s", err)
		return errMsg1
	}

	randomUUID := generateRandomString()
	signedXML, err := cli.SignWSSE(xmlToSign, fmt.Sprintf("id-%s", randomUUID))
	 //signedXML, err := cli.SignWSSE(" <ns2:SendMessage xmlns:ns2=\"http://bip.bee.kz/SyncChannel/v10/Types\"> <request> <requestInfo> <messageId>214b4374-4486-456e-834e-ae06a607a70c</messageId> <serviceId>ISUN_Service2</serviceId> <messageDate>2023-06-06T13:45:21.024Z</messageDate> <sender> <senderId>amangeldygas</senderId> <password>Qazaqstan2023</password> </sender> </requestInfo> <requestData> <data xmlns:cs=\"http://message.persistence.interactive.nat\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:type=\"cs:Request\"> <events> <id>407001</id> <datetime>2023-06-07T00:00:00.000+06:00</datetime> <deviceTypeId>2</deviceTypeId> <operationTypeId>3</operationTypeId> <deviceNameId>1</deviceNameId> <productTypeId>1</productTypeId> <temperature>10.79</temperature> <density>52.5</density> <volume>12.2</volume> <pipelineId>3</pipelineId> <massflowbegin>4.10</massflowbegin> <massflowend>8.10</massflowend> <mass>4.0</mass> </events> <events> <id>407002</id> <datetime>2023-06-07T00:00:00.000+06:00</datetime> <deviceTypeId>2</deviceTypeId> <operationTypeId>3</operationTypeId> <deviceNameId>1</deviceNameId> <productTypeId>1</productTypeId> <temperature>10.79</temperature> <density>52.5</density> <volume>12.2</volume> <pipelineId>3</pipelineId> <massflowbegin>4.10</massflowbegin> <massflowend>4.10</massflowend> <mass>0.0</mass> </events> <events> <id>407003</id> <datetime>2023-06-07T00:00:00.000+06:00</datetime> <deviceTypeId>2</deviceTypeId> <operationTypeId>3</operationTypeId> <deviceNameId>1</deviceNameId> <productTypeId>1</productTypeId> <temperature>10.79</temperature> <density>52.5</density> <volume>12.2</volume> <pipelineId>3</pipelineId> <massflowbegin>4.10</massflowbegin> <massflowend>5.25</massflowend> <mass>1.15</mass> </events> <events> <id>407004</id> <datetime>2023-06-07T00:00:00.000+06:00</datetime> <deviceTypeId>1</deviceTypeId> <operationTypeId>7</operationTypeId> <deviceNameId>2</deviceNameId> <productTypeId>1</productTypeId> <density>52.5</density> <tankLevel>12.5</tankLevel> <temperature>12.79</temperature> <volume>12.5</volume> <mass>0.0</mass> </events> <events> <id>407005</id> <datetime>2023-06-07T00:00:00.000+06:00</datetime> <deviceTypeId>1</deviceTypeId> <operationTypeId>7</operationTypeId> <deviceNameId>2</deviceNameId> <productTypeId>1</productTypeId> <density>52.5</density> <tankLevel>12.5</tankLevel> <temperature>12.79</temperature> <volume>12.5</volume> <mass>0.0</mass> </events> <events> <id>407006</id> <datetime>2023-06-07T00:00:00.000+06:00</datetime> <deviceTypeId>1</deviceTypeId> <operationTypeId>7</operationTypeId> <deviceNameId>2</deviceNameId> <productTypeId>1</productTypeId> <density>52.5</density> <tankLevel>12.5</tankLevel> <temperature>12.79</temperature> <volume>12.5</volume> <mass>0.0</mass> </events> <events> <id>407007</id> <datetime>2023-06-07T00:00:00.000+06:00</datetime> <deviceTypeId>1</deviceTypeId> <operationTypeId>7</operationTypeId> <deviceNameId>2</deviceNameId> <productTypeId>1</productTypeId> <density>52.5</density> <tankLevel>12.5</tankLevel> <temperature>12.79</temperature> <volume>12.5</volume> <mass>0.0</mass> </events> <events> <id>407008</id> <datetime>2023-06-07T00:00:00.000+06:00</datetime> <deviceTypeId>1</deviceTypeId> <operationTypeId>7</operationTypeId> <deviceNameId>2</deviceNameId> <productTypeId>1</productTypeId> <density>52.5</density> <tankLevel>12.5</tankLevel> <temperature>12.79</temperature> <volume>12.5</volume> <mass>0.0</mass> </events> <events> <id>407009</id> <datetime>2023-06-07T00:00:00.000+06:00</datetime> <deviceTypeId>1</deviceTypeId> <operationTypeId>7</operationTypeId> <deviceNameId>2</deviceNameId> <productTypeId>1</productTypeId> <density>52.5</density> <tankLevel>12.5</tankLevel> <temperature>12.79</temperature> <volume>12.5</volume> <mass>0.0</mass> </events> </data> </requestData> </request> </ns2:SendMessage>", "id-31fc1a20-7490-4b66-948d-5b61483dd662")
	 if err != nil {
		 errMsg2 := fmt.Sprintf("Error on signWSSE by KALKAN (Ошибка при подписании XML): %s", err)
		 return errMsg2
	 }

	 fmt.Println("--------------------- Подписанный XML в формате WSSE ------------- \n", signedXML)
	 save_to_file(signedXML)
	
	 return "OK"
}

func save_to_file(data string) {

	// Read the environment variable
	outputDir := os.Getenv("DIR")

	// If the environment variable is not set, use a default value
	if outputDir == "" {
		outputDir = "."
	}

	if DATE_STR_FILE == "" {
		log.Fatal("OUTPUT_DIR environment variable is not set")
	}

	// Generate a timestamp
	//timestamp := time.Now().Format("20060102_150405")

	// Construct the filename with the timestamp
	filename := fmt.Sprintf("%s/output_%s.xml", outputDir, DATE_STR_FILE)

	// Open the file for writing
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write data to the file
	_, err = file.WriteString(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data saved to file:", filename)
}

//----------------------------- // ------------------------------------------ //
//----------------------------- // ------------------------------------------ //

func prepareXmlToSign() string {

	currentTime := time.Now()

	getDataAmangeldi();
	getDataAirakty();
	getDataZharkum();

	eventsStrAmangeldi := getEventsXmlAmangeldi()
	eventsStrAirakty := getEventsXmlAirakty()
	eventsStrZharkum := getEventsXmlZharkum()

	if !updatedAmangeldi {
		errMsg := fmt.Sprintf("ERROR : need to check for update data for Amangeldi : %s", pathToFileAmangeldi)
		logging(errMsg)
		return errMsg
	}
	if !updatedAirakty {
		errMsg1 := fmt.Sprintf("ERROR : need to check for update data for Airakty : %s", pathToFileAirakty)
		logging(errMsg1)
		return errMsg1
	}
	if !updatedJarkum {
		errMsg2 := fmt.Sprintf("ERROR : need to check for update data for Zharkum : %s", pathToFileZharkum)
		logging(errMsg2)
		return errMsg2
	}


	xmlString := "<ns2:SendMessage xmlns:ns2=\"http://bip.bee.kz/SyncChannel/v10/Types\"> <request> <requestInfo> <messageId>"
	xmlString += generateRandomString()
	xmlString += "</messageId> <serviceId>"
	xmlString += serviceId
	xmlString += "</serviceId> <messageDate>"
	xmlString += currentTime.Format(timeFormatForEvent)
	xmlString += "</messageDate> <sender> <senderId>"
	xmlString += senderId
	xmlString += "</senderId> <password>"
	xmlString += senderPassword
	xmlString += "</password> </sender> </requestInfo> <requestData> <data xmlns:cs=\"http://message.persistence.interactive.nat\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:type=\"cs:Request\">"
	xmlString += eventsStrAmangeldi + eventsStrAirakty + eventsStrZharkum
	xmlString += " </data> </requestData> </request> </ns2:SendMessage>"

	updatedAmangeldi = false
	updatedAirakty = false
	updatedJarkum = false
	return xmlString
}

//----------------------------- // ------------------------------------------ //

func getEventTime() (string, error) {

	DATE_STR_FILE := os.Getenv("DATE")
	layout := "02.01.2006"
	date, err := time.Parse(layout,DATE_STR_FILE)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "", err
	}
	currentTime := time.Date(date.Year(), date.Month(), date.Day(), 13, 0, 0, 0, date.Location())
	eventTimeString:= currentTime.Format(timeFormatForEvent)
	fmt.Println("Event time:", eventTimeString)
	return eventTimeString, nil
}


func getDataAmangeldi() {
	// here we will update data of rashod407001, rezer407004, rezer407005
	fileName := pathToFileAmangeldi + generateFileName()
	// Read the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		logging(fmt.Sprintf("Error opening file: %s, err: %s", fileName, err))
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		logging(fmt.Sprintf("Error reading CSV file: %s, err: %s", fileName, err))
		return
	}

//	DATE_STR_FILE := os.Getenv("DATE")
//	layout := "2006.01.02"
//	date, err := time.Parse(layout,DATE_STR_FILE)
//	if err != nil {
//		fmt.Println("Error parsing date:", err)
//		return
//	}
//	currentTime := time.Date(date.Year(), date.Month(), date.Day(), 13, 0, 0, 0, date.Location())
//	localTimeString := currentTime.Format(timeFormatForEvent)

	// get time of event
	eventTime, err := getEventTime()

	rashod407001.DateTime = eventTime
	rezer407004.DateTime = eventTime
	rezer407005.DateTime = eventTime

	for _, row := range records {
		keyName := row[0]
		val := row[1]
		logging("keyName :" + keyName + ", val:" + val)
		if keyName == "dev1_density" {
			rashod407001.Density = val
		}
		if keyName == "dev1_volume" {
			rashod407001.Volume = val
		}
		if keyName == "dev1_temperature" {
			rashod407001.Temperature = val
		}
		if keyName == "dev1_massflowbegin" {
			rashod407001.MassFlowBegin = val
		}
		if keyName == "dev1_massflowend" {
			rashod407001.MassFlowEnd = val
		}
		if keyName == "dev1_mass" {
			rashod407001.Mass = val
		}

		if (keyName == "dev2_density") {
			rezer407004.Density = val
		}
		if (keyName == "dev2_volume") {
			rezer407004.Volume = val
		}
		if (keyName == "dev2_temperature") {
			rezer407004.Temperature = val
		}
		if (keyName == "dev2_tankLevel") {
			rezer407004.TankLevel = val
		}
		if (keyName == "dev2_mass") {
			rezer407004.Mass = val
		}

		if (keyName == "dev3_density") {
			rezer407005.Density = val
		}
		if (keyName == "dev3_volume") {
			rezer407005.Volume = val
		}
		if (keyName == "dev3_temperature") {
			rezer407005.Temperature = val
		}
		if (keyName == "dev3_tankLevel") {
			rezer407005.TankLevel = val
		}
		if (keyName == "dev3_mass") {
			rezer407005.Mass = val
		}
	}
	updatedAmangeldi = true
}

func getEventsXmlAmangeldi() string{

	xmlString := "<events> " +
		"<id>" + rashod407001.ID + "</id> " +
		"<datetime>" + rashod407001.DateTime + "</datetime> " +
		"<deviceTypeId>" + rashod407001.DeviceTypeID + "</deviceTypeId> " +
		"<operationTypeId>" + rashod407001.OperationTypeID + "</operationTypeId> " +
		"<deviceNameId>" + rashod407001.DeviceNameID + "</deviceNameId> " +
		"<productTypeId>" + rashod407001.ProductTypeID + "</productTypeId> " +
		"<temperature>" + rashod407001.Temperature + "</temperature> " +
		"<density>" + rashod407001.Density + "</density> " +
		"<volume>" + rashod407001.Volume + "</volume> " +
		"<pipelineId>" + rashod407001.PipelineID + "</pipelineId> " +
		"<massflowbegin>" + rashod407001.MassFlowBegin + "</massflowbegin> " +
		"<massflowend>" + rashod407001.MassFlowEnd + "</massflowend> " +
		"<mass>" + rashod407001.Mass + "</mass> " +
		"</events>"

	xmlString += "<events> " +
		"<id>" + rezer407004.ID + "</id> " +
		"<datetime>" + rezer407004.DateTime + "</datetime> " +
		"<deviceTypeId>" + rezer407004.DeviceTypeID + "</deviceTypeId> " +
		"<operationTypeId>" + rezer407004.OperationTypeID + "</operationTypeId> " +
		"<deviceNameId>" + rezer407004.DeviceNameID + "</deviceNameId> " +
		"<productTypeId>" + rezer407004.ProductTypeID + "</productTypeId> " +
		"<temperature>" + rezer407004.Temperature + "</temperature> " +
		"<density>" + rezer407004.Density + "</density> " +
		"<volume>" + rezer407004.Volume + "</volume> " +
		"<tankLevel>" + rezer407004.TankLevel + "</tankLevel> " +
		"<mass>" + rezer407004.Mass + "</mass> " +
		"</events>"

	xmlString += "<events> " +
		"<id>" + rezer407005.ID + "</id> " +
		"<datetime>" + rezer407005.DateTime + "</datetime> " +
		"<deviceTypeId>" + rezer407005.DeviceTypeID + "</deviceTypeId> " +
		"<operationTypeId>" + rezer407005.OperationTypeID + "</operationTypeId> " +
		"<deviceNameId>" + rezer407005.DeviceNameID + "</deviceNameId> " +
		"<productTypeId>" + rezer407005.ProductTypeID + "</productTypeId> " +
		"<temperature>" + rezer407005.Temperature + "</temperature> " +
		"<density>" + rezer407005.Density + "</density> " +
		"<volume>" + rezer407005.Volume + "</volume> " +
		"<tankLevel>" + rezer407005.TankLevel + "</tankLevel> " +
		"<mass>" + rezer407005.Mass + "</mass> " +
		"</events>"

	return xmlString
}

//----------------------------- // ------------------------------------------ //

func getDataAirakty() {
	// here we will update data of rashod407001, rezer407004, rezer407005
	fileName := pathToFileAirakty + generateFileName()
	// Read the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		logging(fmt.Sprintf("Error opening file: %s, err: %s", fileName, err))
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		logging(fmt.Sprintf("Error reading CSV file: %s, err: %s", fileName, err))
		return
	}

	//currentTime := time.Now()
	//localTimeString := currentTime.Format(timeFormatForEvent)

	// get time of event
	eventTime, err := getEventTime()

	rashod407002.DateTime = eventTime 
	rezer407006.DateTime = eventTime 
	rezer407007.DateTime = eventTime

	for _, row := range records {
		keyName := row[0]
		val := row[1]
		if keyName == "dev1_density" {
			rashod407002.Density = val
		}
		if keyName == "dev1_volume" {
			rashod407002.Volume = val
		}
		if keyName == "dev1_temperature" {
			rashod407002.Temperature = val
		}
		if keyName == "dev1_massflowbegin" {
			rashod407002.MassFlowBegin = val
		}
		if keyName == "dev1_massflowend" {
			rashod407002.MassFlowEnd = val
		}
		if keyName == "dev1_mass" {
			rashod407002.Mass = val
		}

		if (keyName == "dev2_density") {
			rezer407006.Density = val
		}
		if (keyName == "dev2_volume") {
			rezer407006.Volume = val
		}
		if (keyName == "dev2_temperature") {
			rezer407006.Temperature = val
		}
		if (keyName == "dev2_tankLevel") {
			rezer407006.TankLevel = val
		}
		if (keyName == "dev2_mass") {
			rezer407006.Mass = val
		}

		if (keyName == "dev3_density") {
			rezer407007.Density = val
		}
		if (keyName == "dev3_volume") {
			rezer407007.Volume = val
		}
		if (keyName == "dev3_temperature") {
			rezer407007.Temperature = val
		}
		if (keyName == "dev3_tankLevel") {
			rezer407007.TankLevel = val
		}
		if (keyName == "dev3_mass") {
			rezer407007.Mass = val
		}
	}
	updatedAirakty = true
}

func getEventsXmlAirakty() string{

	xmlString := "<events> " +
		"<id>" + rashod407002.ID + "</id> " +
		"<datetime>" + rashod407002.DateTime + "</datetime> " +
		"<deviceTypeId>" + rashod407002.DeviceTypeID + "</deviceTypeId> " +
		"<operationTypeId>" + rashod407002.OperationTypeID + "</operationTypeId> " +
		"<deviceNameId>" + rashod407002.DeviceNameID + "</deviceNameId> " +
		"<productTypeId>" + rashod407002.ProductTypeID + "</productTypeId> " +
		"<temperature>" + rashod407002.Temperature + "</temperature> " +
		"<density>" + rashod407002.Density + "</density> " +
		"<volume>" + rashod407002.Volume + "</volume> " +
		"<pipelineId>" + rashod407002.PipelineID + "</pipelineId> " +
		"<massflowbegin>" + rashod407002.MassFlowBegin + "</massflowbegin> " +
		"<massflowend>" + rashod407002.MassFlowEnd + "</massflowend> " +
		"<mass>" + rashod407002.Mass + "</mass> " +
		"</events>"

	xmlString += "<events> " +
		"<id>" + rezer407006.ID + "</id> " +
		"<datetime>" + rezer407006.DateTime + "</datetime> " +
		"<deviceTypeId>" + rezer407006.DeviceTypeID + "</deviceTypeId> " +
		"<operationTypeId>" + rezer407006.OperationTypeID + "</operationTypeId> " +
		"<deviceNameId>" + rezer407006.DeviceNameID + "</deviceNameId> " +
		"<productTypeId>" + rezer407006.ProductTypeID + "</productTypeId> " +
		"<temperature>" + rezer407006.Temperature + "</temperature> " +
		"<density>" + rezer407006.Density + "</density> " +
		"<volume>" + rezer407006.Volume + "</volume> " +
		"<tankLevel>" + rezer407006.TankLevel + "</tankLevel> " +
		"<mass>" + rezer407006.Mass + "</mass> " +
		"</events>"

	xmlString += "<events> " +
		"<id>" + rezer407007.ID + "</id> " +
		"<datetime>" + rezer407007.DateTime + "</datetime> " +
		"<deviceTypeId>" + rezer407007.DeviceTypeID + "</deviceTypeId> " +
		"<operationTypeId>" + rezer407007.OperationTypeID + "</operationTypeId> " +
		"<deviceNameId>" + rezer407007.DeviceNameID + "</deviceNameId> " +
		"<productTypeId>" + rezer407007.ProductTypeID + "</productTypeId> " +
		"<temperature>" + rezer407007.Temperature + "</temperature> " +
		"<density>" + rezer407007.Density + "</density> " +
		"<volume>" + rezer407007.Volume + "</volume> " +
		"<tankLevel>" + rezer407007.TankLevel + "</tankLevel> " +
		"<mass>" + rezer407007.Mass + "</mass> " +
		"</events>"

	return xmlString
}

//----------------------------- // ------------------------------------------ //

func getDataZharkum() {
	// here we will update data of rashod407001, rezer407004, rezer407005
	fileName := pathToFileZharkum + generateFileName()
	// Read the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		logging(fmt.Sprintf("Error opening file: %s, err: %s", fileName, err))
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		logging(fmt.Sprintf("Error reading CSV file: %s, err: %s", fileName, err))
		return
	}

	//currentTime := time.Now()
	//localTimeString := currentTime.Format(timeFormatForEvent)
	//rashod407003.DateTime = localTimeString
	//rezer407008.DateTime = localTimeString
	//rezer407009.DateTime = localTimeString

	// get time of event
	eventTime, err := getEventTime()

	rashod407003.DateTime = eventTime
	rezer407008.DateTime = eventTime
	rezer407009.DateTime = eventTime

	for _, row := range records {
		keyName := row[0]
		val := row[1]
		if keyName == "dev1_density" {
			rashod407003.Density = val
		}
		if keyName == "dev1_volume" {
			rashod407003.Volume = val
		}
		if keyName == "dev1_temperature" {
			rashod407003.Temperature = val
		}
		if keyName == "dev1_massflowbegin" {
			rashod407003.MassFlowBegin = val
		}
		if keyName == "dev1_massflowend" {
			rashod407003.MassFlowEnd = val
		}
		if keyName == "dev1_mass" {
			rashod407003.Mass = val
		}

		if (keyName == "dev2_density") {
			rezer407008.Density = val
		}
		if (keyName == "dev2_volume") {
			rezer407008.Volume = val
		}
		if (keyName == "dev2_temperature") {
			rezer407008.Temperature = val
		}
		if (keyName == "dev2_tankLevel") {
			rezer407008.TankLevel = val
		}
		if (keyName == "dev2_mass") {
			rezer407008.Mass = val
		}

		if (keyName == "dev3_density") {
			rezer407009.Density = val
		}
		if (keyName == "dev3_volume") {
			rezer407009.Volume = val
		}
		if (keyName == "dev3_temperature") {
			rezer407009.Temperature = val
		}
		if (keyName == "dev3_tankLevel") {
			rezer407009.TankLevel = val
		}
		if (keyName == "dev3_mass") {
			rezer407009.Mass = val
		}
	}
	updatedJarkum = true
}

func getEventsXmlZharkum() string{

	xmlString := "<events> " +
		"<id>" + rashod407003.ID + "</id> " +
		"<datetime>" + rashod407003.DateTime + "</datetime> " +
		"<deviceTypeId>" + rashod407003.DeviceTypeID + "</deviceTypeId> " +
		"<operationTypeId>" + rashod407003.OperationTypeID + "</operationTypeId> " +
		"<deviceNameId>" + rashod407003.DeviceNameID + "</deviceNameId> " +
		"<productTypeId>" + rashod407003.ProductTypeID + "</productTypeId> " +
		"<temperature>" + rashod407003.Temperature + "</temperature> " +
		"<density>" + rashod407003.Density + "</density> " +
		"<volume>" + rashod407003.Volume + "</volume> " +
		"<pipelineId>" + rashod407003.PipelineID + "</pipelineId> " +
		"<massflowbegin>" + rashod407003.MassFlowBegin + "</massflowbegin> " +
		"<massflowend>" + rashod407003.MassFlowEnd + "</massflowend> " +
		"<mass>" + rashod407003.Mass + "</mass> " +
		"</events>"

	xmlString += "<events> " +
		"<id>" + rezer407008.ID + "</id> " +
		"<datetime>" + rezer407008.DateTime + "</datetime> " +
		"<deviceTypeId>" + rezer407008.DeviceTypeID + "</deviceTypeId> " +
		"<operationTypeId>" + rezer407008.OperationTypeID + "</operationTypeId> " +
		"<deviceNameId>" + rezer407008.DeviceNameID + "</deviceNameId> " +
		"<productTypeId>" + rezer407008.ProductTypeID + "</productTypeId> " +
		"<temperature>" + rezer407008.Temperature + "</temperature> " +
		"<density>" + rezer407008.Density + "</density> " +
		"<volume>" + rezer407008.Volume + "</volume> " +
		"<tankLevel>" + rezer407008.TankLevel + "</tankLevel> " +
		"<mass>" + rezer407008.Mass + "</mass> " +
		"</events>"

	xmlString += "<events> " +
		"<id>" + rezer407009.ID + "</id> " +
		"<datetime>" + rezer407009.DateTime + "</datetime> " +
		"<deviceTypeId>" + rezer407009.DeviceTypeID + "</deviceTypeId> " +
		"<operationTypeId>" + rezer407009.OperationTypeID + "</operationTypeId> " +
		"<deviceNameId>" + rezer407009.DeviceNameID + "</deviceNameId> " +
		"<productTypeId>" + rezer407009.ProductTypeID + "</productTypeId> " +
		"<temperature>" + rezer407009.Temperature + "</temperature> " +
		"<density>" + rezer407009.Density + "</density> " +
		"<volume>" + rezer407009.Volume + "</volume> " +
		"<tankLevel>" + rezer407009.TankLevel + "</tankLevel> " +
		"<mass>" + rezer407009.Mass + "</mass> " +
		"</events>"

	return xmlString
}

//----------------------------- // ------------------------------------------ //

func generateRandomString() string {
	uuid := uuid.New()
	return uuid.String()
}

func generateFileName() string {
	// Get the current date and time
	//layout := "2006-01-02"
	// Format the date as "YYYY-M-D"
	//date, _ := time.Parse(layout, DATE_STR_FILE)

	// Generate the file name
	fileName := fmt.Sprintf("log_%s.csv", DATE_STR_FILE)

	return fileName
}

func sendMail(to string, subject string, body string) {
	// Set up authentication information.

	// Compose the email message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", mailSender, to, subject, body)

	// SMTP authentication setup
	auth := smtp.PlainAuth("", mailSender, appPassword, smtpHost)

	// Send the email using mail.ru SMTP
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, mailSender, []string{to}, []byte(message))
	if err != nil {
		logging(fmt.Sprintf("Error sending email: %s",err))
		return
	}

	logging("Email sent successfully.")
}

func logging(text string) {
	// Get the current date and time
	now := time.Now()

	// Format the date as "YYYY-M-D"
	date := now.Format("2006-1-2")

	// Format the time as "HH:MM:SS"
	time := now.Format("15:04:05")

	// Create the log text
	logText := fmt.Sprintf("%s %s: %s \n", date, time, text)

	// check for file exist if not create one
	checkFileExistIfNotCreate(logFileName)

	// Append the log text to the log file
	file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	// Append text to the file
	_, err = io.WriteString(file, logText)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}


func checkFileExistIfNotCreate(filePath string) {
	// check for file exist if not create one
	// Check if the file exists
	_, err := os.Stat(filePath)
	if err != nil {
		// File does not exist, create it
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		fmt.Println("File created:", filePath)
	}
}