package savings

import (
  "github.com/gotk3/gotk3/gtk"
)


func SavingsPage() (*gtk.Box, error) {
  box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  if err != nil {
    return nil, err
  }

  label, err := gtk.LabelNew("SavingsTab")
  if err != nil {
    return nil, err
  }
  box.PackStart(label, false, false, 0)
  return box, nil
}


