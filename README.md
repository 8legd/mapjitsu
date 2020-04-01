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

Some example Pipeline functions are also provided e.g. a simple ToString function

```go

func ToString(v interface{}) (interface{}, error) {
    if v == nil {
        return "", nil // return nil as empty string
    }
    // otherwise use default fmt
    return fmt.Sprintf("%v", v), nil
}

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

Here is an example test with JSON input and output using [MXJ](http://godoc.org/github.com/clbanning/mxj)

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

	// next we define our mappings
	definition := mapjitsu.Definition{
		Mappings: []mapjitsu.Mapping{
			{
				// MXJ paths are used here, see https://godoc.org/github.com/clbanning/mxj#Map.ValueForPath
				Source: sources.MXJ{Map: input, Path: "user.first_name"},
				Target: targets.MXJ{Map: output, Path: "Customer.FirstName"},
			},
			{
				Source: sources.MXJ{Map: input, Path: "user.last_name"},
				Target: targets.MXJ{Map: output, Path: "Customer.LastName"},
			},
			{
				Source: sources.MXJ{Map: input, Path: "user.title", OnError: onNotExistReturnDefault("")},
				Target: targets.MXJ{Map: output, Path: "Customer.Title"},
			},
			{
				Source:    sources.MXJ{Map: input, Path: "user.dob"},
				Transform: mapjitsu.Pipeline{mapjitsu.ToString}, // will convert nil values to an empty string
				Target:    targets.MXJ{Map: output, Path: "Customer.DOB"},
			},
			{
				// here the builtin Calculated source is used to combine two text fields
				Source: sources.Calculated{
					Formula: func() (interface{}, error) {
						result := input.ValueOrEmptyForPathString("user.first_name")
						if s := input.ValueOrEmptyForPathString("user.last_name"); s != "" {
							if result != "" {
								result = result + " "
							}
							result = result + s
						}
						if result == "" {
							// optionally an error can be returned here e.g. for required data items
							return result, errors.New("could not calculate fullName missing either a user.first_name or user.last_name")
						}
						return result, nil
					},
				},
				Target: targets.MXJ{Map: output, Path: "Customer.FullName"},
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

### Further examples

The [tests](tests) provide further examples


## Contributing

Tests

`GO111MODULE=on go test -v github.com/8legd/mapjitsu/tests`