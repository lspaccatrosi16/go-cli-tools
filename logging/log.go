package logging

import (
	"fmt"
	"strings"
	"sync"
)

type Logger struct {
	verbose bool
	test    bool
	disable bool
}

func (l *Logger) SetVerbose(verbose bool) {
	(*l).verbose = verbose
}

func (l *Logger) SetTestMode(test bool) {
	(*l).test = test
}

func (l *Logger) SetDisable(disable bool) {
	(*l).disable = disable
}

func (l *Logger) Log(strs ...string) {
	if l.disable {
		return
	}
	if !l.test {
		for _, str := range strs {
			fmt.Println(str)
		}
	}
}

func (l *Logger) Debug(strs ...string) {
	if l.disable {
		return
	}
	if l.verbose && !l.test {
		for _, str := range strs {
			fmt.Println(str)
		}
	}
}

func (l *Logger) DebugDivider() {
	if l.disable {
		return
	}
	if l.verbose && !l.test {
		fmt.Println(strings.Repeat("=", 100))
	}
}

func (l *Logger) LogDivider() {
	if l.disable {
		return
	}
	if !l.test {
		fmt.Println(strings.Repeat("=", 100))
	}
}

var instance *Logger

var lock = &sync.Mutex{}

func GetLogger() *Logger {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		if instance == nil {
			instance = &Logger{}
		}
	}

	return instance
}
