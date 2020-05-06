package mapjitsu

type Mapping struct {
	Source    Source
	Transform Pipeline
	Target    Target
}

type Source interface {
	Value() (interface{}, error)
}

// The SourceFunc type is an adapter to allow the use of
// ordinary functions as Sources. If f is a function
// with the appropriate signature, SourceFunc(f) is a
// Source that returns f().
type SourceFunc func() (interface{}, error)

// Value returns f().
func (f SourceFunc) Value() (interface{}, error) {
	return f()
}

type Pipeline []func(interface{}) (interface{}, error)

type Target interface {
	SetValue(interface{}) error
}

// The TargetFunc type is an adapter to allow the use of
// ordinary functions as Targets. If f is a function
// with the appropriate signature, TargetFunc(f) is a
// Source that calls f.
type TargetFunc func(interface{}) error

// SetValue calls f.
func (f TargetFunc) SetValue(v interface{}) error {
	return f(v)
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
