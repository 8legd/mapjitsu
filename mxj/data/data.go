package data

import (
	"fmt"

	"github.com/clbanning/mxj"
)

type Source struct {
	Map     mxj.Map
	Path    string
	OnError func(path string, v interface{}, err error) (interface{}, error)
}

func (s Source) Value() (interface{}, error) {
	v, err := s.Map.ValueForPath(s.Path)
	if err != nil {
		if s.OnError != nil { // optional error handler
			return s.OnError(s.Path, v, err)
		}
		return nil, fmt.Errorf("failed to return %s %v", s.Path, err)
	}
	return v, nil
}

type Target struct {
	Map     mxj.Map
	Path    string
	OnError func(path string, v interface{}, err error) error
}

func (t Target) SetValue(v interface{}) error {
	err := t.Map.SetValueForPath(v, t.Path)
	if err != nil {
		if t.OnError != nil { // optional error handler
			return t.OnError(t.Path, v, err)
		}
		return fmt.Errorf("failed to set %s %v", t.Path, err)
	}
	return nil
}
