package coinbase

import (
	"fmt"
)

func GetCandles(symbol, tf string) () {
} 







func getCBCandlesAuth(apiKey string, apiSecret string, symbol string, tf int, tflbl string, start int64, end int64) []Candle {
	// fmt.Println("AUTH Coinbase Candles", symbol, tf, start, end)

	timestamp := time.Now().Unix()
	path := fmt.Sprintf("/api/v3/brokerage/products/%s/candles", symbol)
	query := fmt.Sprintf("?granularity=%s&start=%d&end=%d", tflbl, start, end)
	// message := fmt.Sprintf("%d%s%s%s", timestamp, "GET", path, "")
	signature := getCBSign(apiSecret, timestamp, "GET", path, "")

	// fmt.Println("VARS", startTime, endTime.Unix(), start, end)
	// fmt.Println("Message", path+query)

	req, err := http.NewRequest("GET", "https://api.coinbase.com"+path+query, nil)
	if err != nil {
		fmt.Println("ERROR", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("CB-ACCESS-SIGN", signature)
	req.Header.Add("CB-ACCESS-TIMESTAMP", fmt.Sprintf("%d", timestamp))
	req.Header.Add("CB-ACCESS-KEY", apiKey)
	req.Header.Add("CB-VERSION", "2015-07-22")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR", err)
	}

	// fmt.Println("Response::", resp)
	defer resp.Body.Close()

	var data ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("ERROR", err)
	}

	if len(data.Candles) == 0 {
		// fmt.Println("UNDEFINED-Coinbase", symbol, tf)
		return nil
	}

	// fmt.Println("DATA", )

	candles := make([]Candle, len(data.Candles))

	for i, el := range data.Candles {

		t, err := strconv.ParseInt(el.Time, 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
		}
		o, err := strconv.ParseFloat(el.Open, 64)
		if err != nil {
			fmt.Println("Error:", err)
		}
		h, err := strconv.ParseFloat(el.High, 64)
		if err != nil {
			fmt.Println("Error:", err)
		}
		l, err := strconv.ParseFloat(el.Low, 64)
		if err != nil {
			fmt.Println("Error:", err)
		}
		c, err := strconv.ParseFloat(el.Close, 64)
		if err != nil {
			fmt.Println("Error:", err)
		}
		v, err := strconv.ParseFloat(el.Volume, 64)
		if err != nil {
			fmt.Println("Error:", err)
		}
		candles[i] = Candle{
			Time:   t,
			Open:   o,
			High:   h,
			Low:    l,
			Close:  c,
			Volume: v,
		}
	}

	// fmt.Println("\n", "Candles Length", len(candles))

	return candles
}
