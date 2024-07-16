package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 31337
    user     = "postgres"
    password = "1234567" // Replace with actual password
    dbname   = "myapp_db"
)

func createDatabase() error {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
        host, port, user, password)
    
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
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

func createTables(db *sql.DB) error {
    // Create account table
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS account (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            password VARCHAR(100) NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        return fmt.Errorf("error creating account table: %v", err)
    }

    // Create todo table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS todo (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES account(id),
            title VARCHAR(100) NOT NULL,
            description TEXT,
            completed BOOLEAN DEFAULT FALSE,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        return fmt.Errorf("error creating todo table: %v", err)
    }

    // Create coinbase_btcusd_5m table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS coinbase_btcusd_5m (
            id SERIAL PRIMARY KEY,
            timestamp TIMESTAMP NOT NULL,
            open DECIMAL(10, 2) NOT NULL,
            high DECIMAL(10, 2) NOT NULL,
            low DECIMAL(10, 2) NOT NULL,
            close DECIMAL(10, 2) NOT NULL,
            volume DECIMAL(14, 6) NOT NULL
        )
    `)
    if err != nil {
        return fmt.Errorf("error creating coinbase_btcusd_5m table: %v", err)
    }

    fmt.Println("All tables created successfully")
    return nil
}

func main() {
    err := createDatabase()
    if err != nil {
        log.Fatal("Error creating database:", err)
    }

    db, err := connectDB()
    if err != nil {
        log.Fatal("Error connecting to the database:", err)
    }
    defer db.Close()

    err = createTables(db)
    if err != nil {
        log.Fatal("Error creating tables:", err)
    }
}

