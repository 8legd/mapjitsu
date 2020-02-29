package mapjitsu

import (
	"fmt"
	"strings"

	"github.com/8legd/mapjitsu/sources"
	"github.com/8legd/mapjitsu/targets"
)

type Pipeline []func(interface{}) (interface{}, error)

type Mapping struct {
	Source    sources.Source
	Transform Pipeline
	Target    targets.Target
}

type Definition struct {
	Mappings []Mapping
}

func (d Definition) Apply() error {
	for _, m := range d.Mappings {

		v, err := m.Source.Value()
		if err != nil {
			return err
		}

		for _, f := range m.Transform {
			v, err = f(v)
			if err != nil {
				return err
			}
		}

		err = m.Target.SetValue(v)
		if err != nil {
			return err
		}

	}
	return nil
}

func ToString(v interface{}) (interface{}, error) {
	if v == nil {
		return "", nil // return nil as empty string
	}
	// otherwise use default fmt
	return fmt.Sprintf("%v", v), nil
}

func ToLower(v interface{}) (interface{}, error) {
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("invalid type %T for function ToLower expected string", v)
	}
	return strings.ToUpper(s), nil
}

func ToUpper(v interface{}) (interface{}, error) {
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("invalid type %T for function ToUpper expected string", v)
	}
	return strings.ToUpper(s), nil
}
