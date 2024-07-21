package postgres

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "os"
    "github.com/joho/godotenv"
    "strings"
)



type Candle struct {
	Time int64	
	Open float64
	High float64
	Low float64
	Close float64
	Volume float64
}

type Timeframe struct {
	Label 	string
	Xch 	string
	Tf	int
}

var host string
var port string
var user string
var password string
var dbname string

func CreateDatabase() error {
	fmt.Println("\n------------------------------\n Create Postgres Database \n------------------------------\n")
	
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	host     := os.Getenv("PG_HOST")
	port     := os.Getenv("PG_PORT")
	user     := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASS") 
	dbname   := os.Getenv("PG_DBNAME")

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
        host, port, user, password)

	fmt.Printf("host: %s\n port: %s\n user: %s\n password: %s\n dbname: %s\n", host, port, user, password, dbname)
    
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
	    fmt.Println("Error opening Postgres", err)
        return err
    }
    defer db.Close()

    _, err = db.Exec("CREATE DATABASE " + dbname)
    if err != nil {
        // Check if the database already exists
        if err.Error() == fmt.Sprintf("pq: database \"%s\" already exists", dbname) {
            fmt.Printf("Database %s already exists\n", dbname)
            return nil
        }
        return err
    }

    fmt.Printf("Database %s created successfully\n", dbname)
    return nil
}

func connectDB() (*sql.DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }
    
    err = db.Ping()
    if err != nil {
        return nil, err
    }
    
    fmt.Println("Successfully connected to the database")
    return db, nil
}


// -------------------------------
// Coinbase
// -------------------------------

// WriteCBCandles
func writeCandles(candles []Candle, exchange string, raw_symbol string, tf Timeframe) error {
    symbol := strings.ReplaceAll(raw_symbol, "-", "")
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    defer db.Close()

    createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s_%s_%d (
        time BIGINT NOT NULL PRIMARY KEY,
        open FLOAT NOT NULL,
        high FLOAT NOT NULL,
        low FLOAT NOT NULL,
        close FLOAT NOT NULL,
        volume FLOAT NOT NULL
    )`, exchange, symbol, tf.Tf)

    _, err = db.Exec(createTableQuery)
    if err != nil {
        return fmt.Errorf("failed to create table: %w", err)
    }

    insertQuery := fmt.Sprintf(`
    INSERT INTO %s_%s_%d (time, open, high, low, close, volume)
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (time) DO UPDATE SET
        open = EXCLUDED.open,
        high = EXCLUDED.high,
        low = EXCLUDED.low,
        close = EXCLUDED.close,
        volume = EXCLUDED.volume
    `, exchange, symbol, tf.Tf)

    for _, candle := range candles {
        _, err = db.Exec(insertQuery, candle.Time, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
        if err != nil {
            return fmt.Errorf("failed to insert candle: %w", err)
        }
    }

    return nil
}

// GetCBCandles
func getCandles(exchange string, symbol string, tf Timeframe, limit int) ([]Candle, error) {
    // Replace any non-alphanumeric characters in the symbol with underscores
    safeSymbol := strings.Map(func(r rune) rune {
        if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
            return r
        }
        return '_'
    }, symbol)

    tableName := fmt.Sprintf("%s_%s_%d", exchange, safeSymbol, tf.Tf)

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    defer db.Close()

    query := fmt.Sprintf(`
        SELECT time, open, high, low, close, volume
        FROM %s
        ORDER BY time DESC
        LIMIT $1
    `, tableName)

    rows, err := db.Query(query, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %w", err)
    }
    defer rows.Close()

    var candles []Candle
    for rows.Next() {
        var c Candle
        err := rows.Scan(&c.Time, &c.Open, &c.High, &c.Low, &c.Close, &c.Volume)
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }
        candles = append(candles, c)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    // Reverse the slice to get ascending order by time
    for i, j := 0, len(candles)-1; i < j; i, j = i+1, j-1 {
        candles[i], candles[j] = candles[j], candles[i]
    }

    return candles, nil
}


// To be done later...
// --------------------
// get portfolio

// get fills








