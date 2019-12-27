package log

import (
	"fmt"
	"os"
)

type consoleWriter struct {
}

func (c *consoleWriter) Write(data interface{}) error {
	str := fmt.Sprintf("%v\n", data)
	_, err := os.Stdout.Write([]byte(str))
	return err
}

func NewConsole() *consoleWriter {
	return &consoleWriter{}
}
