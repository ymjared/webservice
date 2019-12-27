package log

import (
	"errors"
	"fmt"
	"os"
)

var errorFileNotFound = errors.New("file nor created")

type fileWriter struct {
	file *os.File
}

func (f *fileWriter) SetFile(filename string) (err error) {
	if f.file != nil {
		f.file.Close()
	}

	f.file, err = os.Create(filename)
	return err
}

func (f *fileWriter) Write(data interface{}) error {
	if f.file == nil {
		return errorFileNotFound
	}

	str := fmt.Sprintf("%v\n", data)
	_, err := f.file.Write([]byte(str))
	return err
}

func NewFile() *fileWriter {
	return &fileWriter{}
}
