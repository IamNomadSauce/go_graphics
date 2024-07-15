package crypto

import (
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "time"
    "strconv"
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/gdk"
    "github.com/gotk3/gotk3/glib"
    "gogtk/chart/candlestick"
    "gogtk/cbwebsocket"
)

type TickerMessage struct {
    Type      string `json:"type"`
    ProductID string `json:"product_id"`
    Price     string `json:"price"`
    Time      string `json:"time"`
}

var watchlist = []string{
    "BTC-USD",
    "XLM-USD",
}

var chartInstance *candlestick.Candlestick
var priceUpdateChan chan float64
var CurrentAsset = watchlist[1]

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

    watchlistBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10) // 10 pixels spacing
    if err != nil {
        return nil, err
    }

    // Create a map to store labels for each asset
    assetLabels := make(map[string]*gtk.Label)

    for _, asset := range watchlist {
        assetLabel, err := gtk.LabelNew(asset)
        if err != nil {
            return nil, err
        }

        // Create a frame to contain the label
        frame, err := gtk.FrameNew("")
        if err != nil {
            return nil, err
        }
        frame.Add(assetLabel)

        // Add some padding around the label
        frame.SetMarginStart(5)
        frame.SetMarginEnd(5)
        frame.SetMarginTop(5)
        frame.SetMarginBottom(5)

        watchlistBox.PackStart(frame, false, false, 0)
        assetLabels[asset] = assetLabel
    }

    // Add the watchlistBox to the main box
    box.PackStart(watchlistBox, false, false, 0)

    drawingArea, err := gtk.DrawingAreaNew()
    if err != nil {
        return nil, err
    }

    drawingArea.SetSizeRequest(400, 300)
    candles := generateTestData(100)
    chartInstance, priceUpdateChan = candlestick.NewCandlestick(candles, drawingArea)

    drawingArea.Connect("draw", chartInstance.Draw)
    drawingArea.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.POINTER_MOTION_MASK | gdk.BUTTON_RELEASE_MASK | gdk.SCROLL_MASK))
    drawingArea.Connect("button-press-event", chartInstance.OnMousePress)
    drawingArea.Connect("motion-notify-event", chartInstance.OnMouseMove)
    drawingArea.Connect("button-release-event", chartInstance.OnMouseRelease)
    drawingArea.Connect("scroll-event", chartInstance.OnScroll)

    box.PackStart(drawingArea, true, true, 0)

    messageChannel := make(chan string)

    go cbwebsocket.StartWebSocketClient(watchlist, messageChannel)

    go func() {
        for message := range messageChannel {
            var tickerData TickerMessage
            err := json.Unmarshal([]byte(message), &tickerData)
            if err != nil {
                log.Printf("Error parsing message: %v", err)
                continue
            }
            //if tickerData.Type == "ticker" && tickerData.ProductID == CurrentAsset  {
            if tickerData.Type == "ticker" {
                glib.IdleAdd(func() {
                    if label, exists := assetLabels[tickerData.ProductID]; exists {
                        label.SetText(fmt.Sprintf("%s: %s", tickerData.ProductID, tickerData.Price))
                    }

		    price, err := strconv.ParseFloat(tickerData.Price, 64)
		    if err != nil {
			    fmt.Println("Error parsing price: %v", err)
			    return
		    }
		    if tickerData.ProductID == CurrentAsset  {
			    priceUpdateChan <- price
		    }
                })
            }
        }
    }()

    return box, nil
}

func generateTestData(count int) []candlestick.Candle {
    candles := make([]candlestick.Candle, count)
    baseTime := time.Now().AddDate(0, 0, -count)
    basePrice := 0.10

    for i := 0; i < count; i++ {
        open := basePrice + rand.Float64()*10 - 5
        high := open + rand.Float64()*5
        low := open - rand.Float64()*5
        cls := (open + high + low) / 3
        volume := rand.Float64() * 10000

        candles[i] = candlestick.Candle{
            Time:   baseTime.Add(time.Duration(i) * time.Hour * 24).Unix(),
            Open:   open,
            High:   high,
            Low:    low,
            Close:  cls,
            Volume: volume,
        }
        basePrice = cls
    }
    return candles
}

