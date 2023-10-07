#include "apdefap.h"
void UAWriteToArchive()
{
    UAHCONNECT hConnect = 0;
    UAHARCHIVE hArchive = 0;
    LONG IndexArchive;
    LONG FieldLength;
    LONG FieldType;
    LONG NumberOfFields;
    LONG Index;
    long IntValue;
    char StringField[255];
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
    //******* Connect to Archive via Name ****************************
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
    //******* Get Number of Fields **********************************************
    NumberOfFields = uaArchiveGetFields( hArchive );
    printf("Number of Fields = %u\n", NumberOfFields );
    //******* Read Last Data Set ************************************************
    if (uaArchiveMoveLast( hArchive ) == TRUE )
        printf("Number of Fields = %u\n", NumberOfFields );
    else
    {
        printf("uaArchiveMoveLast Error: %d\n", uaGetLastError() );
        goto finish;
    }
    //******* Write into Data Fields *********************************************
    IntValue = 105;//RGB for darkgray
    strcpy(StringField, "darkgray" );
    GetSystemTime( &SysDate );
    for ( Index = 1; Index < NumberOfFields; Index++ )
    {
        printf("Data of Field %u: \n", Index );
        FieldType = uaArchiveGetFieldType( hArchive, Index );
        switch ( FieldType )
        {
            case UA_FIELDTYPE_INTEGER :
            printf("Field Type = Integer\n");
            if (uaArchiveSetFieldValueLong ( hArchive, Index, IntValue ) == TRUE )
            printf( "Field Value = %u\n", IntValue );
            else
            printf("Error callinguaArchiveSetFieldValueLong: %d\n", uaGetLastError() );
            break;
            case UA_FIELDTYPE_FLOAT :
            printf("Field Type = Float\n");
            if (uaArchiveSetFieldValueFloat ( hArchive, Index, FloatValue ) == TRUE )
            printf("Field Value = %f\n", FloatValue );
            else
            printf("Error callinguaArchiveSetFieldValueFloat: %d\n", uaGetLastError() );
            break;
            case UA_FIELDTYPE_DOUBLE :
            printf("Field Type = Double\n");
            if (uaArchiveSetFieldValueDouble (hArchive, Index, DoubleValue ) == TRUE )
            printf("Field Value = %g\n", DoubleValue );
            else
            printf("Error calling uaArchiveSetFieldValueDouble: %d\n", uaGetLastError() );
            break;
            case UA_FIELDTYPE_STRING :
            printf("Field Type = String\n");
            if (uaArchiveSetFieldValueString ( hArchive, Index, StringField ) == TRUE )
            printf("Field Value = %s\n", StringField );
            else
            printf("Error callinguaArchiveSetFieldValueString: %d\n", uaGetLastError() );
            break;
            case UA_FIELDTYPE_DATETIME :
            printf("Field Type = Date & Time\n");
            if (uaArchiveSetFieldValueDate ( hArchive, Index, &SysDate ) == TRUE )
            printf("%d.%d.%d\n ",SysDate.wDay, SysDate.wMonth, SysDate.wYear );
            else
            printf("Error calling uaArchiveSetFieldValueLong: %d\n", uaGetLastError() );
            break;
            case -1 :
            default:
            printf("Error executing uaArchiveSetFieldType\n");
        }
        FieldLength = uaArchiveGetFieldLength( hArchive, Index );
        if ( FieldLength != -1 )
        printf("Field Length = %u\n", FieldLength );
        else
        printf("Error executing uaArchiveGetFieldLength\n");
    }
    // ******* Update Archive *******************************************
    if (uaArchiveUpdate(hArchive) == FALSE )
    {
        printf("uaArchiveUpdate Error:\n" );
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