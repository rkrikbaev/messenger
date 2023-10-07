#include "apdefap.h"
void UAReadFromArchive()
{
UAHCONNECT hConnect = 0;
UAHARCHIVE hArchive = 0;
LONG IndexArchive;
LONG FieldLength;
LONG FieldType;
LONG NumberOfFields;
LONG Index;
long IntValue;
float FloatValue;
double DoubleValue;
char ArchivName[255], StringField[255];
SYSTEMTIME SysDate;
//******* Connect to Component User Archives ****************************
if (uaConnect( &hConnect ) == FALSE )
{
printf("uaConnect error: %d\n", uaGetLastError() );
return;
}
if (hConnect == NULL)
{
printf("Handle UAHCONNECT equals NULL\n" );
return;
}
//******* Connect to Archive via Archive Name ****************************
if (uaQueryArchiveByName( hConnect, "color", &hArchive ) == FALSE )
{
printf("uaQueryArchive Error: %d\n", uaGetLastError() );
goto finish;
}
//******* Opens Archive ******************************************************
if ( uaArchiveOpen( hArchive ) == FALSE )
{
printf("uaArchive Open Error\n" );
goto finish;
}
//******* Move to first record set ****************************************************
if (uaArchiveMoveFirst(hArchive) == FALSE )
{
printf("uaArchiveMoveFirst Error = %d\n" , uaGetLastError() );
goto finish;
}
//******* Get Number of Fields **********************************************
NumberOfFields = uaArchiveGetFields( hArchive );
printf("Number of Fields = %u\n", NumberOfFields );
//******* Read and show Data Fields ****************************************
for ( Index = 1; Index < NumberOfFields; Index++ )
{
printf("Data of Field %u: \n", Index );
FieldType = uaArchiveGetFieldType( hArchive, Index );
switch ( FieldType )
{
case UA_FIELDTYPE_INTEGER :
printf("Field Type = Integer\n");
if ( uaArchiveGetFieldValueLong ( hArchive, Index, &IntValue ) == TRUE )
printf( "Field Value = %u\n", IntValue );
else
printf("Error callinguaArchiveGetFieldValueLong: %d\n", uaGetLastError() );
break;
case UA_FIELDTYPE_FLOAT :
printf("Field Type = Float\n");
if (uaArchiveGetFieldValueFloat ( hArchive, Index, &FloatValue ) == TRUE )
printf("Field Value = %f\n", FloatValue );
else
printf("Error callinguaArchiveGetFieldValueFloat: %d\n", uaGetLastError() );
break;
case UA_FIELDTYPE_DOUBLE :
printf("Field Type = Double\n");
if (uaArchiveGetFieldValueDouble (hArchive, Index, &DoubleValue ) == TRUE )
printf("Field Value = %g\n", DoubleValue );
else
printf("Error calling uaArchiveGetFieldValueDouble: %d\n", uaGetLastError() );
break;
case UA_FIELDTYPE_STRING :
printf("Field Type = String\n");
if (uaArchiveGetFieldValueString ( hArchive, Index, StringField, 20 ) == TRUE )
printf("Field Value = %s\n", StringField );
else
printf("Error callinguaArchiveGetFieldValueString: %d\n", uaGetLastError() );
break;
case UA_FIELDTYPE_DATETIME :
printf("Field Type = Date & Time\n");
if (uaArchiveGetFieldValueDate ( hArchive, Index, &SysDate ) == TRUE )
printf("%d.%d.%d\n ",SysDate.wDay, SysDate.wMonth, SysDate.wYear );
else
printf("Error calling uaArchiveGetFieldValueLong: %d\n", uaGetLastError() );
break;
case -1 :
default:
printf("Error executing uaArchiveGetFieldType\n");
}
//******* Read and show Field Length **************************************
FieldLength = uaArchiveGetFieldLength( hArchive, Index );
if ( FieldLength != -1 )
printf("Field Length = %u\n", FieldLength );
else
printf("Error executing uaArchiveGetFieldLength\n");
}
//******* Close all handles and connections ***************************
finish:;
//******* Close Archive *******************************************************
if( NULL != hArchive )
{
if (uaArchiveClose ( hArchive ) == FALSE )
{
printf("error on closing archive\n" );
}
}
//****** Release Connection to Archive *************************************
if( NULL != hArchive )
{
if (uaReleaseArchive ( hArchive ) == FALSE )
{
printf("error on releasing archive\n" );
}
hArchive = 0;
}
//******* Disconnect to Component User Archives *************************
if( NULL != hConnect )
{
if (uaDisconnect ( hConnect ) == FALSE )
{
printf("error on disconnection\n" );
}
hConnect = 0;
}
}