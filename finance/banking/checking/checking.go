package checking

import (
  "github.com/gotk3/gotk3/gtk"
)


func CheckingPage() (*gtk.Box, error) {
  box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  if err != nil {
    return nil, err
  }

  label, err := gtk.LabelNew("CheckingPage")
  if err != nil {
    return nil, err
  }
  box.PackStart(label, false, false, 0)
  return box, nil
}


