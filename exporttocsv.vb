Sub ExportToXML(strObjectName, strFilePath)               

	Dim strFilename, strLine
	Dim i, k
	Dim fso, objFile, objTag
	Dim dTime, sDateTime, sYear, sMonth, sDay, sHour, sMin, sSec
	
	dTime = Now()
	
	sYear = DatePart("yyyy", dTime)
	sMonth = (DatePart("m", dTime))
	sDay = (DatePart("d", dTime))
	sHour = (DatePart("h", dTime))
	sMin = (DatePart("n", dTime))
	sSec = (DatePart("s", dTime))

    Dim sColumns
    sColumns = "dev1_density," & "dev1_massflowbegin," & "dev1_massflowend," & "dev1_mass," & "dev2_density," & "dev2_volume," & "dev2_temperature," & "dev2_level," & "dev2_mass," & "dev3_density," & "dev3_volume," & "dev3_temperature," & "dev3_level," & "dev3_mass"
	
    Dim dev1_density, dev1_massflowbegin, dev1_massflowend, dev1_mass, dev2_density, dev2_volume, dev2_temperature, dev2_level, dev2_mass, dev3_density, dev3_volume, dev3_temperature, dev3_level, dev3_mass
    
    Set dev1_density = HMIRuntime.Tags(strObjectName & "_dev1.density") 
    Set dev1_massflowbegin= HMIRuntime.Tags(strObjectName & "_dev1_massflowbegin")
    Set dev1_massflowend = HMIRuntime.Tags(strObjectName & "_dev1_massflowend")
    Set dev1_mass = HMIRuntime.Tags(strObjectName & "_dev1_mass")
    Set dev2_density = HMIRuntime.Tags(strObjectName & "_dev2_density")
    Set dev2_volume = HMIRuntime.Tags(strObjectName & "_dev2_volume")
    Set dev2_temperature = HMIRuntime.Tags(strObjectName & "_dev2_temperature")
    Set dev2_level = HMIRuntime.Tags(strObjectName & "_dev2_level")
    Set dev2_mass = HMIRuntime.Tags(strObjectName & "_dev2_mass")
    Set dev3_density = HMIRuntime.Tags(strObjectName & "_dev3_density")
    Set dev3_volume = HMIRuntime.Tags(strObjectName & "_dev3_volume")
    Set dev3_temperature = HMIRuntime.Tags(strObjectName & "_dev3_temperature")
    Set dev3_level = HMIRuntime.Tags(strObjectName & "_dev3_level")
    Set dev3_mass = HMIRuntime.Tags(strObjectName & "_dev3_mass") 

    Dim sValues
    sValues = dev1_density & "," &  dev1_massflowbegin & "," &  dev1_massflowend & "," &  dev1_mass & "," &  dev2_density & "," &  dev2_volume & "," &  dev2_temperature & "," &  dev2_level & "," &  dev2_mass & "," &  dev3_density & "," &  dev3_volume & "," &  dev3_temperature & "," &  dev3_level & "," &  dev3_mass

    Dim arrLines(2)

    arrLines(0) = sColumns
    arrLines(1) = sValues

	'HMIRuntime.Tags("Filename").Write strFilename
	Set fso = CreateObject("Scripting.FileSystemObject")
	Set objFile = fso.CreateTextFile(strFilename,True)

    For k=0 To UBound(arrLines) - 1
        objFile.WriteLine arrLines(k)
		'for control
		HMIRuntime.Trace(strLine & vbCrLf)
    Next

    objFile.Close

	Dim sLastChange
	sLastChange = sYear & "/" & sMonth & "/" & sDay & " " & sHour & ":" & sMin
	HMIRuntime.Tags("Last_change").Write sLastChange	
	
	'create file
	sDateTime = sYear & "." & sMonth & "." & sDay
	
	Dim sToday
	sToday = Date
	strFilename = strFilePath & "raw_" & sToday & ".csv"
	
	sLastChange = dTime
	HMIRuntime.Tags("Last_change").Write sLastChange

    MassFlowDayRestart strObjectName, "dev1"
	
End Sub