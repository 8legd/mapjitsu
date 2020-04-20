package tests

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/8legd/mapjitsu"
	"github.com/8legd/mapjitsu/sources"
	"github.com/8legd/mapjitsu/targets"
)

// Example test with CSV input and output
func TestCSV(t *testing.T) {

	inputCSV := `first_name, last_name, dob
Tim,Test,
Tina,Test,01/01/2000`

	inputHeader := []string{"first_name", "last_name", "dob"}
	outputHeader := []string{"Customer DOB", "Customer FirstName", "Customer FullName", "Customer LastName", "Customer Title"}

	// initialise output
	var output [][]string
	output = append(output, outputHeader)

	r := csv.NewReader(strings.NewReader(inputCSV))
	inputRecord, err := r.Read() // read header
	if err != nil {
		t.Fatalf("failed to read header %v", err)
	}
	row := 0
	for { // read next input row
		inputRecord, err = r.Read()
		if err == io.EOF {
			break
		}
		row = row + 1 // keep track of the row number
		if err != nil {
			t.Fatalf("failed to read input at row %d", err)
		}

		// initialise output
		outputRecord := []string{"", "", "", "", ""}

		// define row mappings
		definition := mapjitsu.Definition{
			Mappings: []mapjitsu.Mapping{
				{
					// CSV sources and targets can use a column number for the mapping
					Source: sources.CSV{Record: inputRecord, ColumnNumber: 1},  // first_name
					Target: targets.CSV{Record: outputRecord, ColumnNumber: 2}, // Customer FirstName
				},
				{
					// CSV sources and targets can also use column names for the mapping if the header is provided
					Source: sources.CSV{Record: inputRecord, ColumnName: "last_name", Header: inputHeader},
					Target: targets.CSV{Record: outputRecord, ColumnName: "Customer LastName", Header: outputHeader},
				},
				{
					Source: sources.CSV{Record: inputRecord, ColumnName: "dob", Header: inputHeader},
					Target: targets.CSV{Record: outputRecord, ColumnName: "Customer DOB", Header: outputHeader},
				},
				{
					// here the builtin Calculated source is used to combine two text fields
					Source: sources.Calculated{
						Formula: func() (interface{}, error) {
							firstName := inputRecord[0]
							result := firstName
							lastName := inputRecord[1]
							if result != "" && lastName != "" {
								result = result + " "
							}
							result = result + lastName
							if result == "" {
								// optionally an error can be returned here e.g. for required data items
								return result, errors.New("could not calculate Customer FullName missing either a first_name or last_name")
							}
							return result, nil
						},
					},
					Target: targets.CSV{Record: outputRecord, ColumnName: "Customer FullName", Header: outputHeader},
				},
			},
		}

		// once our mappings are defined we can apply them
		err = definition.Apply()
		if err != nil {
			t.Fatalf("failed to apply row mappings at row %d %v", row, err)
		}

		// append output
		output = append(output, outputRecord)

	}

	assert := func(expected string, actual string) {
		if actual != expected {
			t.Errorf("resulting output \n%s does not match expected \n%s", actual, expected)
			return
		}
		t.Logf("%s", actual)
	}

	var outputCSV strings.Builder

	w := csv.NewWriter(&outputCSV)
	w.WriteAll(output)
	if err := w.Error(); err != nil {
		t.Fatalf("failed to write output %v", err)
	}

	expected := `Customer DOB,Customer FirstName,Customer FullName,Customer LastName,Customer Title
,Tim,Tim Test,Test,
01/01/2000,Tina,Tina Test,Test,
`

	assert(expected, outputCSV.String())

}
