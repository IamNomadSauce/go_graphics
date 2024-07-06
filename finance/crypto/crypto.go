package crypto

import (
    "fmt"
    "math"
    "math/rand"
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/gdk"
    "github.com/gotk3/gotk3/cairo"
    "gogtk/chart/scatter"
)

type Point struct {
    X, Y float64
}
type Square struct {
    X, Y float64
}

var points []Point
var squares []Square
var scale float64 = 1.0
var offsetX, offsetY float64 = 0.0, 0.0
var dragging bool = false
var lastX, lastY float64
var hoveredPoint *Point
var hoveredSquare *Square

var clickedSquare *Square
var clickedPoint *Point

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

    drawingArea.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.POINTER_MOTION_MASK | gdk.BUTTON_RELEASE_MASK | gdk.SCROLL_MASK))
    drawingArea.Connect("button-press-event", onMousePress)
    drawingArea.Connect("motion-notify-event", onMouseMove)
    drawingArea.Connect("button-release-event", onMouseRelease)
    drawingArea.Connect("scroll-event", onScroll)

    for i := 0; i < 100; i++ {
        points = append(points, Point{X: rand.Float64(), Y: rand.Float64()})
    }

    box.PackStart(drawingArea, true, true, 0)

    return box, nil
}

func drawScatterPlot(da *gtk.DrawingArea, cr *cairo.Context) {
    width := float64(da.GetAllocatedWidth())
    height := float64(da.GetAllocatedHeight())

    cr.SetSourceRGB(0, 0, 0)
    cr.Paint()

    // Apply transformations
    cr.Save()
    cr.Translate(40, height - 40)  // Move origin to bottom-left corner of plot area
    cr.Scale(scale, -scale)  // Flip Y-axis and apply zoom
    cr.Translate(offsetX/scale, offsetY/scale)  // Apply pan

    // Draw scatter points
    cr.SetSourceRGB(1, 0, 0)
    size := 50.0
    for _, p := range points {
      x := p.X * (width-80)
      y := p.Y * (height-80)
      if hoveredPoint != nil && p == *hoveredPoint {
        cr.SetSourceRGB(0,1,0)
      } else {
        cr.SetSourceRGB(1,0,0)
      }
      cr.Arc(x, y, 10/scale, 0, 2*math.Pi)
      cr.Fill()
    }
    for _, s := range squares {
      x := s.X * (width-80)
      y := s.Y * (height-80)
      if hoveredSquare != nil && s == *hoveredSquare {
        cr.SetSourceRGB(0,1,0)
      } else {
        cr.SetSourceRGB(1,0,0)
      }
      cr.Rectangle(x, y, size, size)
      cr.Fill()
    }

    cr.Restore()

    // Draw Axes
    cr.SetSourceRGB(1, 1, 1)
    cr.SetLineWidth(2)
    cr.MoveTo(40, height-40)
    cr.LineTo(width-40, height-40) // X-axis
    cr.MoveTo(40, 40)
    cr.LineTo(40, height-40) // Y-axis
    cr.Stroke()

    // Calculate the visible range
    minX := -offsetX / scale
    maxX := (width - 80 - offsetX) / scale
    minY := offsetY / scale
    maxY := (height - 80 + offsetY) / scale

    // Setting up Axis Labels
    cr.SetFontSize(12)

    // X axis labels
    for i := 0; i <= 10; i++ {
        x := 40 + (width-80)*float64(i)/10
        y := height - 30
        label := fmt.Sprintf("%.2f", minX + (maxX-minX)*float64(i)/10)
        cr.MoveTo(x, y)
        cr.ShowText(label)
    }

    // Y-axis labels
    for i := 0; i <= 10; i++ {
        x := 10.0
        y := height - 40 - (height-80)*float64(i)/10
        label := fmt.Sprintf("%.2f", minY + (maxY-minY)*float64(i)/10)
        cr.MoveTo(x, y)
        cr.ShowText(label)
    }
}

func onMousePress(da *gtk.DrawingArea, event *gdk.Event) {
    buttonEvent := gdk.EventButtonNewFromEvent(event)
    if buttonEvent.Button() == gdk.BUTTON_PRIMARY {
        dragging = true
        lastX, lastY = buttonEvent.MotionVal()

        x := (buttonEvent.X() - 40 - offsetX) / scale
        y := (float64(da.GetAllocatedHeight()) - buttonEvent.Y() - 40 - offsetY) / scale

        // Check if a point is clicked
        for _, p := range points {
          px := p.X * (float64(da.GetAllocatedWidth()) - 80)
          py := p.Y * (float64(da.GetAllocatedHeight()) - 80)
          if math.Hypot(px-x, py-y) <= 10/scale {
            clickedPoint = &p
            squares = append(squares, Square{X: p.X, Y: p.Y})
            da.QueueDraw()
            return
          }
        }
        // Check if a point is clicked
        for _, s := range squares {
          px := s.X * (float64(da.GetAllocatedWidth()) - 80)
          py := s.Y * (float64(da.GetAllocatedHeight()) - 80)
          if math.Hypot(px-x, py-y) <= 10/scale {
            clickedSquare = &s
            da.QueueDraw()
            return
          }
        }

    }
}

func onMouseMove(da *gtk.DrawingArea, event *gdk.Event) {
  motionEvent := gdk.EventMotionNewFromEvent(event)
  x, y := motionEvent.MotionVal()
  if dragging {
      dx := x - lastX
      dy := y - lastY
      offsetX += dx
      offsetY -= dy  // Invert the Y-axis movement
      lastX, lastY = x, y

      da.QueueDraw()
  } else {
    // Check if a point is hovered
    hx := (x - 40 - offsetX) / scale
    hy := (float64(da.GetAllocatedHeight()) - y - 40 - offsetY) / scale
    for _, p := range points {
      px := p.X * (float64(da.GetAllocatedWidth()) - 80)
      py := p.Y * (float64(da.GetAllocatedHeight()) - 80)
      if math.Hypot(px-hx, py-hy) <= 10/scale {
        hoveredPoint = &p
        da.QueueDraw()
        return
      }
    }
    for _, s := range squares {
      px := s.X * (float64(da.GetAllocatedWidth()) - 80)
      py := s.Y * (float64(da.GetAllocatedHeight()) - 80)
      if math.Hypot(px-hx, py-hy) <= 50/scale {
        hoveredSquare = &s
        da.QueueDraw()
        return
      }
    }
    hoveredPoint = nil
    hoveredSquare = nil
    da.QueueDraw()
  }
}

func onMouseRelease(da *gtk.DrawingArea, event *gdk.Event) {
    buttonEvent := gdk.EventButtonNewFromEvent(event)
    if buttonEvent.Button() == gdk.BUTTON_PRIMARY {
        dragging = false
    }
}

func onScroll(da *gtk.DrawingArea, event *gdk.Event) {
    scrollEvent := gdk.EventScrollNewFromEvent(event)
    direction := scrollEvent.Direction()

    //width := float64(da.GetAllocatedWidth())
    height := float64(da.GetAllocatedHeight())

    // Get mouse position relative to the drawing area
    x := scrollEvent.X() - 40
    y := height - scrollEvent.Y() - 40  // Flip Y coordinate

    // Calculate world coordinates before zooming
    worldX := (x - offsetX) / scale
    worldY := (y - offsetY) / scale

    oldScale := scale
    if direction == gdk.SCROLL_UP {
        scale *= 1.1
    } else if direction == gdk.SCROLL_DOWN {
        scale /= 1.1
    }

    // Adjust offsets to keep the point under the cursor fixed
    offsetX += worldX * (oldScale - scale)
    offsetY += worldY * (oldScale - scale)

    da.QueueDraw()
}

