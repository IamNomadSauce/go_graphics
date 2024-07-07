package candlestick

import (
  "fmt"
  "math"
  "github.com/gotk3/gotk3/gtk"
  "github.com/gotk3/gotk3/gdk"
  "github.com/gotk3/gotk3/cairo"
  
)
type Candle struct {
  Time    int64
  Open    float64
  High    float64
  Low    float64
  Close    float64
  Volume    float64
}

type Candlestick struct {
    Candles []Candle
    Scale   float64
    OffsetX float64
    OffsetY float64
    Dragging bool
    LastX, LastY float64
    HoveredCandle *Candle
    ClickedCandle *Candle
    // Config, candle colors, linetools, etc
}
func NewCandlestick(candles []Candle) *Candlestick {
    return &Candlestick{
        Candles: candles,
        Scale:   1.0,
    }
}

func (c *Candlestick) Draw(da *gtk.DrawingArea, cr *cairo.Context) {
    width := float64(da.GetAllocatedWidth())
    height := float64(da.GetAllocatedHeight())

    cr.SetSourceRGB(0, 0, 0)
    cr.Paint()

    // Calculate price and time ranges
    minTime := c.Candles[0].Time
    maxTime := c.Candles[len(c.Candles)-1].Time
    minPrice := math.Inf(1)
    maxPrice := math.Inf(-1)
    for _, candle := range c.Candles {
        minPrice = math.Min(minPrice, candle.Low)
        maxPrice = math.Max(maxPrice, candle.High)
    }

    // Calculate scaling factors
    timeRange := float64(maxTime - minTime)
    if timeRange == 0 {
        timeRange = 1 // Prevent division by zero
    }
    timeScale := (width - 80) / timeRange
    priceRange := maxPrice - minPrice
    if priceRange == 0 {
        priceRange = 1 // Prevent division by zero
    }
    priceScale := (height - 80) / priceRange

    // Apply Transformations
    cr.Save()
    cr.Translate(40, height-40)
    cr.Scale(c.Scale, -c.Scale)
    cr.Translate(c.OffsetX/c.Scale, c.OffsetY/c.Scale)

    // Draw candlesticks
    candleWidth := timeScale * 0.8 // 80% of the time slot
    if candleWidth < 1 {
        candleWidth = 10 / c.Scale // Ensure a minimum width for visibility
    }
    for _, candle := range c.Candles {
        x := float64(candle.Time-minTime) * timeScale
        yOpen := (candle.Open - minPrice) * priceScale
        yClose := (candle.Close - minPrice) * priceScale
        yHigh := (candle.High - minPrice) * priceScale
        yLow := (candle.Low - minPrice) * priceScale

        // Draw candle wick
        cr.SetSourceRGB(1, 1, 1) // White for the wick
        cr.SetLineWidth(1 / c.Scale)
        cr.MoveTo(x+candleWidth/2, yLow)
        cr.LineTo(x+candleWidth/2, yHigh)
        cr.Stroke()

        // Draw candle body
        if candle.Close > candle.Open {
            cr.SetSourceRGB(0, 1, 0) // Green for bullish
        } else {
            cr.SetSourceRGB(1, 0, 0) // Red for bearish
        }
        cr.Rectangle(x, math.Min(yOpen, yClose), candleWidth, math.Abs(yClose-yOpen))
        cr.Fill()
    }

    cr.Restore()

    // Draw Axes
    cr.SetSourceRGB(1, 1, 1)
    cr.SetLineWidth(2)
    cr.MoveTo(40, height-40)
    cr.LineTo(width-40, height-40) // X-Axis
    cr.MoveTo(40, 40)
    cr.LineTo(40, height-40) // Y-Axis
    cr.Stroke()

    // Axis Labels setup
    cr.SetFontSize(12)

    // Calculate visible range
    visibleMinTime := minTime + int64((-c.OffsetX / c.Scale) / timeScale)
    visibleMaxTime := minTime + int64(((width - 80) / c.Scale - c.OffsetX / c.Scale) / timeScale)
    visibleMinPrice := minPrice + (-c.OffsetY / c.Scale) / priceScale
    visibleMaxPrice := minPrice + ((height - 80) / c.Scale - c.OffsetY / c.Scale) / priceScale

    // X Axis Labels (Time)
    for i := 0; i <= 10; i++ {
        x := 40 + (width-80) * float64(i)/10
        y := height - 30
        time := visibleMinTime + int64(float64(visibleMaxTime-visibleMinTime)*float64(i)/10)
        label := fmt.Sprintf("%d", time)
        cr.MoveTo(x, y)
        cr.ShowText(label)
    }

    // Y Axis Labels (Price)
    for i := 0; i <= 10; i++ {
        x := 10
        y := height - 40 - (height-80)*float64(i)/10
        price := visibleMinPrice + (visibleMaxPrice-visibleMinPrice)*float64(i)/10
        label := fmt.Sprintf("%.2f", price)
        cr.MoveTo(float64(x), y)
        cr.ShowText(label)
    }
}



func (c *Candlestick) OnMousePress(da *gtk.DrawingArea, event *gdk.Event) {
  buttonEvent := gdk.EventButtonNewFromEvent(event)
  if buttonEvent.Button() == gdk.BUTTON_PRIMARY {
    c.Dragging = true
    c.LastX, c.LastY = buttonEvent.MotionVal()

    //x := (buttonEvent.X() - 40 - c.OffsetY) / c.Scale
    //y := (float64(da.GetAllocatedHeight()) - buttonEvent.Y() - 40 - c.OffsetY) / c.Scale

    // Check if a point is clicked
    //for _, p := range c.Candles {
      //px := p.X * (float64(da.GetAllocatedWidth()) - 80)
      //py := p.Y * (float64(da.GetAllocatedHeight()) - 80)
      //if math.Hypot(px-x, py-y) <= 10/c.Scale {
      //  c.ClickedCandle = &p
      //  da.QueueDraw()
      //  return
      //}
    //}
  }
}

func (c *Candlestick) OnMouseMove(da *gtk.DrawingArea, event *gdk.Event) {
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
    //hx := (x-40-c.OffsetX) / c.Scale
    //hy := (float64(da.GetAllocatedHeight()) - y - 40 - c.OffsetY) / c.Scale
    //for _, p := range c.Candles {
    //  px := p.X * (float64(da.GetAllocatedWidth())-80)
    //  py := p.Y * (float64(da.GetAllocatedHeight())-80)
    //  if math.Hypot(px-hx, py-hy) <= 10/c.Scale {
    //    c.HoveredCandle = &p
    //    da.QueueDraw()
    //    return
    //  }
    //}
    c.HoveredCandle = nil
    da.QueueDraw()
  }
}

func (c *Candlestick) OnMouseRelease(da *gtk.DrawingArea, event *gdk.Event) {
  buttonEvent := gdk.EventButtonNewFromEvent(event)
  if buttonEvent.Button() == gdk.BUTTON_PRIMARY {
    c.Dragging = false
  }
}

func (c *Candlestick) OnScroll(da *gtk.DrawingArea, event *gdk.Event) {
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
