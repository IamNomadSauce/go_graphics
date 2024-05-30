package common

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

type DrawState struct {
	Shape string
	X, Y  float64
}

func ClearContainer(container *gtk.Box) {
	children := container.GetChildren()
	children.Foreach(func(item interface{}) {
		child, ok := item.(*gtk.Widget)
		if ok {
			container.Remove(child)
		}
	})
}

func UpdateGUI(container *gtk.Box, drawState *DrawState) {
	fmt.Println("UpdateGUI")
	ClearContainer(container)

	shapeLabel, _ := gtk.LabelNew(drawState.Shape)
	xLabel, _ := gtk.LabelNew(fmt.Sprintf("%.2f", drawState.X))
	yLabel, _ := gtk.LabelNew(fmt.Sprintf("%.2f", drawState.Y))

	container.PackStart(shapeLabel, false, false, 0)
	container.PackStart(xLabel, false, false, 0)
	container.PackStart(yLabel, false, false, 0)

	container.ShowAll()
}
