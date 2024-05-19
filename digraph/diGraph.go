package digraph

import (
	"fmt"
	"go_graphics/common"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// ----------------------------------------------------
// DiGraph Page
// ----------------------------------------------------
func DiGraphPage() *gtk.Box {
	fmt.Println("------------------------------")
	fmt.Println("DiGraph")
	fmt.Println("------------------------------\n")

	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		fmt.Println("Error Creating DiGraph", err)
	}

	sidebar, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		fmt.Println("Error Creating DiGraph:SideBar", err)
	}
	dynamicContent, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		fmt.Println("Error creating dynamic content container", err)
	}

	drawState := &common.DrawState{Shape: "Square", X: 100, Y: 100}

	setupShapeButtons(sidebar, drawState, dynamicContent)
	setupDynamicContent(dynamicContent, drawState)

	container.PackStart(sidebar, false, false, 10)
	container.PackStart(dynamicContent, true, true, 10)

	// updateGUI(dynamicContent, common.drawState)

	// container.ShowAll()

	return container
}

func setupDynamicContent(dynamicContent *gtk.Box, drawState *common.DrawState) {
	dynamicContent.Connect("button-press-event", func(da *gtk.Box, event *gdk.Event) {
		buttonEvent := gdk.EventButton{Event: event}
		drawState.X, drawState.Y = buttonEvent.X(), buttonEvent.Y()
		common.UpdateGUI(dynamicContent, drawState)
	})
}

// ----------------------------------------------------
// Shape Button Click Handler
func setupShapeButtons(sidebar *gtk.Box, ds *common.DrawState, dynamicContent *gtk.Box) {
	fmt.Println("SetupShapeButtons")
	menuItems := []string{
		"Square",
		"Circle",
		"Line",
		"Box",
	}

	for _, item := range menuItems {
		item := item
		btn, err := gtk.ButtonNewWithLabel(item)
		if err != nil {
			fmt.Println("Error Creating Shape Button", err)
			continue
		}
		btn.Connect("clicked", func() {
			ds.Shape = item
			fmt.Println("Shape Selected:", item)
			common.UpdateGUI(dynamicContent, ds)
		})
		sidebar.PackStart(btn, false, false, 5)
	}
}

// ----------------------------------------------------
// Drawing Area Setup/staging
func setupDrawingArea(drawArea *gtk.DrawingArea, drawState *common.DrawState, container *gtk.Box) {
	drawArea.Connect("button-press-event", func(da *gtk.DrawingArea, event *gdk.Event) {
		buttonEvent := gdk.EventButton{Event: event}
		drawState.X, drawState.Y = buttonEvent.X(), buttonEvent.Y()
		drawArea.QueueDraw()
		common.UpdateGUI(container, drawState)
	})
	drawArea.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		drawShape(cr, drawState)
	})
}

func drawShape(cr *cairo.Context, ds *common.DrawState) {
	cr.SetLineWidth(2)
	fmt.Println(ds.Shape, ds.X, ds.Y)
	switch ds.Shape {
	case "Square":
		cr.Rectangle(ds.X, ds.Y, 50, 50)
	case "Circle":
		cr.Arc(ds.X, ds.Y, 25, 0, 2*3.14159)
	case "Line":
		cr.MoveTo(ds.X, ds.Y)
		cr.LineTo(ds.X-25, ds.Y)
	case "Box":
		cr.Rectangle(ds.X, ds.Y, 100, 50)
	}
}
