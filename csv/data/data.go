package data

import (
	"fmt"
)

type Source struct {
	Header       []string
	Record       []string
	ColumnNumber uint
	ColumnName   string
}

func (s Source) Value() (interface{}, error) {
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

type Target struct {
	Header       []string
	Record       []string
	ColumnNumber uint
	ColumnName   string
}

func (t Target) SetValue(v interface{}) error {
	if t.ColumnNumber < 1 {
		if t.Header == nil || len(t.Header) < 1 || t.ColumnName == "" {
			return fmt.Errorf("either a ColumnNumber must be specifed in the range 1 to %d or a Header and ColumnName provided", len(t.Record))
		}
		for index, value := range t.Header {
			if value == t.ColumnName {
				t.ColumnNumber = uint(index + 1)
				break
			}
		}
		if t.ColumnNumber < 1 {
			return fmt.Errorf("ColumnName %s does not exist in Header", t.ColumnName)
		}
	}
	if int(t.ColumnNumber) > len(t.Record) {
		return fmt.Errorf("invalid column %d, record only contains %d columns", t.ColumnNumber, len(t.Record))
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("value has invalid type %T, expected string", v)
	}
	t.Record[t.ColumnNumber-1] = s
	return nil
}
