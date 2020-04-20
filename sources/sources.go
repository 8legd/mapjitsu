package sources

import (
	"fmt"

	"github.com/clbanning/mxj"
)

type Source interface {
	Value() (interface{}, error)
}

type Calculated struct {
	Formula func() (interface{}, error)
}

func (s Calculated) Value() (interface{}, error) {
	return s.Formula()
}

type MXJ struct {
	Map     mxj.Map
	Path    string
	OnError func(path string, v interface{}, err error) (interface{}, error)
}

func (s MXJ) Value() (interface{}, error) {
	v, err := s.Map.ValueForPath(s.Path)
	if err != nil {
		if s.OnError != nil { // optional error handler
			return s.OnError(s.Path, v, err)
		}
		return nil, fmt.Errorf("failed to return %s %v", s.Path, err)
	}
	return v, nil
}

type CSV struct {
	Header       []string
	Record       []string
	ColumnNumber uint
	ColumnName   string
}

func (s CSV) Value() (interface{}, error) {
	if s.ColumnNumber < 1 {
		if s.Header == nil || len(s.Header) < 1 || s.ColumnName == "" {
			return nil, fmt.Errorf("either a ColumnNumber must be specifed in the range 1 to %d or a Header and ColumnName provided", len(s.Record))
		}
		for index, value := range s.Header {
			if value == s.ColumnName {
				s.ColumnNumber = uint(index + 1)
				break
			}
		}
		if s.ColumnNumber < 1 {
			return nil, fmt.Errorf("ColumnName %s does not exist in Header", s.ColumnName)
		}
	}
	if int(s.ColumnNumber) > len(s.Record) {
		return nil, fmt.Errorf("invalid column %d, record only contains %d columns", s.ColumnNumber, len(s.Record))
	}
	return s.Record[s.ColumnNumber-1], nil
}
