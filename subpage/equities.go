package subpage

import (
	"github.com/gotk3/gotk3/gtk"
)

func CreateEquities() (*gtk.Stack, error) {
	equitiespage, err := gtk.StackNew()
	if err != nil {
		return nil, err
	}

	return equitiespage, nil
}
