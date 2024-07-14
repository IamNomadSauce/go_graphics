package crypto

import (
	"encoding/json"
    "fmt"
    _"math"
    "math/rand"
    "time"
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/gdk"
    "github.com/gotk3/gotk3/glib"
    _"github.com/gotk3/gotk3/cairo"
    //"gogtk/chart/scatter"
    //"gogtk/db/mysql"
    "gogtk/chart/candlestick"
    "gogtk/cbwebsocket"
    
)

type TickerMessage struct {
	Type string `json:"type"`
	ProductID string `json:"product_id"`
	Price string `json:"price"`
	Time string `json:"Time"`
}

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
    candles := generateTestData(100)
    chartInstance := candlestick.NewCandlestick(candles)

    drawingArea.Connect("draw", chartInstance.Draw)
    drawingArea.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.POINTER_MOTION_MASK | gdk.BUTTON_RELEASE_MASK | gdk.SCROLL_MASK))
    drawingArea.Connect("button-press-event", chartInstance.OnMousePress)
    drawingArea.Connect("motion-notify-event", chartInstance.OnMouseMove)
    drawingArea.Connect("button-release-event", chartInstance.OnMouseRelease)
    drawingArea.Connect("scroll-event", chartInstance.OnScroll)


    box.PackStart(drawingArea, true, true, 0)

    messageChannel := make(chan string)


    go cbwebsocket.StartWebSocketClient(messageChannel) 

    go func() {
	    for message := range messageChannel {
		    var tickerData TickerMessage
		    err := json.Unmarshal([]byte(message), &tickerData)
		    if err != nil {
			    fmt.Println("Error parsing message: %v", err)
			    continue
		    }
		    if tickerData.Type == "ticker" {
			    glib.IdleAdd(func() {
				    label.SetText(fmt.Sprintf("Product: %s, Price %s, Time: %s", tickerData.ProductID, tickerData.Price, tickerData.Time))

			    })

		    }
	    }
    }()


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
