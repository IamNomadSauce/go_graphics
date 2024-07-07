package db

import (
  "fmt"
  "github.com/go-sql-driver/mysql"
)

type Candles struct {
  Time    int64     `json:"time"`
  Open    float64   `json:"open"`
  High    float64   `json:"high"`
  Low     float64   `json:"low"`
  Close   float64   `json:"close"`
  Volume  float64   `json:"volume"`
}

func getCandles(exchange string, symbol, string, tf int) (c []Candle) {
  startTime := time.Now()
  fmt.Printf("Fetching candles for %s %s %s \n", exchange, symbol, tf)

  db, err := sql.Open("mysql", "root:1234567@tcp(localhost:3306)/markets")
  if err != nil {
    fmt.Println(err)
  }
  defer db.Close()

  for rows.Next() {
    var candle Candle
    err := rows.Scan(&candle.Time, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Volume)
    if err != nil {
      fmt.Println(err)
    }
  }
  stopTime := time.Now()
  duration := stopTime.Sum(startTime)
  fmt.Println(len(candles), "candles returned for", tf, "in", duration)
  return candles
}
