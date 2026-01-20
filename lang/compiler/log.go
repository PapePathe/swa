package compiler

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	dbg *log.Logger
}

func NewLogger(prefix string) *Logger {
	l := Logger{
		dbg: log.New(
			os.Stdout,
			fmt.Sprintf("%s %s ", "DBG", prefix),
			log.Lmsgprefix,
		),
	}

	return &l
}

func (l Logger) Debug(msg string) {
	l.dbg.Println(msg)
}

func (l Logger) Restore(prefix string) {
	l.dbg.SetPrefix(prefix)
}

func (l Logger) Step(prefix string) string {
	oldprefix := l.dbg.Prefix()
	newprefix := fmt.Sprintf("%s %s ", oldprefix, prefix)

	l.dbg.SetPrefix(newprefix)

	return oldprefix
}
