package scatter

import (
  "fmt"
  "math"
  "github.com/gotk3/gotk3/gtk"
  "github.com/gotk3/gotk3/gdk"
  "github.com/gotk3/gotk3/cairo"
)

type Point struct {
  X, Y float64
}

type Scatter struct {
  Points          []Point
  Scale           float64
  OffsetX         float64
  OffsetY         float64
  Dragging        bool
  LastX           float64
  LastY           float64
  HoveredPoint    *Point
  ClickedPoint    *Point
}

func NewCart(points []Point) *Scatter {
  return &Scatter{
    Points: points,
    Scale: 1.0,
  }
}

func (c *Scatter) Draw(da *gtk.DrawingArea, cr *cairo.Context) {
  width := float64(da.GetAllocatedWidth())
  height := float64(da.GetAllocatedHeight())

  cr.SetSourceRGB(0,0,0)
  cr.Paint()

  // Apply Transforations
  cr.Save()
  cr.Translate(40, height-40)
  cr.Scale(c.Scale, -c.Scale)
  cr.Translate(c.OffsetX/c.Scale, c.OffsetY/c.Scale)

  cr.SetSourceRGB(1,0,0)
  for _, p := range c.Points {
    x := p.X * (width-80)
    y := p.Y * (height-80)
    if c.HoveredPoint != nil && p == *c.HoveredPoint {
      cr.SetSourceRGB(0,1,0)
    } else {
      cr.SetSourceRGB(1,0,0)
    }
    cr.Arc(x, y, 10/c.Scale, 0, 2*math.Pi)
    cr.Fill()
  }
  cr.Restore()

  // Draw Axes
  cr.SetSourceRGB(1,1,1)
  cr.SetLineWidth(2)
  cr.MoveTo(40, height-40)
  cr.LineTo(width-40, height-40) // X-Axis
  cr.MoveTo(40,40)
  cr.LineTo(40, height-40) // Y-Axis
  cr.Stroke()

  minX := -c.OffsetX / c.Scale
  maxX := (width-80 - c.OffsetX) / c.Scale
  minY := c.OffsetY / c.Scale
  maxY := (height-80+c.OffsetY) / c.Scale

  // Axis Labels setup
  cr.SetFontSize(12)

  // X Axis Labels
  for i := 0; i <= 10; i++ {
    x := 40 + (width-80) * float64(i)/10
    y := height - 30
    label := fmt.Sprintf("%.2f", minX+(maxX-minX)*float64(i)/10)
    cr.MoveTo(x,y)
    cr.ShowText(label)
  }
  // Y Axis Labels
  for i := 0; i <= 10; i++ {
    x := 40 + (width-80) * float64(i)/10
    y := height - 30
    label := fmt.Sprintf("%.2f", minY+(maxY-minY)*float64(i)/10)
    cr.MoveTo(x,y)
    cr.ShowText(label)
  }
}

func (c *Scatter) OnMousePress(da *gtk.DrawingArea, event *gdk.Event) {
  buttonEvent := gdk.EventButtonNewFromEvent(event)
  if buttonEvent.Button() == gdk.BUTTON_PRIMARY {
    c.Dragging = true
    c.LastX, c.LastY = buttonEvent.MotionVal()

    x := (buttonEvent.X() - 40 - c.OffsetY) / c.Scale
    y := (float64(da.GetAllocatedHeight()) - buttonEvent.Y() - 40 - c.OffsetY) / c.Scale

    // Check if a point is clicked
    for _, p := range c.Points {
      px := p.X * (float64(da.GetAllocatedWidth()) - 80)
      py := p.Y * (float64(da.GetAllocatedHeight()) - 80)
      if math.Hypot(px-x, py-y) <= 10/c.Scale {
        c.ClickedPoint = &p
        da.QueueDraw()
        return
      }
    }
  }
}

func (c *Scatter) OnMouseMove(da *gtk.DrawingArea, event *gdk.Event) {
  motionEvent := gdk.EventMotionNewFromEvent(event)
  x, y := motionEvent.MotionVal()
  if c.Dragging {
    dx := x - c.LastX
    dy := y - c.LastY
    c.OffsetX += dx
    c.OffsetY -= dy
    c.LastX, c.LastY = x, y

    da.QueueDraw()
  } else {
    // Check if a point is hovered
    hx := (x-40-c.OffsetX) / c.Scale
    hy := (float64(da.GetAllocatedHeight()) - y - 40 - c.OffsetY) / c.Scale
    for _, p := range c.Points {
      px := p.X * (float64(da.GetAllocatedWidth())-80)
      py := p.Y * (float64(da.GetAllocatedHeight())-80)
      if math.Hypot(px-hx, py-hy) <= 10/c.Scale {
        c.HoveredPoint = &p
        da.QueueDraw()
        return
      }
    }
    c.HoveredPoint = nil
    da.QueueDraw()
  }
}

func (c *Scatter) OnMouseRelease(da *gtk.DrawingArea, event *gdk.Event) {
  buttonEvent := gdk.EventButtonNewFromEvent(event)
  if buttonEvent.Button() == gdk.BUTTON_PRIMARY {
    c.Dragging = false
  }
}

func (c *Scatter) OnScroll(da *gtk.DrawingArea, event *gdk.Event) {
  scrollEvent := gdk.EventScrollNewFromEvent(event) 
  direction := scrollEvent.Direction()

  height := float64(da.GetAllocatedHeight())

  // Get Mouse position relative to the drawing area
  x := scrollEvent.X() - 40
  y := height - scrollEvent.Y() - 40

  // Calculate world coordinates before zooming
  worldX := (x - c.OffsetX) / c.Scale
  worldY := (y - c.OffsetY) / c.Scale

  oldScale := c.Scale
  if direction == gdk.SCROLL_UP {
    c.Scale *= 1.1
  } else if direction == gdk.SCROLL_DOWN {
    c.Scale /= 1.1
  }

  // Adjust offsets to keep the point under the cursor foxed
  c.OffsetX += worldX * (oldScale - c.Scale)
  c.OffsetY += worldY * (oldScale - c.Scale)

  da.QueueDraw()
}
