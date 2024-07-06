package crypto

import (
    _"fmt"
    _"math"
    "math/rand"
    "github.com/gotk3/gotk3/gtk"
    _"github.com/gotk3/gotk3/gdk"
    _"github.com/gotk3/gotk3/cairo"
    "gogtk/chart/scatter"
)

func CryptoPage() (*gtk.Box, error) {
    box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        return nil, err
    }

    label, err := gtk.LabelNew("CryptoTab")
    if err != nil {
        return nil, err
    }
    box.PackStart(label, false, false, 0)

    drawingArea, err := gtk.DrawingAreaNew()
    if err != nil {
        return nil, err
    }

    drawingArea.SetSizeRequest(400, 300)

    points := []scatter.Point{}
    for i := 0; i < 100; i++ {
      points = append(points, scatter.Point{X: rand.Float64(), Y: rand.Float64()})
    }

    chartInstance := scatter.NewChart(points)

    drawingArea.Connect("draw", chartInstance.Draw)
    //drawingArea.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.POINTER_MOTION_MASK | gdk.BUTTON_RELEASE_MASK | gdk.SCROLL_MASK))
    drawingArea.Connect("button-press-event", chartInstance.OnMousePress)
    drawingArea.Connect("motion-notify-event", chartInstance.OnMouseMove)
    drawingArea.Connect("button-release-event", chartInstance.OnMouseRelease)
    drawingArea.Connect("scroll-event", chartInstance.OnScroll)

    box.PackStart(drawingArea, true, true, 0)

    return box, nil
}
