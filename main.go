package main

import (
	"fmt"
	"net/http"
  "html/template"
	"hbw/views"
  "hbw/db"
  "time"
  "strconv"
  "strings"
  "sync"
)

var index *views.View
var contact *views.View
var projects_page *views.View
var finance_page *views.View
var chart *views.View

var count int = 0
type Project struct {
	Id int
	Title string
	Description string
	Created_at time.Time
  Selected bool
}

func GetAllProjects(w http.ResponseWriter) {
  projects, err := db.GetProjects()
  if err != nil {
    fmt.Println("Error getting projects from DB", err)
  }
  t, _ := template.ParseFiles("views/projects.html")
  err = t.Execute(w, projects)
  if err != nil {
    fmt.Println("Error executing projects template with projects", err)
  }
}

var (
  chartCount  int
  mu          sync.Mutex
)

func updateChartCount() {
  mu.Lock()
  chartCount++

  mu.Unlock()

}
func main() {
	fmt.Println("Starting Server on port 3000")

  ticker := time.NewTicker(1 * time.Second)
  go func() {
    for range ticker.C {
      updateChartCount()
    }
  }()
	
	index = views.NewView("bootstrap", "views/index.html")
  finance_page = views.NewView("bootstrap", "views/finance.html")
	projects_page = views.NewView("bootstrap", "views/projects.html")
	contact = views.NewView("bootstrap", "views/contact.html")

  // chart = views.NewView("bootstrap", "views/chart.html")

  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

  // Views || Pages
	http.HandleFunc("/", indexHandler)
  http.HandleFunc("/finance", financeHandler)
  http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/contact", contactHandler)
  http.HandleFunc("/chart", chartHandler)

  // Function endpoints
	http.HandleFunc("/createProject", create_project_handler)
	http.HandleFunc("/deleteProject", delete_project_handler)
	http.HandleFunc("/selectProject", select_project_handler)
  //
	http.HandleFunc("/newTodo", create_todo_handler)
	http.HandleFunc("/toggleTodoCompleted", toggle_todo_handler)
	http.HandleFunc("/deleteTodo", delete_todo_handler)

	http.ListenAndServe(":3000", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n-----------------------\n Home Page \n-----------------------\n")
	index.Render(w, nil)
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n-----------------------\n Projects Page \n-----------------------\n")
  all_projects, err := db.GetProjects()
  if err != nil {
    fmt.Println("Error Getting All Projects from DB", err)
  }
	projects_page.Render(w, all_projects)
}

func create_project_handler(w http.ResponseWriter, r *http.Request) {
  title := r.FormValue("title")
  description := r.FormValue("description")
  fmt.Println("Create Project Handler", title, description)
  err := db.CreateProject(title, description)
  if err != nil {
    fmt.Println("create_projects_handler failed to add project to db", err)
  }
  http.Redirect(w, r, "/projects", http.StatusSeeOther)
}

func create_todo_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n-----------------------\n create_todo_handler \n-----------------------\n")
  todo := r.FormValue("todo")
  project_idstr := r.FormValue("project_id")
  project_id, err := strconv.Atoi(project_idstr)
  if err != nil {
    fmt.Println("Error converting str to int", err)
  }

  fmt.Println(todo, project_id)

  err = db.CreateTodo(todo, project_id)
  if err != nil {
    fmt.Println("Error Creating Todo\n", err)
  } 
  http.Redirect(w, r, "/projects", http.StatusSeeOther)
}

func toggle_todo_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n-----------------------\n toggle_todo_handler \n-----------------------\n")
  err := r.ParseForm()
  if err != nil {
    fmt.Println("Error Parsing form toggle_todo_handler", err)
  }
  idstr := r.FormValue("id")
  // id, err := strconv.ParseInt(idstr, 10, 64)
  // if err != nil {
  //   fmt.Println("Error parsing int (toggle_todo_handler")
  // }

  boolstr := r.FormValue("completed")
  completed, err := strconv.ParseBool(boolstr)
  if err != nil {
    fmt.Println("Error parsing bool (toggle_todo_handler")
  }
  fmt.Println(idstr, boolstr)

  db.ToggleTodoCompleted(idstr, completed)

  http.Redirect(w, r, "/projects", http.StatusSeeOther)
}

func delete_todo_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n-----------------------\n delete_todo_handler \n-----------------------\n")
  idstr := r.FormValue("id")
  id, err := strconv.ParseInt(idstr, 10, 64)
  if err != nil {
    fmt.Println("Error converting string to int", err)
  }
  fmt.Println(id)

  db.DeleteTodo(id)
  http.Redirect(w, r, "/projects", http.StatusSeeOther)
}
// ---------------------

func delete_project_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n-----------------------\n delete_project_handler \n-----------------------\n")
  idstr := r.FormValue("id")
  id, err := strconv.ParseInt(idstr, 10, 64)
  if err != nil {
    http.Error(w, "Invalid project ID", http.StatusBadRequest)
    return
  }

  err = db.DeleteProject(id)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return 
  }
  http.Redirect(w, r, "/projects", http.StatusSeeOther)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n-----------------------\n Contact Page \n-----------------------\n")
	contact.Render(w, nil)
}


func select_project_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n-----------------------\n select_project_handler \n-----------------------\n")
  idstr := r.FormValue("id")
  id, err := strconv.ParseInt(idstr, 10, 64)
  if err != nil {
    fmt.Println("Error converting id to string (select_project_handler)", err)
  }
  selstr := r.FormValue("selected")
  sel, err := strconv.ParseBool(selstr)
  fmt.Println("Selected ID", id, sel)

  //db.UpdateProject(id, sel)
  http.Redirect(w, r, "/projects", http.StatusSeeOther)


}



// Finance Chart
var data = []int{10, 20, 30, 40, 50}

func generateSVG(data []int) string {
	var svg strings.Builder

	// SVG header
	svg.WriteString(`<svg width="300" height="200" xmlns="http://www.w3.org/2000/svg">`)

	// Bar properties
	barWidth := 40
	barSpacing := 10
	maxHeight := 100

	for i, value := range data {
		height := (value * maxHeight) / 50 // Scale height based on max value
		x := i * (barWidth + barSpacing)
		y := maxHeight - height

		// Create a rectangle for each bar
		svg.WriteString(fmt.Sprintf(
			`<rect x="%d" y="%d" width="%d" height="%d" fill="blue"/>`,
			x, y, barWidth, height,
		))
	}

	// SVG footer
	svg.WriteString(`</svg>`)

	return svg.String()
}


func financeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n-----------------------\n Finance Page \n-----------------------\n")
  count++
	finance_page.Render(w, count)
  // chartHandler(w, r)
}


type PageData struct {
  Tick int
}
func chartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n-----------------------\n Chart \n-----------------------\n")
  //chart.Render(w, generateSVG(data))

  // mu.Lock()
  // currentCount := chartCount
  // mu.Unlock()

  // data := PageData{
  //   Tick: currentCount,
  // }
  chart.Render(w, nil)
}

// ---------------------------------------------------------------------------------------------------
// Node Graph 
// ---------------------------------------------------------------------------------------------------

type Node struct {
  ID            int
  Value         string
  Children      []Node
}

func GenerateTree() {

}






