package reader

import (
	"os"
)

type TxtReader struct {
	contentReader
}

func NewTxtReader() *TxtReader {
	return &TxtReader{}
}

func (txt *TxtReader) Load(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	txt.setContent(string(data))
	return nil
}
