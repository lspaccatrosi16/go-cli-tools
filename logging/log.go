package logging

import (
	"fmt"
	"strings"
	"sync"
)

type Logger struct {
	verbose bool
	test    bool
}

func (l *Logger) SetVerbose(verbose bool) {
	(*l).verbose = verbose
}

func (l *Logger) SetTestMode(test bool) {
	(*l).test = test
}

func (l *Logger) Log(strs ...string) {
	if !l.test {
		for _, str := range strs {
			fmt.Println(str)
		}
	}
}

func (l *Logger) Debug(strs ...string) {
	if l.verbose && !l.test {
		for _, str := range strs {
			fmt.Println(str)
		}
	}
}

func (l *Logger) DebugDivider() {
	if l.verbose && !l.test {
		fmt.Println(strings.Repeat("=", 100))
	}
}

func (l *Logger) LogDivider() {
	if !l.test {
		fmt.Println(strings.Repeat("=", 100))
	}
}

var instance *Logger

var lock = &sync.Mutex{}

func GetLogger() *Logger {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &Logger{}
		}
	}

	return instance
}
