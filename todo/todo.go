package todo

import (
	"fmt"
	"gogtk/common"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/gdk"
	"database/sql"
	"gogtk/db/postgres"
	"time"
  "strconv"

)

type Project struct {
	Id int64
	Title string
	Description string
	Created_at time.Time
}

type ToDo struct {
	Todo      string
	Completed bool
}

var projects []Project
var todos []ToDo
var project_new bool = false

// ---------------------------------------------------------


func GetProjects(db *sql.DB) ([]Project, error) {
	fmt.Println("\n---------------------------------------------------\n Get Projects \n---------------------------------------------------\n")
	//var projects []Project
	rows, err := db.Query("SELECT * FROM projects;")
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
		projects = append(projects, project)
		fmt.Println(" -", project)
	}
	return projects, nil
}

// ---------------------------------------------------------

func createProject(title, description string, db *sql.DB) error {
	fmt.Println("\n---------------------------------------------------\n Create Project \n---------------------------------------------------\n")

	stmt, err := db.Prepare("INSERT INTO projects (title, description, created_at) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("Error Preparing statement: %v", err)

	}
	defer stmt.Close()

	createdAt := time.Now()
	result, err := stmt.Exec(title, description, createdAt)
	if err != nil {
		return fmt.Errorf("Error getting last insert ID: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("Error getting last insert ID: %v", err)
	}

	newProject := Project{
		Id: 		id,
		Title:		title,
		Description:	description,
		Created_at:	createdAt,
	}

	projects = append(projects, newProject)

	fmt.Printf("Project created successfully. ID: %d\n", id)

	return nil

}

// css implementation
func cssWdgScnBytes(data []byte) error {

	cssProv, err := gtk.CssProviderNew()
	if err == nil {
		if err = cssProv.LoadFromData(string(data)); err == nil {
			screen, err := gdk.ScreenGetDefault()
			if err != nil {
				return err
			}
			gtk.AddProviderForScreen(screen, cssProv, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
		}
	}
	return err
}


var css = []byte(`
  #frame-white {
    background-color: #3e3e3e;
  }
  #entry, .entry {
    background-color: #ffffff;
  }
  #project-label {
    border: 1px solid #3e3e3e;
    padding: 10px;
    border-radius: 5px;
    margin: 3px;
  }
  `)

// ---------------------------------------------------------
func ToDoPage() *gtk.Box {

	fmt.Println("\n----------------------------------------------------\n Project Management Init \n----------------------------------------------------\n")



	db, err := postgres.DBConnect()
	if err != nil {
		fmt.Println("Error Connecting to Database", err)
	}
	defer db.Close()
	projects, err := GetProjects(db)
	if err != nil {
		fmt.Println("Error retrieving Projects", err)
	}
	fmt.Println("Current Projects", len(projects))

	//
	page_box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	//page_title, _ := gtk.LabelNew("Project Management")
	//
  sidebar, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	projects_box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	//projects_header, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	prj_stmt := fmt.Sprintf("Projects: %d", len(projects))
	projects_lbl, _ := gtk.LabelNew(prj_stmt)

  prj_new_btn, _ := gtk.ButtonNewWithLabel("New Project")

	sidebar.PackStart(projects_lbl, false, false, 0)
  sidebar.PackStart(prj_new_btn, false, false, 0)

  new_project_label, _ := gtk.LabelNew(strconv.FormatBool(project_new))
  sidebar.PackStart(new_project_label, false, false, 0)
  sidebar.SetName("project-label")
  sidebar.SetSizeRequest(350, 1200)

  //projects_box.SetName("frame-white")


  //cssProvider, _ := gtk.CssProviderNew()
  all_projects, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  for _, project := range projects {
    project_box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
    projectLabel, _ := gtk.LabelNew(fmt.Sprintf("Title:\n%s", project.Title))
    projectDescription, _ := gtk.LabelNew(fmt.Sprintf("Description:\n%s", project.Description))
    project_box.SetName("project-label")
    all_projects.SetName("project-label")
    projectLabel.SetSizeRequest(650, 150)
    projectLabel.SetXAlign(0)
    projectDescription.SetXAlign(0)
    project_box.PackStart(projectLabel, false, false, 0)
    project_box.PackStart(projectDescription, false, false, 0)
    all_projects.PackStart(project_box, false, false, 0)
  }
  //all_projects.SetName("frame-white")
  cssWdgScnBytes(css)


  // All projects
  projects_box.PackStart(all_projects, false, false, 0)

  // New Project box
  new_prj_box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  title_input, _ := gtk.EntryNew()
  description_input, _ := gtk.EntryNew()
  submit_new_project, _ := gtk.ButtonNewWithLabel("Submit")

  title_input.SetName("entry")
  description_input.SetName("entry")
  //

  //title_input_style, _ := title_input.GetStyleContext()
  //title_input_style.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

  //description_input_style, _ := description_input.GetStyleContext()
  //description_input_style.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

  new_prj_box.PackStart(title_input, false, false, 10)
  new_prj_box.PackStart(description_input, false, false, 10)
  new_prj_box.PackStart(submit_new_project, false, false, 10)



  submit_new_project.SetNoShowAll(true)
  title_input.SetNoShowAll(true)
  description_input.SetNoShowAll(true)

  title_input.Hide()
  description_input.Hide()
  submit_new_project.Hide()

  sidebar.PackStart(new_prj_box, false, false, 10)
  projects_box.ShowAll()

  //projects_box.SetVisible(false)

  //
  
  prj_new_btn.Connect("clicked", func () {
    project_new = !project_new
    new_project_label.SetText(strconv.FormatBool(project_new))
    fmt.Println("New Project", project_new)

    if project_new {
      prj_new_btn.Hide()
      title_input.Show()
      description_input.Show()
      submit_new_project.Show()
    } else {
      prj_new_btn.Show()
      description_input.Hide()
      title_input.Hide()
      submit_new_project.Hide()

    }
  })

  submit_new_project.Connect("clicked", func() {
    var title string
    var description string
    title, err = title_input.GetText()
    if err != nil {
      fmt.Printf("Error getting text for title_input \n%v", err)
    }
    description, err = description_input.GetText()
    if err != nil {
      fmt.Printf("Error getting text for description_input \n%v", err)
    }
    err := postgres.CreateProject(title, description)
    if err != nil {
      fmt.Printf("Error inserting new project %v", err)
    }
    project_new = !project_new
    submit_new_project.Hide()
    title_input.Hide()
    description_input.Hide()
    prj_new_btn.Show()
    fmt.Println("Submit", project_new)
  })

  //
  todos_box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  todos_stmt := fmt.Sprintf("Todos: %d", len(todos))
  todos_lbl, _ := gtk.LabelNew(todos_stmt)
  todo_title, _ := gtk.EntryNew()
  
  
  todos_box.PackStart(todos_lbl, false, false, 0)
  todos_box.PackStart(todo_title, false, false, 0)
  todos_box.SetName("project-label")
  todos_box.SetSizeRequest(700,450)

	//

	//header, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	//header.PackStart(page_title, false, false, 0)

	//page_box.PackStart(header, false, false, 0)
  page_box.PackStart(sidebar, false, false, 15)
	page_box.PackStart(projects_box, false, false, 15)
	page_box.PackStart(todos_box, false, false, 15)

	//todos_box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 8)

	// text, _ := gtk.LabelNew("Add Todos to your list")

	// New Todo
	//new_todo_btn, _ := gtk.ButtonNewWithLabel("Add ToDo")
	//todo, _ := gtk.EntryNew()

	// List ToDos
	// todoList, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	// if err != nil {
	// 	fmt.Println("Todolist Error", err)
	// }

	// ----------------------------------------------------------------
	// Todo Column Container
	// todoColumns, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 20)
	// if err != nil {
	// 	fmt.Println("Error creating Todo Columns", err)
	// }
	// // Container for Incomplete Todos
	// completedTodos, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	// if err != nil {
	// 	fmt.Println("Error creating Incomplete Column", err)
	// }
	//
	// // Container for Completed Todos
	// incompleteTodos, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	// if err != nil {
	// 	fmt.Println("Error creating Completed Column", err)
	// }
	//
	// ----------------------------------------------------------------

	// new_todo_btn.Connect("clicked", func() {
	// 	todoText, err := todo.GetText()
	// 	if err != nil {
	// 		fmt.Println("Error GetText for todo", err)
	// 	}
	// 	fmt.Println("Add ToDo", todoText)
	// 	todos = append(todos, ToDo{Todo: todoText, Completed: false})
	// 	updateTodoList(incompleteTodos, completedTodos)
	// 	fmt.Println("\n-----------------------\n")
	// 	todoList.ShowAll() // Render the list
	// })
	//
	// ----------------------------------------

	//header.PackStart(new_todo_btn, false, false, 0)
	//header.PackStart(todo, false, false, 0)
	//page_box.PackStart(new_todo_btn, false, false, 0)
	// page_box.PackStart(todoColumns, false, false, 0)
	// // --------
	// todoColumns.PackStart(incompleteTodos, true, true, 0)
	// todoColumns.PackStart(completedTodos, true, true, 0)
	//
	// updateTodoList(incompleteTodos, completedTodos) // Initial updade to display todos
	return page_box
}



func redrawProjectsPage(projectsBox *gtk.Box) {
    // Clear existing content
    children := projectsBox.GetChildren()
    children.Foreach(func(item interface{}) {
      widget := item.(*gtk.Widget)
      projectsBox.Remove(widget)
    })

    // Fetch updated projects from the database
    db, err := postgres.DBConnect()
    if err != nil {
        fmt.Println("Error connecting to database:", err)
        return
    }
    defer db.Close()

    // projects, err := postgres.GetProjects(db)
    // if err != nil {
    //     fmt.Println("Error retrieving projects:", err)
    //     return
    // }
    //
    // // Recreate project list
    // for _, project := range projects {
    //     projectLabel, _ := gtk.LabelNew(fmt.Sprintf("%s: %s", project.Title, project.Description))
    //     projectsBox.PackStart(projectLabel, false, false, 0)
    // }

    // Show all new widgets
    projectsBox.ShowAll()
}

func updateTodoList(incompleteTodos, completedTodos *gtk.Box) {
	// Clear existing todos from both containers
	common.ClearContainer(incompleteTodos)
	common.ClearContainer(completedTodos)

	// Create Labels for containers
	completedLbl, _ := gtk.LabelNew("Completed")
	incompleteLbl, _ := gtk.LabelNew("Incomplete")

	completedTodos.PackStart(completedLbl, false, false, 0)
	incompleteTodos.PackStart(incompleteLbl, false, false, 0)

	// ------------------------------

	// Add updated list of todos to the appropriate container based on completed status
	for i, todo := range todos {
		todoLbl, err := gtk.LabelNew(todo.Todo)
		if err != nil {
			fmt.Println("Error creating Todo label", err)
			continue
		}
		todoBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 3)
		if err != nil {
			fmt.Println("Error creating todoBox", err)
		}
		complete_btn, err := gtk.ButtonNewWithLabel("X")
		if err != nil {
			fmt.Println("Error creating Completed button", err)
		}
		if todo.Completed {
			todoBox.PackStart(complete_btn, false, false, 5)
			todoBox.PackStart(todoLbl, false, false, 0)
			completedTodos.PackStart(todoBox, false, false, 0)
		} else {
			todoBox.PackStart(complete_btn, false, false, 5)
			todoBox.PackStart(todoLbl, false, false, 0)
			incompleteTodos.PackStart(todoBox, false, false, 0)
		}
		complete_btn.Connect("clicked", func() {
			todos[i].Completed = !todos[i].Completed
			fmt.Println("X Clicked", todos[i].Todo, todos[i].Completed)
			updateTodoList(incompleteTodos, completedTodos)
		})
	}

	incompleteTodos.ShowAll()
	completedTodos.ShowAll()
}
