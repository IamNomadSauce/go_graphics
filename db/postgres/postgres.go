package postgres

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "os"
    "github.com/joho/godotenv"
    "strconv"
    "time"
)

type Candle struct {
    Time   int64
    Open   float64
    High   float64
    Low    float64
    Close  float64
    Volume float64
}

type Timeframe struct {
    Label string
    Xch   string
    Tf    int
}


type Project struct {
	Id int64
	Title string
	Description string
	Created_at time.Time
}

var host string
var port int
var user string
var password string
var dbname string

func DBConnect() (*sql.DB, error) {
	
	fmt.Println("\n------------------------------\n DBConnect \n------------------------------\n")
    host = os.Getenv("PG_HOST")
    portStr := os.Getenv("PG_PORT")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        fmt.Printf("Invalid port number: %v\n", err)
        return nil, err
    }
    user = os.Getenv("PG_USER")
    password = os.Getenv("PG_PASS")
    dbname = os.Getenv("PG_DBNAME")

    // Connect to the default 'postgres' database to check for the existence of the target database
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        fmt.Println("Error opening Postgres", err)
        return nil, err
    }
    //defer db.Close()

    return db, nil

}

func CreateDatabase() (*sql.DB, error) {
    fmt.Println("\n------------------------------\n CreateDatabase \n------------------------------\n")

    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
        return nil, err
    }
    host = os.Getenv("PG_HOST")
    portStr := os.Getenv("PG_PORT")
    port, err = strconv.Atoi(portStr)
    if err != nil {
        fmt.Printf("Invalid port number: %v\n", err)
        return nil, err
    }
    user = os.Getenv("PG_USER")
    password = os.Getenv("PG_PASS")
    dbname = os.Getenv("PG_DBNAME")

    // Connect to the default 'postgres' database to check for the existence of the target database
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        fmt.Println("Error opening Postgres", err)
        return nil, err
    }
    defer db.Close()

    // Check if the database already exists
    var exists bool
    query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '%s')", dbname)
    err = db.QueryRow(query).Scan(&exists)
    if err != nil {
        fmt.Println("Error checking database existence", err)
        return nil, err
    }

    if exists {
        fmt.Printf("Database %s already exists\n", dbname)

	    psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	    newDB, err := sql.Open("postgres", psqlInfo)
	    if err != nil {
		    return nil, fmt.Errorf("Error connecting to existing database: %v", err)
	    }
	    return newDB, nil
    }

    // Create the database if it does not exist
    _, err = db.Exec("CREATE DATABASE " + dbname)
    if err != nil {
        fmt.Println("Error creating database", err)
        return nil, err
    }

    fmt.Printf("Database %s created successfully\n", dbname)


    newDB, err := sql.Open("postgres", psqlInfo)
    if err != nil {
	    return nil, fmt.Errorf("Error connecting to new database: %v\n", err)
    }
    // Create Tables 
    //err = CreateTables(db)
    //if err != nil {
	    //return fmt.Errorf("Error creating tables")
    //}
    return newDB, nil
}

func ShowDatabases(db *sql.DB) error {
	fmt.Println("\n------------------------------\n ShowDatabases \n------------------------------\n")
	rows, err := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false")
	if err != nil {
		fmt.Println("Error listing Databases", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var datname string
		if err := rows.Scan(&datname); err != nil {
			fmt.Println("Error scanning database name", err)
			return err
		}
		fmt.Println(" -", datname)
	}

	return nil
}

func CreateTables(db *sql.DB) error {
	fmt.Println("\n------------------------------\n CreatTables \n------------------------------\n")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("Error creating Projects table")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			project_id INTEGER REFERENCES projects(id),
			title VARCHAR(100) NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("Error creating TODOs table")
	}

	fmt.Println("All Tables Created Successfully")

	return nil
}

func ListTables(db *sql.DB) error {
	fmt.Println("\n------------------------------\n ListTables \n------------------------------\n")
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		fmt.Println("Error listing tables", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil{
			fmt.Println("Error scanning table name", err)
			return err
		}
		fmt.Println(" -", tableName)
	}

	return nil
}

func GetProjects(db *sql.DB) ([]Project, error) {
	fmt.Println("\n---------------------------------------------------\n Get Projects \n---------------------------------------------------\n")

	var projects []Project
	rows, err := db.Query("SELECT * FROM projects;")
	if err != nil {
		fmt.Println("Error listing projects", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var project Project
		if err := rows.Scan(project); err != nil {
			fmt.Println("Error scanning Projects table", err)
			return nil, err
		}
		projects = append(projects, project)
		fmt.Println(" -", project)
	}
	return projects, nil
}
















