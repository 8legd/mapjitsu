package sources

import (
	"fmt"

	"github.com/clbanning/mxj"
)

type Source interface {
	Value() (interface{}, error)
}

type Calculated struct {
	Result interface{}
}

func (s Calculated) Value() (interface{}, error) {
	return s.Result, nil
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
