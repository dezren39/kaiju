package properties

import (
	"errors"
	"kaiju/engine"
	"kaiju/markup/css/rules"
	"kaiju/markup/document"
	"kaiju/ui"
)

// length|auto|initial|inherit
func (p MarginTop) Process(panel *ui.Panel, elm document.DocElement, values []rules.PropertyValue, host *engine.Host) error {
	if len(values) != 1 {
		return errors.New("MarginTop requires exactly 1 value")
	} else {
		current := panel.Layout().Margin()
		size := marginSizeFromStr(values[0].Str, host.Window)
		panel.Layout().SetMargin(current.X(), size, current.Z(), current.W())
		return nil
	}
}