package candlestick

import (
  "fmt"
  "math"
  "github.com/gotk3/gotk3/gtk"
  _"github.com/gotk3/gotk3/gdk"
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

    // Apply Transformations
    cr.Save()
    cr.Translate(40, height-40)
    cr.Scale(c.Scale, -c.Scale)
    cr.Translate(c.OffsetX/c.Scale, c.OffsetY/c.Scale)

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
    timeScale := (width - 80) / float64(maxTime - minTime)
    priceScale := (height - 80) / (maxPrice - minPrice)

    // Draw candlesticks
    candleWidth := timeScale * 0.95 // 80% of the time slot
    for _, candle := range c.Candles {
        x := float64(candle.Time-minTime) * timeScale
        yOpen := (candle.Open - minPrice) * priceScale
        yClose := (candle.Close - minPrice) * priceScale
        yHigh := (candle.High - minPrice) * priceScale
        yLow := (candle.Low - minPrice) * priceScale

        // Draw candle body
        if candle.Close > candle.Open {
            cr.SetSourceRGB(0, 1, 0) // Green for bullish
        } else {
            cr.SetSourceRGB(1, 0, 0) // Red for bearish
        }
        cr.Rectangle(x, math.Min(yOpen, yClose), candleWidth, math.Abs(yClose-yOpen))
        cr.Fill()

        // Draw candle wick
        cr.SetSourceRGB(1, 1, 1) // White for the wick
        cr.SetLineWidth(1 / c.Scale)
        cr.MoveTo(x+candleWidth/2, yLow)
        cr.LineTo(x+candleWidth/2, yHigh)
        cr.Stroke()
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

    // X Axis Labels (Time)
    for i := 0; i <= 10; i++ {
        x := 40 + (width-80) * float64(i)/10
        y := height - 30
        time := minTime + int64(float64(maxTime-minTime)*float64(i)/10)
        label := fmt.Sprintf("%d", time)
        cr.MoveTo(x, y)
        cr.ShowText(label)
    }

    // Y Axis Labels (Price)
    for i := 0; i <= 10; i++ {
        x := 10
        y := height - 40 - (height-80)*float64(i)/10
        price := minPrice + (maxPrice-minPrice)*float64(i)/10
        label := fmt.Sprintf("%.2f", price)
        cr.MoveTo(float64(x), y)
        cr.ShowText(label)
    }
}

