package crypto

import (
    _"fmt"
    _"math"
    "math/rand"
    "time"
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/gdk"
    _"github.com/gotk3/gotk3/cairo"
    //"gogtk/chart/scatter"
    //"gogtk/db/mysql"
    "gogtk/chart/candlestick"
    "gogtk/cbwebsocket"
    
)

func CryptoPage() (*gtk.Box, error) {
    box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    if err != nil {
        return nil, err
    }


    // Start websocket server
    go cbwebsocket.StartWebSocketClient()

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

    candles := generateTestData(100)

    chartInstance := candlestick.NewCandlestick(candles)

    drawingArea.Connect("draw", chartInstance.Draw)



    //points := []scatter.Point{}
    //for i := 0; i < 100; i++ {
     // points = append(points, scatter.Point{X: rand.Float64(), Y: rand.Float64()})
    //}

    //chartInstance := scatter.NewChart(points)

    drawingArea.Connect("draw", chartInstance.Draw)
    drawingArea.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.POINTER_MOTION_MASK | gdk.BUTTON_RELEASE_MASK | gdk.SCROLL_MASK))
    drawingArea.Connect("button-press-event", chartInstance.OnMousePress)
    drawingArea.Connect("motion-notify-event", chartInstance.OnMouseMove)
    drawingArea.Connect("button-release-event", chartInstance.OnMouseRelease)
    drawingArea.Connect("scroll-event", chartInstance.OnScroll)

    box.PackStart(drawingArea, true, true, 0)

    return box, nil
}

func generateTestData(count int) []candlestick.Candle {
  candles := make([]candlestick.Candle, count)
  baseTime := time.Now().AddDate(0,0, -count)
  basePrice := 100.0

  for i := 0; i < count; i++ {
    open := basePrice + rand.Float64()*10-5
    high := open + rand.Float64()*5
    low := open - rand.Float64()*5
    cls := (open + high + low) / 3
    volume := rand.Float64() * 10000

    candles[i] = candlestick.Candle{
      Time:   baseTime.Add(time.Duration(i) * time.Hour * 24).Unix(),
      Open: open,
      High: high,
      Low: low,
      Close: cls,
      Volume: volume,
    }
    basePrice = cls
  }
  return candles
}
