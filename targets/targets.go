package targets

import (
	"fmt"

	"github.com/clbanning/mxj"
)

type Target interface {
	SetValue(interface{}) error
}

type MXJ struct {
	Map     mxj.Map
	Path    string
	OnError func(path string, v interface{}, err error) error
}

func (t MXJ) SetValue(v interface{}) error {
	err := t.Map.SetValueForPath(v, t.Path)
	if err != nil {
		if t.OnError != nil { // optional error handler
			return t.OnError(t.Path, v, err)
		}
		return fmt.Errorf("failed to set %s %v", t.Path, err)
	}
	return nil
}

type CSV struct {
	Header       []string
	Record       []string
	ColumnNumber uint
	ColumnName   string
}

func (s CSV) SetValue(v interface{}) error {
	if s.ColumnNumber < 1 {
		if s.Header == nil || len(s.Header) < 1 || s.ColumnName == "" {
			return fmt.Errorf("either a ColumnNumber must be specifed in the range 1 to %d or a Header and ColumnName provided", len(s.Record))
		}
		for index, value := range s.Header {
			if value == s.ColumnName {
				s.ColumnNumber = uint(index + 1)
				break
			}
		}
		if s.ColumnNumber < 1 {
			return fmt.Errorf("ColumnName %s does not exist in Header", s.ColumnName)
		}
	}
	if int(s.ColumnNumber) > len(s.Record) {
		return fmt.Errorf("invalid column %d, record only contains %d columns", s.ColumnNumber, len(s.Record))
	}
	t, ok := v.(string)
	if !ok {
		return fmt.Errorf("value has invalid type %T, expected string", v)
	}
	s.Record[s.ColumnNumber-1] = t
	return nil
}
