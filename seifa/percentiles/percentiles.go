package percentiles

import (
	"github.com/8legd/mapjitsu/seifa/data"
)

// Index of Relative Socio-economic Disadvantage
// NOTE: Postal Areas (POAs) are an ABS approximation of postcodes
// see https://www.abs.gov.au/ausstats/abs@.nsf/Lookup/by%20Subject/1270.0.55.003~July%202016~Main%20Features~Postal%20Areas%20(POA)~8
func IRSD(poa string) (int, bool, error) {
	return data.ParseIntFromFirstMatchingRecord(2, poa, 5)
}

// Index of Relative Socio-economic Advantage and Disadvantage
// NOTE: Postal Areas (POAs) are an ABS approximation of postcodes
// see https://www.abs.gov.au/ausstats/abs@.nsf/Lookup/by%20Subject/1270.0.55.003~July%202016~Main%20Features~Postal%20Areas%20(POA)~8
func IRSAD(poa string) (int, bool, error) {
	return data.ParseIntFromFirstMatchingRecord(3, poa, 5)
}

// Index of Economic Resources
// NOTE: Postal Areas (POAs) are an ABS approximation of postcodes
// see https://www.abs.gov.au/ausstats/abs@.nsf/Lookup/by%20Subject/1270.0.55.003~July%202016~Main%20Features~Postal%20Areas%20(POA)~8
func IER(poa string) (int, bool, error) {
	return data.ParseIntFromFirstMatchingRecord(4, poa, 5)
}

// Index of Education and Occupation
// NOTE: Postal Areas (POAs) are an ABS approximation of postcodes
// see https://www.abs.gov.au/ausstats/abs@.nsf/Lookup/by%20Subject/1270.0.55.003~July%202016~Main%20Features~Postal%20Areas%20(POA)~8
func IEO(poa string) (int, bool, error) {
	return data.ParseIntFromFirstMatchingRecord(5, poa, 5)
}
