package accounts

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

var menu = []string{
	"All",
	"Crypto",
	"Equities",
}

func AccountsTab() *gtk.Box {
	fmt.Println("\n-----------------------\nAccounts\n------------------------\n")

	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		fmt.Println("Error Creating AccountsTab", err)
	}

	return container
}
