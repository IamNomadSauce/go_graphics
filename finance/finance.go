package finance

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

// Assuming tabLabel is the label you want to dynamically update.
var tabLabel *gtk.Label

type Tab struct {
	Label   string
	Content *gtk.Box
}

var menu = []Tab{
	{
		Label: "Accounts",
		Content: func() *gtk.Box {
			box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
			if err != nil {
				fmt.Println("Error creating box:", err)
				return nil
			}
			return box
		}(),
	},
	{
		Label: "All",
		Content: func() *gtk.Box {
			box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
			if err != nil {
				fmt.Println("Error creating box:", err)
				return nil
			}
			return box
		}(),
	},
	// Add more main tabs as needed
}

func FinancePage() *gtk.Box {
	fmt.Println("\n-------------------------\nFinances\n-------------------------\n")
	notebook, err := gtk.NotebookNew()
	if err != nil {
		fmt.Println("Error creating notebook", err)
	}

	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		fmt.Println("Error creating Finances container", err)
	}

	navbar, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		fmt.Println("Error creating FinancePage navbar", err)
	}

	// Initialize the label with default text.
	tabLabel, err = gtk.LabelNew("Account") // Default tab name
	if err != nil {
		fmt.Println("Error creating label", err)
	}

	// setupNavbar(navbar, notebook)
	setupTabs(notebook)
	container.PackStart(navbar, false, false, 10)
	container.PackStart(tabLabel, false, false, 10)
	container.PackStart(notebook, false, false, 0)

	return container
}

func setupSubTabs(parentNotebook *gtk.Notebook, subTabs []Tab) {
	subNotebook, err := gtk.NotebookNew()
	if err != nil {
		fmt.Println("Error creating sub-notebook", err)
		return
	}

	for _, tab := range subTabs {
		tabLabel, _ := gtk.LabelNew(tab.Label)
		subNotebook.AppendPage(tab.Content, tabLabel)
	}

	parentNotebook.AppendPage(subNotebook, nil)
}

func setupTabs(notebook *gtk.Notebook) {
	// Define sub-tabs for each main setupTabs
	accountSubTabs := []Tab{
		{
			Label: "Checking",
			Content: func() *gtk.Box {
				box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
				if err != nil {
					fmt.Println("Error creating box:", err)
					return nil
				}
				return box
			}(),
		},
		{
			Label: "Savings",
			Content: func() *gtk.Box {
				box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
				if err != nil {
					fmt.Println("Error creating box:", err)
					return nil
				}
				return box
			}(),
		},
		// Add more sub-tabs as needed
	}

	allSubTabs := []Tab{
		{
			Label: "All Accounts",
			Content: func() *gtk.Box {
				box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
				if err != nil {
					fmt.Println("Error creating box:", err)
					return nil
				}
				return box
			}(),
		},
		// Add more sub-tabs as needed
	}

	// Add main tabs and their sub-tabs
	for _, tab := range menu {
		tabLabel, _ := gtk.LabelNew(tab.Label)
		tabContent, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
		notebook.AppendPage(tabContent, tabLabel)

		switch tab.Label {
		case "Accounts":
			setupSubTabs(notebook, accountSubTabs)
		case "All":
			setupSubTabs(notebook, allSubTabs)
			// Add cases for other main tabs and their sub-tabs
		}
	}
}

// var menu = []Tab{
// {
// 	Label:   "Accounts",
// 	Content: accounts.AccountsTab(),
// },
// {
// 	Label:   "All",
// 	Content: all.AllTab(),
// },
// {
// 	Label:   "Equities",
// 	Content: equities.EquitiesTab(),
// },
// {
// 	Label:   "Crypto",
// 	Content: crypto.CryptoTab(),
// },
// }

// func setupTabs(notebook *gtk.Notebook) {
// 	for _, tab := range menu {
// 		tabLabel, _ := gtk.LabelNew(tab.Label)
// 		// if err != nil {
// 		// 	fmt.Println("Error creating label:", err)
// 		// 	continue
// 		// }
// 		notebook.AppendPage(tab.Content, tabLabel)
// 	}
// }

// func setupNavbar(sb *gtk.Box, notebook *gtk.Notebook) {
// 	fmt.Println("Setup Finance navbar")
// 	menuItems := []string{"API", "Account", "Trade"}
//
// 	for index, item := range menuItems {
// 		btn, err := gtk.ButtonNewWithLabel(item)
// 		if err != nil {
// 			fmt.Println("Error creating Finance navbar Button", err)
// 			continue
// 		}
// 		btn.Connect("clicked", func() {
// 			// Update the label text based on the button clicked.
// 			tabLabel.SetText(item)
// 			fmt.Println("Menu Item Selected", item)
// 			notebook.SetCurrentPage(index)
// 		})
// 		sb.PackStart(btn, false, false, 5)
// 	}
// 	sb.ShowAll()
// }
