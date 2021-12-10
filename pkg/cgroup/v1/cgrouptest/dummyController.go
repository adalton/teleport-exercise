package cgrouptest

type Controller interface {
	Name() string
	Apply(path string) error
}

type DummyController struct {
	ControllerName   string
	ApplyReturnValue error
}

func (d *DummyController) Name() string {
	return d.ControllerName
}

func (d *DummyController) Apply(string) error {
	return d.ApplyReturnValue
}
