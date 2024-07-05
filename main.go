package main

import (
	"fmt"
	"gogtk/digraph"
	"gogtk/todo"
  "gogtk/finance"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

type Component struct{}

type Page struct {
	Title   string
	Content *gtk.Box
}

func main() {
	startT := time.Now()
	fmt.Println("Running go_graphics")
	gtk.Init(nil)

	// Create a new top-level window
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		fmt.Println("Unable to create window", err)
	}
	win.SetTitle("go_graphics")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// --------------------------------------------------

	// vbox
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		fmt.Println("Unable to create vbox", err)
	}

	// Create a frame/border for navbar
	frame, err := gtk.FrameNew("")
	if err != nil {
		fmt.Println("Unable to create border", err)
	}
	frame.SetShadowType(gtk.SHADOW_IN)

	// ---------------------------------------------------

	// Create a drawing area and set a minimum size
	drawArea, err := gtk.DrawingAreaNew()
	if err != nil {
		fmt.Println("Unable to create drawing area:", err)
	}
	drawArea.SetSizeRequest(1200, 800)

  // Comment for magic temp


	notebook, err := gtk.NotebookNew()
	if err != nil {
		fmt.Println("Unable to create notebook", err)
	}
	pages := []Page{
		{
		  Title:   "Finances",
		  Content: finance.FinancePage(),
		},
		{
			Title:   "DiGraph",
			Content: digraph.DiGraphPage(),
		},
		// {
		// 	Title:   "Home",
		// 	Content: home.HomePage(),
		// },

		// {
		// 	Title:   "Page2",
		// 	Content: page2.Page2(),
		// },
		{
			Title:   "Todo",
			Content: todo.ToDoPage(),
		},
	}

	// Add pages to the notebook
	for _, page := range pages {
		label, _ := gtk.LabelNew(page.Title)
		notebook.AppendPage(page.Content, label)
	}

	// Create a horizontal box for the navigation bar
	navbar, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 15)
	if err != nil {
		fmt.Println("Unable to create Navbar", err)
	}

	for i, page := range pages {
		button, err := gtk.ButtonNewWithLabel(page.Title)
		if err != nil {
			fmt.Println("Unable to create button", page, err)
		}
		currIndex := i
		button.Connect("clicked", func() {
			fmt.Println(page, "Button was clicked")
			notebook.SetCurrentPage(currIndex)
		})
		navbar.PackStart(button, false, false, 0)
	}

	// Move size_labe stuff here
	size_label, err := gtk.LabelNew("")
	if err != nil {
		fmt.Println("Unable to create label")
	}

	// frame.Add(navbar)
	// vbox.PackStart(frame, false, false, 0)
	// vbox.PackStart(navbar, false, false, 0)
	vbox.PackStart(notebook, true, true, 0)
	// vbox.PackStart(size_label, false, false, 0)

	win.Connect("configure-event", func() {
		width, height := win.GetSize()
		sizetext := fmt.Sprintf("%d %d", width, height)
		size_label.SetText(sizetext)
		// fmt.Println("\n Window resized to :", width, height)
	})

	// Add the drawing area to the window
	// win.Add(drawArea)
	win.Add(vbox)

	// Show all widgets contained in the window
	win.ShowAll()

	// Begin executing the GTK main loop
	stopT := time.Now()
	duration := stopT.Sub(startT)
	fmt.Println(duration)
	win.Maximize()
	win.ShowAll()
	gtk.Main()
}
