package subpage

import (
	"github.com/gotk3/gotk3/gtk"
)

type CryptoPage struct{}

func CreateCrypto() (*gtk.Stack, error) {
	cryptopage, err := gtk.StackNew()
	if err != nil {
		return nil, err
	}

	return cryptopage, nil
}
