package util

import (
	"fmt"
	"strings"
)

//Logger For Systemd logs
var prefixes = []string{"<0>", "<1>", "<2>", "<3>", "<4>", "<5>", "<6>", "<7>"}

// Debug levels corresponding to the syslog levels
const (
	EMERG int = iota
	ALERT
	CRIT
	ERR
	WARNING
	NOTICE
	INFO
	DEBUG
)

// func log(level DebugLevel, args ...string) {
// 	fmt.Println(prefixes[level], args)
// }

type logger struct{}

// Log is the namespace for the logger functions
var Log = logger{}

func (l logger) Log(level int, str string, args ...interface{}) {
	fmt.Println(prefixes[level] +
		strings.Replace(fmt.Sprintf(str, args...),
			"\n", "\n"+prefixes[level], -1))
}

func (l logger) Emerg(str string, args ...interface{}) {
	l.Log(EMERG, str, args...)
}

func (l logger) Alert(str string, args ...interface{}) {
	l.Log(ALERT, str, args...)
}

func (l logger) Crit(str string, args ...interface{}) {
	l.Log(CRIT, str, args...)
}

func (l logger) Err(str string, args ...interface{}) {
	l.Log(ERR, str, args...)
}

func (l logger) Warning(str string, args ...interface{}) {
	l.Log(WARNING, str, args...)
}

func (l logger) Notice(str string, args ...interface{}) {
	l.Log(NOTICE, str, args...)
}

func (l logger) Info(str string, args ...interface{}) {
	l.Log(INFO, str, args...)
}

func (l logger) Debug(str string, args ...interface{}) {
	l.Log(DEBUG, str, args...)
}
