func TestReadCSVFile(t *testing.T) {
	// Create a temporary CSV file for testing
	tempFile, err := ioutil.TempFile("", "test.csv")
	if err != nil {
		t.Fatalf("failed to create temporary CSV file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the temporary CSV file
	testData := []string{"header1,header2,header3", "value1,value2,value3"}
	if _, err := tempFile.WriteString(strings.Join(testData, "\n")); err != nil {
		t.Fatalf("failed to write test data to temporary CSV file: %v", err)
	}

	// Call the function under test
	headers, records, err := readCSVFile(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to read CSV file: %v", err)
	}

	// Verify the results
	expectedHeaders := []string{"header1", "header2", "header3"}
	if !reflect.DeepEqual(headers, expectedHeaders) {
		t.Errorf("unexpected headers, got: %v, want: %v", headers, expectedHeaders)
	}

	expectedRecords := [][]string{{"value1", "value2", "value3"}}
	if !reflect.DeepEqual(records, expectedRecords) {
		t.Errorf("unexpected records, got: %v, want: %v", records, expectedRecords)
	}
}