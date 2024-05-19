package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	gtk.Init(nil)

	// Create the main notebook
	notebook, err := gtk.NotebookNew()
	if err != nil {
		log.Fatal("Failed to create notebook:", err)
	}

	// Create the first page with two stacks
	stack1, err := gtk.StackNew()
	if err != nil {
		log.Fatal("Failed to create stack:", err)
	}

	stack1Page1, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Failed to create box:", err)
	}
	stack1.AddTitled(stack1Page1, "Page 1", "Page 1")

	stack1Page2, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Failed to create box:", err)
	}
	stack1.AddTitled(stack1Page2, "Page 2", "Page 2")

	notebook.AppendPage(stack1, nil)

	// Create the second page with one stack
	stack2, err := gtk.StackNew()
	if err != nil {
		log.Fatal("Failed to create stack:", err)
	}

	stack2Page1, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Failed to create box:", err)
	}
	stack2.AddTitled(stack2Page1, "Page 1", "Page 1")

	notebook.AppendPage(stack2, nil)

	// Create the third page with three stacks
	stack3, err := gtk.StackNew()
	if err != nil {
		log.Fatal("Failed to create stack:", err)
	}

	stack3Page1, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Failed to create box:", err)
	}
	stack3.AddTitled(stack3Page1, "Page 1", "Page 1")

	stack3Page2, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Failed to create box:", err)
	}
	stack3.AddTitled(stack3Page2, "Page 2", "Page 2")

	stack3Page3, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal("Failed to create box:", err)
	}
	stack3.AddTitled(stack3Page3, "Page 3", "Page 3")

	notebook.AppendPage(stack3, nil)

	// Add the notebook to the window and show everything
	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Failed to create window:", err)
	}
	window.Add(notebook)
	window.SetDefaultSize(800, 600)
	window.ShowAll()

	gtk.Main()
}

// package main
//
// import (
// 	"fmt"
// 	"go_graphics/home"
// 	"go_graphics/todo"
// 	"time"
//
// 	"github.com/gotk3/gotk3/gtk"
// )
//
// type Component struct{}
//
// type Page struct {
// 	Title   string
// 	Content *gtk.Box
// }
//
// func main() {
// 	startT := time.Now()
// 	fmt.Println("Running go_graphics")
// 	gtk.Init(nil)
//
// 	// Create a new top-level window
// 	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
// 	if err != nil {
// 		fmt.Println("Unable to create window", err)
// 	}
// 	win.SetTitle("go_graphics")
// 	win.Connect("destroy", func() {
// 		gtk.MainQuit()
// 	})
//
// 	// --------------------------------------------------
//
// 	// vbox
// 	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
// 	if err != nil {
// 		fmt.Println("Unable to create vbox", err)
// 	}
//
// 	// Create a frame/border for navbar
// 	frame, err := gtk.FrameNew("")
// 	if err != nil {
// 		fmt.Println("Unable to create border", err)
// 	}
// 	frame.SetShadowType(gtk.SHADOW_IN)
//
// 	// ---------------------------------------------------
//
// 	// Create a drawing area and set a minimum size
// 	drawArea, err := gtk.DrawingAreaNew()
// 	if err != nil {
// 		fmt.Println("Unable to create drawing area:", err)
// 	}
// 	drawArea.SetSizeRequest(1200, 800)
//
// 	notebook, err := gtk.NotebookNew()
// 	if err != nil {
// 		fmt.Println("Unable to create notebook", err)
// 	}
// 	pages := []Page{
// 		// {
// 		// 	Title:   "DiGraph",
// 		// 	Content: digraph.DiGraphPage(),
// 		// },
// 		{
// 			Title:   "Home",
// 			Content: home.HomePage(),
// 		},
//
// 		// {
// 		// 	Title:   "Page2",
// 		// 	Content: page2.Page2(),
// 		// },
// 		{
// 			Title:   "Todo",
// 			Content: todo.ToDoPage(),
// 		},
// 		// {
// 		// 	Title:   "Finances",
// 		// 	Content: finance.FinancePage(),
// 		// },
// 	}
//
// 	// Add pages to the notebook
// 	for _, page := range pages {
// 		label, _ := gtk.LabelNew(page.Title)
// 		notebook.AppendPage(page.Content, label)
// 	}
//
// 	// Create a horizontal box for the navigation bar
// 	navbar, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 15)
// 	if err != nil {
// 		fmt.Println("Unable to create Navbar", err)
// 	}
//
// 	for i, page := range pages {
// 		button, err := gtk.ButtonNewWithLabel(page.Title)
// 		if err != nil {
// 			fmt.Println("Unable to create button", page, err)
// 		}
// 		currIndex := i
// 		button.Connect("clicked", func() {
// 			fmt.Println(page, "Button was clicked")
// 			notebook.SetCurrentPage(currIndex)
// 		})
// 		navbar.PackStart(button, false, false, 0)
// 	}
//
// 	// Move size_labe stuff here
// 	size_label, err := gtk.LabelNew("")
// 	if err != nil {
// 		fmt.Println("Unable to create label")
// 	}
//
// 	// frame.Add(navbar)
// 	// vbox.PackStart(frame, false, false, 0)
// 	// vbox.PackStart(navbar, false, false, 0)
// 	vbox.PackStart(notebook, true, true, 0)
// 	// vbox.PackStart(size_label, false, false, 0)
//
// 	win.Connect("configure-event", func() {
// 		width, height := win.GetSize()
// 		sizetext := fmt.Sprintf("%d %d", width, height)
// 		size_label.SetText(sizetext)
// 		// fmt.Println("\n Window resized to :", width, height)
// 	})
//
// 	// Add the drawing area to the window
// 	// win.Add(drawArea)
// 	win.Add(vbox)
//
// 	// Show all widgets contained in the window
// 	win.ShowAll()
//
// 	// Begin executing the GTK main loop
// 	stopT := time.Now()
// 	duration := stopT.Sub(startT)
// 	fmt.Println(duration)
// 	win.Maximize()
// 	win.ShowAll()
// 	gtk.Main()
// }
