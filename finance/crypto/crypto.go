package crypto

import (
  "github.com/gotk3/gotk3/gtk"
)


func CryptoPage() (*gtk.Box, error) {
  box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
  if err != nil {
    return nil, err
  }

  label, err := gtk.LabelNew("CryptoTab")
  if err != nil {
    return nil, err
  }
  box.PackStart(label, false, false, 0)
  return box, nil
}


