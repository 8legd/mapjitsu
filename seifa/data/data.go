// SOCIO-ECONOMIC INDEXES FOR AREAS (SEIFA) 2016
// Copyright Commonwealth of Australia
// See https://www.abs.gov.au/websitedbs/D3310114.nsf/Home/%A9+Copyright?opendocument/
package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// NOTE: Postal Areas (POAs) are an ABS approximation of postcodes
// see https://www.abs.gov.au/ausstats/abs@.nsf/Lookup/by%20Subject/1270.0.55.003~July%202016~Main%20Features~Postal%20Areas%20(POA)~8
func FirstMatchingRecord(table int, poa string) ([]string, error) {

	// 6 SEIFA data tables are available to read from
	var csvdata string
	switch table {
	case 1:
		csvdata = Table1
	case 2:
		csvdata = Table2
	case 3:
		csvdata = Table3
	case 4:
		csvdata = Table4
	case 5:
		csvdata = Table5
	case 6:
		csvdata = Table6
	default:
		return nil, fmt.Errorf("%d is not a valid SEIFA data table", table)
	}

	r := csv.NewReader(strings.NewReader(csvdata))
	record, err := r.Read() // read header
	if err != nil {
		return nil, fmt.Errorf("failed to read header %v", err)
	}
	for {
		record, err = r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) > 0 && record[0] != "" && strings.EqualFold(record[0], poa) {
			return record, nil
		}
	}

	return nil, nil
}

func ParseIntFromFirstMatchingRecord(table int, poa string, column int) (int, bool, error) {
	matchingRecord, err := FirstMatchingRecord(table, poa)
	if err != nil {
		return 0, false, err
	}
	if matchingRecord == nil {
		return 0, false, nil
	}
	var i int
	i, err = strconv.Atoi(matchingRecord[column])
	if err != nil {
		return 0, false, err
	}
	return i, true, nil
}
