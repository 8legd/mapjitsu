package tests

import (
	"testing"

	"github.com/8legd/mapjitsu/seifa/percentiles"
)

func TestSeifa(t *testing.T) {

	// Test percentiles

	assertMatch := func(index string, poa string, expected int, actual int, exists bool, err error) {
		if err != nil {
			t.Errorf("error reading SEFIA index %s percentile for POA %s %v\n", index, poa, err)
			return
		}
		if exists == false {
			t.Errorf("no matching percentile for POA %s, expected value of %d for SEFIA index %s\n", poa, expected, index)
			return
		}
		if actual != expected {
			t.Errorf("resulting percentile %d for POA %s does not match expected value of %d for SEFIA index %s\n", actual, poa, expected, index)
			return
		}
		t.Logf("SEFIA index %s has a percentile of %d for POA %s", index, actual, poa)
	}

	index := "IRSD" // Index of Relative Socio-economic Disadvantage
	poa := "2000"
	expected := 38
	percentile, exists, err := percentiles.IRSD(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	poa = "6000"
	expected = 75
	percentile, exists, err = percentiles.IRSD(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	poa = "6799"
	expected = 4
	percentile, exists, err = percentiles.IRSD(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	index = "IRSAD" // Index of Relative Socio-economic Advantage and Disadvantage
	poa = "2000"
	expected = 85
	percentile, exists, err = percentiles.IRSAD(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	poa = "6000"
	expected = 88
	percentile, exists, err = percentiles.IRSAD(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	poa = "6799"
	expected = 8
	percentile, exists, err = percentiles.IRSAD(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	index = "IER" // Index of Economic Resources
	poa = "2000"
	expected = 2
	percentile, exists, err = percentiles.IER(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	poa = "6000"
	expected = 8
	percentile, exists, err = percentiles.IER(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	poa = "6799"
	expected = 25
	percentile, exists, err = percentiles.IER(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	index = "IEO" // Index of Education and Occupation
	poa = "2000"
	expected = 94
	percentile, exists, err = percentiles.IEO(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	poa = "6000"
	expected = 92
	percentile, exists, err = percentiles.IEO(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

	poa = "6799"
	expected = 1
	percentile, exists, err = percentiles.IEO(poa)
	assertMatch(index, poa, expected, percentile, exists, err)

}
