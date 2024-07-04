package crypto

import (
    "fmt"
    "math"
    "math/rand"
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/gdk"
    "github.com/gotk3/gotk3/cairo"
)

type Point struct {
    X, Y float64
}

var points []Point
var scale float64 = 1.0
var offsetX, offsetY float64 = 0.0, 0.0
var dragging bool = false
var lastX, lastY float64

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

    for i := 0; i < 500; i++ {
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
    cr.Translate(40 + offsetX, height - 40 + offsetY)
    cr.Scale(scale, -scale)

    // Draw scatter points
    cr.SetSourceRGB(1, 0, 0)
    for _, p := range points {
        cr.Arc(p.X, p.Y, 5/scale, 0, 2*math.Pi)
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
    minY := -(height - 80 + offsetY) / scale
    maxY := -offsetY / scale

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
    }
}

func onMouseMove(da *gtk.DrawingArea, event *gdk.Event) {
    if dragging {
        motionEvent := gdk.EventMotionNewFromEvent(event)
        x, y := motionEvent.MotionVal()
        dx := x - lastX
        dy := y - lastY
        offsetX += dx
        offsetY += dy
        lastX, lastY = x, y

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

    x := scrollEvent.X() - 40
    y := scrollEvent.Y() - 40

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


