package subpage

import (
	"github.com/gotk3/gotk3/gtk"
)

func CreateSavings() (*gtk.Stack, error) {
	savingspage, err := gtk.StackNew()
	if err != nil {
		return nil, err
	}

	return savingspage, nil
}
