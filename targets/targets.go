package targets

import (
	"fmt"

	"github.com/clbanning/mxj"
)

type Target interface {
	SetValue(interface{}) error
}

type MXJ struct {
	Map  mxj.Map
	Path string
}

func (t MXJ) SetValue(v interface{}) error {
	err := t.Map.SetValueForPath(v, t.Path)
	if err != nil {
		return fmt.Errorf("failed to set %s %v", t.Path, err)
	}
	return nil
}
