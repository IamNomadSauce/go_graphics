package crypto

import (
  "math/rand"
  "github.com/gotk3/gotk3/gtk"
  "github.com/gotk3/gotk3/gdk"
  "github.com/gotk3/gotk3/cairo"

)

type Point struct {
  X, Y float64
}

var points []Point

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

  drawingArea.Connect("draw", drawScatterPlot)

  drawingArea.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.POINTER_MOTION_MASK))
  drawingArea.Connect("button-press-event", onMousePress)
  drawingArea.Connect("motion-notify-event", onMouseMove)

  for i := 0; i < 20; i++ {
    points = append(points, Point{X: rand.Float64(), Y: rand.Float64()})
  }

  box.PackStart(drawingArea, true, true, 0)

  return box, nil
}

func drawScatterPlot(da *gtk.DrawingArea, cr *cairo.Context){
  width := float64(da.GetAllocatedWidth())
  height := float64(da.GetAllocatedHeight())

  cr.SetSourceRGB(1,1,1)
  cr.Paint()
  
  cr.SetSourceRGB(0,0,0)
  cr.MoveTo(10,10)
  cr.LineTo(10,height-10)
  cr.LineTo(width-10,height-10)
  cr.Stroke()
  cr.SetSourceRGB(1,0,0)
  for _, p := range points {
    cr.Arc(10+p.X*(width-20), height-10-p.Y*(height-20), 5, 0, 2*3.14159)
    cr.Fill()
  }

}

func onMousePress(da *gtk.DrawingArea, event *gdk.Event) {
  buttonEvent := gdk.EventButtonNewFromEvent(event) 
  x, y := buttonEvent.X(), buttonEvent.Y()
  width := float64(da.GetAllocatedWidth())
  height := float64(da.GetAllocatedHeight())

  newPoint := Point{
    X: (x-10) / (width - 20),
    Y: 1 - (y - 10) / (height-20),
  }
  points = append(points, newPoint)

  da.QueueDraw()
}

func onMouseMove(da *gtk.DrawingArea, event *gdk.Event) {
  motionEvent := gdk.EventMotionNewFromEvent(event)
  x, y := motionEvent.MotionVal()
  xRoot, yRoot := motionEvent.MotionValRoot()

  width := float64(da.GetAllocatedWidth())
  height := float64(da.GetAllocatedHeight())

  println("Mouse at (window):", x, y)
  println("Mouse at (root):", xRoot, yRoot)

  normalizedX := x / width
  normalizedY := y / height
  println("Normalized coordinates:", normalizedX, normalizedY)

  da.QueueDraw()
}


