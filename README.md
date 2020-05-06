# mapjitsu

[![GoDoc](https://godoc.org/github.com/8legd/mapjitsu?status.svg)](https://godoc.org/github.com/8legd/mapjitsu)

`mapjitsu` is a Go library for wrestling data 

## Overview

`mapjitsu` provides a way of mapping data with Go providing builtin types and functions for common behaviour

Documentation can be found at [godoc.org](http://godoc.org/github.com/8legd/mapjitsu) but as with many abstract concepts it is perhaps best explained by way of example

## Getting Started

An understanding of `mapjitsu` starts with three standard [ETL](https://en.wikipedia.org/wiki/Extract,_transform,_load) terms

### Sources

We start with a source for each item of data. A simple interface is provided to represent this 

```go

type Source interface {
    Value() (interface{}, error)
}

```

You can implement this interface yourself or use one of the builtin [sources](http://godoc.org/github.com/8legd/mapjitsu/sources)

### Transforms

When wrestling data a simple 1 to 1 mapping is often not sufficient and some form of transformation is required

This can be carried out through a series of functions referred to as a Pipeline

A type is provided to represent this as a slice of functions

```go

type Pipeline []func(interface{}) (interface{}, error)

```

### Targets

Finally we have a target which is the destination for the wrestled data item. A simple interface is provided to represent this 

```go

type Target interface {
    SetValue(interface{}) error
}

```

Again you can implement this interface yourself or use one of the builtin [targets](http://godoc.org/github.com/8legd/mapjitsu/targets)

### Putting the Sources, Transforms & Targets together to wrestle some data

The [tests](tests) provide examples

Here is a test of JSON input and output using [MXJ](http://godoc.org/github.com/clbanning/mxj)

```go

func TestMXJ(t *testing.T) {

	// start by unmarshaling some JSON to an MXJ Map
	input, err := mxj.NewMapJson([]byte(`{
		"user": {
			"first_name": "Tim",
			"last_name": "Test",
			"dob": null
		}
	}`))
	if err != nil {
		t.Fatalf("failed to unmarshal input %v", err)
	}

	// we also create an MXJ Map for the target
	output := mxj.Map{
		"Customer": make(map[string]interface{}),
	}

	// optional data items will need an error handler as MXJ will return a PathNotExistError
	// here is an example error handler returning a default value
	// (see user.title mapping below for example usage)
	onNotExistReturnDefault := func(defaultvalue interface{}) func(path string, v interface{}, err error) (interface{}, error) {
		return func(path string, v interface{}, err error) (interface{}, error) {
			if err == mxj.PathNotExistError {
				return defaultvalue, nil
			}
			return v, err
		}
	}

	// transformations can be added through pipeline functions
	// this simple example converts values to strings
	// NOTE: a nil value is converted to an empty string
	toString := func(v interface{}) (interface{}, error) { // this transformation converts nil values to an empty string
		if v == nil {
			return "", nil // return nil as empty string
		}
		// otherwise use default fmt
		return fmt.Sprintf("%v", v), nil
	}

	// next we define our mappings
	definition := mapjitsu.Definition{
		Mappings: []mapjitsu.Mapping{
			{
				// MXJ paths are used here, see https://godoc.org/github.com/clbanning/mxj#Map.ValueForPath
				Source: mxjData.Source{Map: input, Path: "user.first_name"},
				Target: mxjData.Target{Map: output, Path: "Customer.FirstName"},
			},
			{
				Source: mxjData.Source{Map: input, Path: "user.last_name"},
				Target: mxjData.Target{Map: output, Path: "Customer.LastName"},
			},
			{
				Source: mxjData.Source{Map: input, Path: "user.title", OnError: onNotExistReturnDefault("")},
				Target: mxjData.Target{Map: output, Path: "Customer.Title"},
			},
			{
				Source:    mxjData.Source{Map: input, Path: "user.dob"},
				Transform: mapjitsu.Pipeline{toString},
				Target:    mxjData.Target{Map: output, Path: "Customer.DOB"},
			},
			{
				// here a function is used as the Source to combine two text fields
				Source: mapjitsu.SourceFunc(func() (interface{}, error) {
					result := input.ValueOrEmptyForPathString("user.first_name")
					if s := input.ValueOrEmptyForPathString("user.last_name"); s != "" {
						if result != "" {
							result = result + " "
						}
						result = result + s
					}
					if result == "" {
						// optionally an error can be returned here e.g. for required data items
						return result, errors.New("could not calculate Customer.FullName missing either a user.first_name or user.last_name")
					}
					return result, nil
				}),
				Target: mxjData.Target{Map: output, Path: "Customer.FullName"},
			},
		},
	}

	// once our mappings are defined we can apply them
	err = definition.Apply()
	if err != nil {
		t.Fatalf("failed to apply mappings %v", err)
	}

	var json []byte
	json, err = output.JsonIndent("", "\t")
	if err != nil {
		t.Fatalf("failed to marshal output %v", err)
	}

	assert := func(expected string, actual []byte) {
		jsonString := "\n" + string(json)
		if jsonString != expected {
			t.Errorf("resulting json string \n%s\n does not match expected \n%s\n", jsonString, expected)
			return
		}
		t.Logf("%s", jsonString)
	}

	expected := `
{
	"Customer": {
		"DOB": "",
		"FirstName": "Tim",
		"FullName": "Tim Test",
		"LastName": "Test",
		"Title": ""
	}
}`

	assert(expected, json)

}

```

and here using CSV input and output instead


```go

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
					// CSV csvData and csvData can use a column number for the mapping
					Source: csvData.Source{Record: inputRecord, ColumnNumber: 1},  // first_name
					Target: csvData.Target{Record: outputRecord, ColumnNumber: 2}, // Customer FirstName
				},
				{
					// CSV csvData and csvData can also use column names for the mapping if the header is provided
					Source: csvData.Source{Record: inputRecord, ColumnName: "last_name", Header: inputHeader},
					Target: csvData.Target{Record: outputRecord, ColumnName: "Customer LastName", Header: outputHeader},
				},
				{
					Source: csvData.Source{Record: inputRecord, ColumnName: "dob", Header: inputHeader},
					Target: csvData.Target{Record: outputRecord, ColumnName: "Customer DOB", Header: outputHeader},
				},
				{
					// here a function is used as the Source to combine two text fields
					Source: mapjitsu.SourceFunc(func() (interface{}, error) {
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
					}),
					Target: csvData.Target{Record: outputRecord, ColumnName: "Customer FullName", Header: outputHeader},
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

```

## Contributing

Tests

`GO111MODULE=on go test -v github.com/8legd/mapjitsu/tests`

Correctness

`GO111MODULE=on go vet ./...`

Coding style

`golint ./...`