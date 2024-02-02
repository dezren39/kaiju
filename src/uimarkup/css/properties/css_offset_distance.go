package properties

import (
	"errors"
	"kaiju/engine"
	"kaiju/ui"
	"kaiju/uimarkup/css/rules"
	"kaiju/uimarkup/markup"
)

func (p OffsetDistance) Process(panel *ui.Panel, elm markup.DocElement, values []rules.PropertyValue, host *engine.Host) error {
	problems := []error{errors.New("OffsetDistance not implemented")}

	return problems[0]
}