package main

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

/* Log levels.
const (
	CRITICAL Level = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)
*/
///`%{color}[%{module}] %{time:15:04:05.000} %{shortfile:.10s} %{callpath:30s} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
var LOG = logging.MustGetLogger("DFMTP")
var format = logging.MustStringFormatter(
	`%{color}[%{module}:%{pid:03x}] %{time:15:04:05.000} %{shortfile:10s}  %{level:.5s} %{shortfunc:20s}: %{color:reset} %{message}`,
)

type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

/**
* log初始化程序，包括输出log的格式及打印级别
*@logLevel     打印级别
 */
func logInit(logLevel uint8) {
	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Println(err)
	}

	var level logging.Level = logging.Level(logLevel)

	backend1 := logging.NewLogBackend(logFile, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(level, "")
	logging.SetBackend(backend1Leveled)
}

func DEBUG(format string, args ...interface{}) {
	LOG.Debugf(format, args...)
}

func ERROR(format string, args ...interface{}) {
	LOG.Errorf(format, args...)
}

func NOTICE(format string, args ...interface{}) {
	LOG.Noticef(format, args...)
}

func WARNING(format string, args ...interface{}) {
	LOG.Warningf(format, args...)
}

func INFO(format string, args ...interface{}) {
	LOG.Infof(format, args...)
}

func FATAL(format string, args ...interface{}) {
	LOG.Criticalf(format, args...)
}

/*func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, &format, args...)
}*/

/**
* log初始化程序，包括输出log的格式及打印级别
 */
func logTest1() {
	LOG.Debugf("debug %s", Password("secret"))
	LOG.Debugf("debug %s", Password("secret"))
	LOG.Infof("info")
	LOG.Noticef("notice")
	LOG.Warningf("warning")
	LOG.Errorf("xiaorui.cc")
	LOG.Criticalf("太严重了%d", 2)
	LOG.Debugf("debug %d", 3)
}
