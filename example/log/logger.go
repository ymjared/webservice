package log

import "fmt"

type LogWriter interface {
	Write(data interface{}) error
}

type Logger struct {
	writes []LogWriter
}

func (l *Logger) RegistryWriter(Writer LogWriter) {
	l.writes = append(l.writes, Writer)
}

func (l *Logger) Log(data interface{}) {
	for _, write := range l.writes {
		write.Write(data)
	}
}

func NewLog() *Logger {
	return &Logger{}
}

func CreateLogger() *Logger {
	log := NewLog()

	cw := NewConsole()
	log.RegistryWriter(cw)

	fw := NewFile()
	if err := fw.SetFile("./log.txt"); err != nil {
		fmt.Println(err)
	} else {
		log.RegistryWriter(fw)
	}
	return log
}
