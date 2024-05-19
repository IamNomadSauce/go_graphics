package all

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

var menu = []string{
	"All",
	"Crypto",
	"Equities",
}

func AllTab() *gtk.Box {
	fmt.Println("\n-----------------------\nAllTab\n------------------------\n")

	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		fmt.Println("Error Creating AccountsTab", err)
	}

	return container
}
