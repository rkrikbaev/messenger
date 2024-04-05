'VBS352 
Dim objTrendControl
Dim objTrendWindow
Dim objTimeAxis
Dim objValueAxis
Dim objTrend

'create reference to TrendControl
Set objTrendControl = ScreenItems("TrendControl")
'---- reference trend ----
'create reference to new windows, time and value axis
'uploading window
Set objTrendWindow = objTrendControl.GetTrendWindowCollection.AddItem("Uploading")
Set objTimeAxis = objTrendControl.GetTimeAxisCollection.AddItem("uploadingTimeAxis")
Set objValueAxis = objTrendControl.GetValueAxisCollection.AddItem("uploadingValueAxis")
'assign time and value axis to the window
objTimeAxis.TrendWindow = objTrendWindow.Name
objTimeAxis.ShowDate = False
objValueAxis.TrendWindow = objTrendWindow.Name
'add new trend and assign properties
Set objTrend = objTrendControl.GetTrendCollection.AddItem("myTrend")
objTrend.Provider = 1
objTrend.TagName = "G_Archive\Trend_1"
objTrend.Color = RGB(255,200,0)
objTrend.Fill = True
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name

'capacity window
Set objTrendWindow = objTrendControl.GetTrendWindowCollection.AddItem("Capacity")
Set objTimeAxis = objTrendControl.GetTimeAxisCollection.AddItem("uploadingTimeAxis")
Set objValueAxis = objTrendControl.GetValueAxisCollection.AddItem("uploadingValueAxis")
'assign time and value axis to the window
objTimeAxis.TrendWindow = objTrendWindow.Name
objTimeAxis.ShowDate = False
objValueAxis.TrendWindow = objTrendWindow.Name
'add trend and assign propertys
Set objTrend = objTrendControl.GetTrendCollection.AddItem("capacityTankATrend")
objTrend.Provider = 1
objTrend.TagName = "G_Archive\Trend_1"
objTrend.Color = RGB(0,0,0)
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name
'add trend and assign propertys
Set objTrend = objTrendControl.GetTrendCollection.AddItem("capacityTankBTrend")
objTrend.Provider = 1
objTrend.TagName = "G_Archive\Trend_1"
objTrend.Color = RGB(0,0,0)
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name


' 'generate values for reference trend
' dtCurrent = CDate("23.11.2006 00:00:00")
' For lIndex = 0 To 360
' vValues(lIndex) = ( Sin(dblCurrent) * 60 ) + 60
' vTimeStamps(lIndex) = dtCurrent
' dblCurrent = dblCurrent + 0.105
' dtCurrent = dtCurrent + CDate ("00:00:01")
' Next
' 'insert data to the reference trend
' objTrend.RemoveData
' objTrend.InsertData vTimeStamps, vValues
' '---- data trend ----
' 'add time and value axis to the existing window
' Set objTimeAxis = objTrendControl.GetTimeAxisCollection.AddItem("myTimeAxis")
' Set objValueAxis = objTrendControl.GetValueAxisCollection.AddItem("myValueAxis")
' 'assign time and value axis to the window
' objTimeAxis.TrendWindow = objTrendWindow.Name
' objTimeAxis.ShowDate = False
' objValueAxis.TrendWindow = objTrendWindow.Name
' 'add new trend and assign properties
' Set objTrend = objTrendControl.GetTrendCollection.AddItem("myTrend")
' objTrend.Provider = 0
' objTrend.Color = RGB(255,200,0)
' objTrend.Fill = True
' objTrend.TrendWindow = objTrendWindow.Name
' objTrend.TimeAxis = objTimeAxis.Name
' objTrend.ValueAxis = objValueAxis.Name
' 'generate values for data trend
' dtCurrent = CDate("23.11.2006 00:00:00")
' For lIndex = 0 To 360
' vValues(lIndex) = ( Sin(dblCurrent) * 60 ) + 60
' vTimeStamps(lIndex) = dtCurrent
' dblCurrent = dblCurrent + 0.106
' dtCurrent = dtCurrent + CDate ("00:00:01")
' Next
' 'insert values to the data trend
' objTrend.RemoveData
' objTrend.InsertData vValues, vTimeStamps




'VBS361
Sub OnClick(ByVal Item) 
Dim objTrendControl
Dim objTrendWindow
Dim objTimeAxis
Dim objValueAxis
Dim objTrend
'create reference to TrendControl
Set objTrendControl = ScreenItems("TrendControl")
'create reference to new window, time and value axis
Set objTrendWindow = objTrendControl.GetTrendWindowCollection.AddItem("myWindow")
Set objTimeAxis = objTrendControl.GetTimeAxisCollection.AddItem("myTimeAxis")
Set objValueAxis = objTrendControl.GetValueAxisCollection.AddItem("myValueAxis")
'assign time and value axis to the window
objTimeAxis.TrendWindow = objTrendWindow.Name
objValueAxis.TrendWindow = objTrendWindow.Name
' assign properties to trendwindow
objTrendWindow.HorizontalGrid = False
' add new trend and assign properties
Set objTrend = objTrendControl.GetTrendCollection.AddItem("myTrend1")
objTrend.Provider = 1
objTrend.TagName = "G_Archive\Trend_1"
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name
objTrend.Color = RGB(255,0,0)
objTrend.PointStyle = 0
'add new trend and assign properties
Set objTrend = objTrendControl.GetTrendCollection.AddItem("myTrend2")
objTrend.Provider = 1
objTrend.TagName = "G_Archive\Trend_2"
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name
objTrend.Color = RGB(0,255,0)
objTrend.LineWidth = 3
'add new trend and assign properties
Set objTrend = objTrendControl.GetTrendCollection.AddItem("myTrend3")
objTrend.Provider = 1
objTrend.TagName = "G_Archive\Trend_3"
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name
objTrend.Color = RGB(0,0,255)
objTrend.LineType = 2
End Sub


‘VBS362
Sub OnClick(ByVal Item) 
Dim objControl
Dim objTimeColumn
Dim objValueColumn
Set objControl = ScreenItems("Control1")
' Control wide specification
objControl.ColumnResize = False
objControl.TimeBase = 1
objControl.TimeColumnTimeFormat = "HH:mm:ss tt"
objControl.TimeColumnLength = 20
' properties for Time column
Set objTimeColumn = objControl.GetTimeColumn("Time column 1")
objTimeColumn.DateFormat = "dd/MM/yy"
' properties for a new 4th value column with connection to archive tag "Trend_4"
Set objValueColumn = objControl.GetValueColumnCollection.AddItem("Trend 4") 
objValueColumn.Caption = "Trend 4"
objValueColumn.Length = 10
objValueColumn.Align = 1
objValueColumn.Provider = 1
objValueColumn.TagName = "G_Archive\Trend_4"
objValueColumn.TimeColumn = "Time column 1"
End Sub


'VBS346
Dim ctrl
Dim objValCol1
Dim objValCol2
Dim coll
Dim valcol
Set ctrl = ScreenItems("TableControl")
Set objValCol1 = ctrl.GetValueColumnCollection.AddItem("ValueColumn1")
Set objValCol2 = ctrl.GetValueColumnCollection.AddItem("ValueColumn2")
objValCol1.Caption = "Value Archive"
objValCol1.Provider = 1
objValCol1.TagName = "ProcessValueArchive\arch1"
objValCol1.TimeColumn = "TimeColumn1"
objValCol2.Caption = "Value Tag"
objValCol2.Provider = 2
objValCol2.TagName = "tagxx" 
objValCol2.TimeColumn = "TimeColumn2"
Set coll = ctrl.GetValueColumnCollection
For Each valcol In coll
valcol.Align = 2
valcol.Length = 10
valcol.AutoPrecisions = TRUE
Next






' ----------------------------------

'VBS361
Sub OnClick(ByVal Item) 
Dim objTrendControl
Dim objTrendWindow
Dim objTimeAxis
Dim objValueAxis
Dim objTrend
'create reference to TrendControl
Set objTrendControl = ScreenItems("TrendControl")
'find all window objects
Dim collWindowObjects
Set collWindowObjects = objTrendControl.GetTrendWindowCollection
For Each window In collWindowObjects
window.Align = 2
window.Length = 10
window.AutoPrecisions = TRUE
Next
Set objTrendWindow = objTrendControl.GetTrendWindowCollection.AddItem("unloadWindow")
Set objTimeAxis = objTrendControl.GetTimeAxisCollection.AddItem("unloadTimeAxis")
Set objValueAxis = objTrendControl.GetValueAxisCollection.AddItem("unloadValueAxis")
'assign time and value axis to the window
objTimeAxis.TrendWindow = objTrendWindow.Name
objValueAxis.TrendWindow = objTrendWindow.Name
' assign properties to trendwindow
objTrendWindow.HorizontalGrid = True
' add new trend and assign properties
Set objTrend = objTrendControl.GetTrendCollection.AddItem("unloadMassTrend")
objTrend.Provider = 1
objTrend.TagName = "G_Archive\Trend_1"
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name
objTrend.Color = RGB(255,0,0)
objTrend.PointStyle = 0

'create reference to new MASSFLOW window, time and value axis
Set objTrendWindow = objTrendControl.GetTrendWindowCollection.AddItem("reserveWindow")
Set objTimeAxis = objTrendControl.GetTimeAxisCollection.AddItem("reserveTimeAxis")
Set objValueAxis = objTrendControl.GetValueAxisCollection.AddItem("reserveValueAxis")
'add new reserve trend and assign properties
Set objTrend = objTrendControl.GetTrendCollection.AddItem("reserveMassTrend1")
objTrend.Provider = 1
objTrend.TagName = "G_Archive\Trend_2"
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name
objTrend.Color = RGB(0,255,0)
objTrend.LineWidth = 3
'add new trend and assign properties
Set objTrend = objReserveTrendControl.GetTrendCollection.AddItem("reserveMassTrend2")
objTrend.Provider = 1
objTrend.TagName = "G_Archive\Trend_3"
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name
objTrend.Color = RGB(0,0,255)
objTrend.LineType = 2
End Sub



' Rustam Krikbayev
' rkrikbaev@gmail.com

Dim objTrendControl
Dim objTrendWindow
Dim objTimeAxis
Dim objValueAxis
Dim objTrend

'create reference to TrendControl
Set objTrendControl = ScreenItems("TrendControl")
'---- reference trend ----
'create reference to new windows, time and value axis
'uploading window
Set objTrendWindow = objTrendControl.GetTrendWindowCollection.AddItem("Uploading")
Set objTimeAxis = objTrendControl.GetTimeAxisCollection.AddItem("uploadingTimeAxis")
Set objValueAxis = objTrendControl.GetValueAxisCollection.AddItem("uploadingValueAxis")
'assign time and value axis to the window
objTimeAxis.TrendWindow = objTrendWindow.Name
objTimeAxis.ShowDate = False
objValueAxis.TrendWindow = objTrendWindow.Name
'add new trend and assign properties
Set objTrend = objTrendControl.GetTrendCollection.AddItem("Uploading")
objTrend.Provider = 1
objTrend.TagName = "calculated\aman_dev1.mass"
objTrend.Color = RGB(255,200,0)
objTrend.Fill = True
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name

'capacity window
Set objTrendWindow = objTrendControl.GetTrendWindowCollection.AddItem("Capacity")
Set objTimeAxis = objTrendControl.GetTimeAxisCollection.AddItem("capaityTimeAxis")
Set objValueAxis = objTrendControl.GetValueAxisCollection.AddItem("capacityValueAxis")
'assign time and value axis to the window
objTimeAxis.TrendWindow = objTrendWindow.Name
objTimeAxis.ShowDate = False
objValueAxis.TrendWindow = objTrendWindow.Name
'add trend and assign propertys
Set objTrend = objTrendControl.GetTrendCollection.AddItem("capacityTankATrend")
objTrend.Provider = 1
objTrend.TagName = "calculated\aman_dev2.mass"
objTrend.Color = RGB(0,0,0)
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name
'add trend and assign propertys
Set objTrend = objTrendControl.GetTrendCollection.AddItem("capacityTankBTrend")
objTrend.Provider = 1
objTrend.TagName = "calculated\aman_dev3.mass"
objTrend.Color = RGB(0,0,0)
objTrend.TrendWindow = objTrendWindow.Name
objTrend.TimeAxis = objTimeAxis.Name
objTrend.ValueAxis = objValueAxis.Name






'---------------TableControl----------------
'
'
'VBS351 

Dim objTableControl

Dim objTimeColumn

Dim objValueColumn

Dim objTrend


'create reference to TableControl and enable BackColor

Set objTableControl = ScreenItems("Control")

objTableControl.UseColumnBackColor = True


'create reference to new TimeColumn and assign column length

Set objTimeColumn = objTableControl.GetTimeColumnCollection.AddItem("TimeAxis")

objTimeColumn.Length = 20


'add new ValueColumn and assign propertys
Dim arrObjectsList
arrObjectsList = Array(Array("MassTotalUploadBegin", "calculated\aman_dev1.mass"))

For i=0 To UBound(arrObjectsList)

	strColumnName = arrObjectsList(i)(0)
	strArchiveName =  arrObjectsList(i)(1)
	
	'for control
	HMIRuntime.Trace("VB-Script: Archive name: " & strArchiveName & vbCrLf)

	CreateTableColumn strColumnName, strArchiveName

Next


'---------------TableControl----------------
' create column procedure
'
'VBS351 

Set objValueColumn = objTableControl.GetValueColumnCollection.AddItem(strColumnName)

objValueColumn.Provider = 1

objValueColumn.TagName = strArchiveName

objValueColumn.BackColor = RGB(255,255,255)

objValueColumn.TimeColumn = objTimeColumn.Name

' 'add new ValueColumn and assign propertys

' Set objValueColumn = objTableControl.GetValueColumnCollection.AddItem("myValueTable2")

' objValueColumn.Provider = 1

' objValueColumn.TagName = "Process value archive\PDL_ZT_2"

' objValueColumn.BackColor = RGB(0,255,255)

' objValueColumn.TimeColumn = objTimeColumn.Name


' 'add new ValueColumn and assign propertys

' Set objValueColumn = objTableControl.GetValueColumnCollection.AddItem("myValueTable3")

' objValueColumn.Provider = 1

' objValueColumn.TagName = "Process value archive\PDL_ZT_3"

' objValueColumn.BackColor = RGB(255,255,0)

' objValueColumn.TimeColumn = objTimeColumn.Name



'---------------TableControl----------------
' create value trends
'
'VBS351 

Dim arrObjectsList

arrObjectsList = Array(10)
arrObjectsList(0) = Array("Отгрузка", "calculated\aman_dev1.mass")
arrObjectsList(1) = Array("Запас в резерв. А", "calculated\aman_dev2.mass")
arrObjectsList(2) = Array("Запас в резерв. В", "calculated\aman_dev3.mass")

createTableControl arrObjectsList