package coinbase

import (
	"fmt"
	"time"
	"sync"
	"os"
	"strconv"
	"net/http"
	"encoding/json"
	_"github.com/joho/godotenv"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GetCandles(symbol, tf string) () {
} 

type Candle struct {
	Time int64	`json:"time"`
	Open float64	`json:"open"`
	High float64	`json:"high"`
	Low float64	`json:"low"`
	Close float64	`json:"close"`
	Volume float64	`json:"volume"`
}

type CoinbaseCandle struct {
	Time   string `json:"start"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

type Timeframe struct {
	Label 	string
	Xch	string
	Tf	int
}

type CandlesData struct {
	Candles []Candle
	Symbol	string
	TF	Timeframe
}


type ApiResponse struct {
	Candles []CoinbaseCandle `json:"candles"`
}

func DoCoinbase(full bool) {
	fmt.Println("=======================================")
	fmt.Println("\n", "Do Coinbase")
	cbKey := os.Getenv("CBAPIKEY")
	cbSecret := os.Getenv("CBAPISECRET")

	//
	//accounts, folio_val := getCBAccts()
	//writeAccountsToDB(accounts, "coinbase")
	//cb_AddFolioVal(folio_val)

	//orders := getCBOrders()
	// for order := range orders {
	// 	fmt.Println("Order", orders[order])
	// }
	//cb_WriteOrders(orders)

	//fills := getCBFills()

	// for fill := range fills {
	// 	fmt.Println("Fills", fills[fill].ProductID, fills[fill].Size, fills[fill].Price)
	// }

	//cb_WriteFills(fills)
	// fmt.Println("ENVIRONMENT VARIABLE", cbKey, cbSecret)

	type ExchangeSchema struct {
		Coinbase struct {
			Watchlist []string    `json:"watchlist"`
			TF        []Timeframe `json:"tf"`
		} `json:"coinbase"`
	}
	var ExchangeSchemaObj = ExchangeSchema{
		Coinbase: struct {
			Watchlist []string    `json:"watchlist"`
			TF        []Timeframe `json:"tf"`
		}{
			Watchlist: []string{
				"BTC-USD",
				//"XRP-USD",
				//"ETH-USD",
				//"XLM-USD",
				////"LTC-USD",
				//"SOL-USD",
				//"ADA-USD",
				//"DOGE-USD",
				//"SHIB-USD",
			},
			TF: []Timeframe{
				{
					Label: "1m",
					Tf:    1,
					Xch:   "ONE_MINUTE",
				},
				// {
				// 	Label: "5m",
				// 	Tf:    5,
				// 	Xch:   "FIVE_MINUTE",
				// },
				// {
				// 	Label: "15m",
				// 	Tf:    15,
				// 	Xch:   "FIFTEEN_MINUTE",
				// },
				// {
				// 	Label: "30m",
				// 	Tf:    30,
				// 	Xch:   "THIRTY_MINUTE",
				// },
				// {
				// 	Label: "1H",
				// 	Tf:    60,
				// 	Xch:   "ONE_HOUR",
				// },
				// {
				// 	Label: "2H",
				// 	Tf:    120,
				// 	Xch:   "TWO_HOUR",
				// },
				// {
				// 	Label: "6H",
				// 	Tf:    360,
				// 	Xch:   "SIX_HOUR",
				// },
				// {
				// 	Label: "1D",
				// 	Tf:    1440,
				// 	Xch:   "ONE_DAY",
				// },
			},
		},
	}

	startT := time.Now()

	// Slice to hold all candles data
	var allCandlesData []CandlesData

	// Loop through each symbol and timeframe synchronously
	for _, symbol := range ExchangeSchemaObj.Coinbase.Watchlist {
		for _, tf := range ExchangeSchemaObj.Coinbase.TF {
			candles := doCBRequest(symbol, tf, cbKey, cbSecret, full)
			allCandlesData = append(allCandlesData, CandlesData{Candles: candles, Symbol: symbol, TF: tf})
		}
	}

	fmt.Println("Writing To Database")

	// Use goroutines to write candles to the database in parallel
	var wg sync.WaitGroup
	for _, data := range allCandlesData {
		wg.Add(1)
		go func(d CandlesData) {
			defer wg.Done()
			//writeCandles(d.Candles, "coinbase", d.Symbol, d.TF)
		}(data)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("All tasks completed.")

	endT := time.Now()
	duration := endT.Sub(startT)
	fmt.Println("DURATION", duration)
}

func getCBSign(apiSecret string, timestamp int64, method string, path string, body string) string {
	fmt.Println("\n-------------------------\n getCBSign \n-------------------------\n")
	fmt.Println("API:", apiSecret)
	message := fmt.Sprintf("%d%s%s%s", timestamp, method, path, body)
	fmt.Println("Message\n", message)
	hasher := hmac.New(sha256.New, []byte(apiSecret))
	hasher.Write([]byte(message))
	signature := hex.EncodeToString(hasher.Sum(nil))
	fmt.Println("Signature\n", signature)
	return signature
}

func doCBRequest(symbol string, tf Timeframe, apiKey string, apiSecret string, full bool) []Candle {
	fmt.Println("\n-------------------------\n doCBRequest \n-------------------------\n")
	numCandles := 300
	totalDuration := time.Duration(tf.Tf*numCandles) * time.Minute
	endTime := time.Now()
	startTime := endTime.Add(-totalDuration).Unix()

	var candles []Candle
	if full {
		//candles = getFullCBCandlesAuth(apiKey, apiSecret, symbol, tf.Tf, tf.Xch, startTime, endTime.Unix())
	} else {
		candles = getCBCandlesAuth(apiKey, apiSecret, symbol, tf.Tf, tf.Xch, startTime, endTime.Unix())
	}

	return candles
}


func getCBCandlesAuth(apiKey string, apiSecret string, symbol string, tf int, tflbl string, start int64, end int64) []Candle {
	fmt.Println("\n-------------------------\n getCBCandlesAuth \n", symbol, tf, "\n-------------------------\n")
	fmt.Println("AUTH Coinbase Candles", symbol, tf, start, end)

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





