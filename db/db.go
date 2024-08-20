package db

import (
    "database/sql"
    "fmt"
    "github.com/lib/pq"
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

type Node struct {
  ID int
  Value string
  Children []Node
}


var host string
var port int
var user string
var password string
var dbname string

func DBConnect() (*sql.DB, error) {
	
	fmt.Println("\n------------------------------\n DBConnect \n------------------------------\n")
  err := godotenv.Load()
  if err != nil {
    fmt.Printf("Error loading .env file %v\n", err)

  }
    host = os.Getenv("PG_HOST")
    portStr := os.Getenv("PG_PORT")
    // fmt.Printf("Host:\n%s\nPort:\n%d\nUser:\n%s\nPW:\n%s\nDB:\n%s\n", host, port, user, password, dbname)
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
    err = CreateTables(db)
    if err != nil {
	    return nil, fmt.Errorf("Error creating tables")
    }
    return newDB, nil
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
		CREATE TABLE IF NOT EXISTS nodes (
			id SERIAL PRIMARY KEY,
      value VARCHAR(255) NOT NULL,
      parent_id INTEGER REFERENCES nodes(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("Error creating Nodes table")
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			title VARCHAR(100) NOT NULL,
			description TEXT,
      selected bool,
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
      children INTEGER[],
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("Error creating TODOs table")
	}

  fmt.Println("All Tables Created successfully")

	return nil
}

type Project struct {
	Id int64
	Title string
	Description string
	Created_at time.Time
  Selected  bool
  Todos []Todo
}

func GetProjects() ([]Project, error) {
    fmt.Println("\n---------------------------------------------------\n Get Projects \n---------------------------------------------------\n")

    db, err := DBConnect()
    if err != nil {
        fmt.Println("Error connecting to DB (GetProjects)", err)
    }

    var projects []Project
    rows, err := db.Query("SELECT id, title, description, created_at FROM projects;")
    if err != nil {
        fmt.Println("Error listing projects", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var project Project
        if err := rows.Scan(&project.Id, &project.Title, &project.Description, &project.Created_at); err != nil {
            fmt.Println("Error scanning Projects table", err)
            return nil, err
        }
        project.Selected = false
        todoRows, err := db.Query("SELECT id, project_id, title, description, completed, children, created_at FROM todos WHERE project_id = $1;", project.Id)
        if err != nil {
          fmt.Println("Error listing todos from db", err)
        }
        defer todoRows.Close()

        var todos []Todo
        for todoRows.Next() {
          var todo Todo
          if err := todoRows.Scan(&todo.Id, &todo.Project_id, &todo.Title, &todo.Description, &todo.Completed, &todo.Children, &todo.Created_at); err != nil {
              fmt.Println("Error scanning Todos table", err)
              return nil, err
          }
          todos = append(todos, todo)
        }
        project.Todos = todos
        projects = append(projects, project)
        // fmt.Println(" -", project)
    }
    fmt.Println(len(projects), "projects found\n")
    return projects, nil
}

type Todo struct {
  Id int64
  Project_id int64
  Title string
  Description string
  Completed bool
  Children pq.Int64Array 
  Created_at time.Time
}

func GetTodos() ([]Todo, error) {
	fmt.Println("\n---------------------------------------------------\n GetTodos \n---------------------------------------------------\n")

  db, err := DBConnect()

  todoRows, err := db.Query("SELECT id, project_id, title, description, completed, children, created_at FROM todos ;")
  if err != nil {
    fmt.Println("Error listing todos from db", err)
  }
  defer todoRows.Close()

  var todos []Todo
  for todoRows.Next() {
    var todo Todo
    if err := todoRows.Scan(&todo.Id, &todo.Project_id, &todo.Title, &todo.Description, &todo.Completed, &todo.Children, &todo.Created_at); err != nil {
        fmt.Println("Error scanning Todos table", err)
        return nil, err
    }
    todos = append(todos, todo)
  }
  return todos, nil

}

func CreateTodo(title string, projectId int, parent_id int) error {
	fmt.Println("\n---------------------------------------------------\n CreateTodo \n---------------------------------------------------\n")
	fmt.Printf("Title: %v\nProject ID: %v\n", title, projectId)
  
  description := ""
  completed := false
  var children []int

  db, err := DBConnect()
  if err != nil {
    fmt.Printf("Error Connecting to DB %v", err)
  }
  defer db.Close()
  sqlStatement := `
		INSERT INTO todos (project_id, parent_id, title, description, completed, children)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	var newTodoID int

	// Execute the statement and get the new todo ID
	err = db.QueryRow(sqlStatement, projectId, parent_id, title, description, completed, pq.Array(children)).Scan(&newTodoID)
	if err != nil {
		return fmt.Errorf("Error inserting todo: %v", err)
	}

  err = Update_todo_Children(parent_id, newTodoID)
  if err != nil {
    return fmt.Errorf("Error updating todo children: %v", err)
  }

	return nil
}

func Update_todo_Children(parent_id, child_id int) error{
 	fmt.Println("\n---------------------------------------------------\n UpdateTodo \n---------------------------------------------------\n")
  fmt.Printf("Todo_id:\n%v\nChild_id:\n%v\n", parent_id, child_id)
  // fmt.Printf("Todo:\n%v\nProject_id:\n%v\nParent:\n%\n", todo, projectId, parent_id)
  db, err := DBConnect()
	if err != nil {
		return fmt.Errorf("Error connecting to DB: %v", err)
	}
	defer db.Close()

	// Use array_append to add the child_id to the children array
	query := `
		UPDATE todos
		SET children = array_append(children, $1)
		WHERE id = $2
	`
	_, err = db.Exec(query, child_id, parent_id)
	if err != nil {
		return fmt.Errorf("Error updating children array: %v", err)
	}

  return nil
}

// func GetTodo(id) error {
// 	fmt.Println("\n---------------------------------------------------\n GetTodo \n---------------------------------------------------\n")
//   fmt.Printf("%v\n", id)
//
//   db, err := DBConnect()
//   if err != nil {
//     fmt.Println("Error connecting to db GetTodo", err)
//   }
//   defer db.Close()
//
//   rows, err := db.Query("SELECT id, title, description, project_id, completed, created_at FROM todos where id = $1;")
//   if err != nil {
//       fmt.Println("Error listing projects", err)
//       return nil, err
//   }
//   defer rows.Close()
// }

func ToggleTodoCompleted(id string, completed bool) error {
	fmt.Println("\n---------------------------------------------------\n ToggleTodoCompleted \n---------------------------------------------------\n")
  fmt.Printf("%v\n%v", id, completed)

  db, err := DBConnect()
    if err != nil {
        return fmt.Errorf("Error connecting to DB: %v", err)
    }
    defer db.Close()

    completed = !completed
    query := "UPDATE todos SET completed = $1 WHERE id = $2"
    _, err = db.Exec(query, completed, id)
    if err != nil {
        return fmt.Errorf("Error updating todo: %v", err)
    }

    return nil
}

func CreateProject(title, description string) error {
	fmt.Println("\n---------------------------------------------------\n CreateProject \n---------------------------------------------------\n")
  fmt.Printf("\n Title: %s \n Description: %s", title, description)
  
  db, err := DBConnect()
  if err != nil {
    fmt.Printf("Error Connecting to DB %v", err)
  }
  defer db.Close()

  sqlStatement := `
    INSERT INTO projects (title, description, created_at)
    VALUES ($1, $2, $3)
    RETURNING id`

  var id int64
  err = db.QueryRow(sqlStatement, title, description, time.Now()).Scan(&id)
  if err != nil {
    return fmt.Errorf("Error inserting new project into database: \n %v", err)
  }
  fmt.Printf("New project created successfully with ID: %v\n\n", id)
  return nil
}

func DeleteProject(id int64) error {
	fmt.Println("\n---------------------------------------------------\n DeleteProject \n---------------------------------------------------\n")
  db, err := DBConnect()
  if err != nil {
    fmt.Println("Error connecting to DB: %v", err)
    return err 
  }
  defer db.Close()

  _, err = db.Exec("DELETE FROM projects where id = $1", id)
  if err != nil {
    fmt.Println("Error deleting project: %v", err)
    return err
  }

  fmt.Printf("Project with ID %d deleted successfully\n", id)
  return nil
}


func DeleteTodo(id int64) error {
	fmt.Println("\n---------------------------------------------------\n DeleteTodo \n---------------------------------------------------\n")
  db, err := DBConnect()
  if err != nil {
    fmt.Println("Error connecting to DB: %v", err)
    return err 
  }
  defer db.Close()

  _, err = db.Exec("DELETE FROM todos where id = $1", id)
  if err != nil {
    fmt.Println("Error deleting todo: %v", err)
    return err
  }

  fmt.Printf("Todo with ID %d deleted successfully\n", id)
  return nil
}

func UpdateProject(id int64, sel bool) error {
	fmt.Println("\n---------------------------------------------------\n UpdateProject \n---------------------------------------------------\n")
  db, err := DBConnect()
  if err != nil {
    fmt.Println("Error Connecting to DB", err)
  }
  query, err := db.Prepare("UPDATE projects SET (selected) = (?) where id = ?")
  defer db.Close()
  
  if err != nil {
    fmt.Println("Error updating project", err)
    return err
  }
  _, err = query.Exec(sel, id)
  fmt.Println("Project updated")
  return nil
}

// --------------------------------------------------------------------------
// Nodes
// --------------------------------------------------------------------------

func GetNodes() ([]Node, error) {
	fmt.Println("\n---------------------------------------------------\n GetNodes \n---------------------------------------------------\n")

  db, err := DBConnect()
  if err != nil {
    return nil, fmt.Errorf("Error connecting to db: %v", err)
  }
  rows, err := db.Query("SELECT id, value, parent_id FROM nodes")
  if err != nil {
    return nil, fmt.Errorf("Error fetching nodes: %v", err)
  }
  defer rows.Close()

  nodeMap := make(map[int]*Node)
  var rootNodes []Node

  for rows.Next() {
    var id, parentID sql.NullInt64
    var value string
    if err := rows.Scan(&id, &value, &parentID); err != nil {
      return nil, fmt.Errorf("Error scanning node: %v", err)
    }

    node := Node{
      ID: int(id.Int64),
      Value: value,
    }

    nodeMap[node.ID] = &node

    if parentID.Valid {
      parentNode := nodeMap[int(parentID.Int64)]
      parentNode.Children = append(parentNode.Children, node)
    } else {
      rootNodes = append(rootNodes, node)
    }
  }
  fmt.Println("Nodes:", len(rootNodes))
  return rootNodes, nil
}



