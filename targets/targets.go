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
