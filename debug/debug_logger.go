package debug

import (
	"io"
	logging "log"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func DbgInit() {
	f, err := os.OpenFile("application.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		logging.Fatalf("error opening file: %v", err)
	}

	log = logrus.New()

	//log.Formatter = &logrus.JSONFormatter{}

	log.SetReportCaller(true)

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
}

// Fatalf ...
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

// Fatalln ...
func Fatalln(v ...interface{}) {
	log.Fatalln(v...)
}

// Panicf ...
func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}

// Fatal ...
func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

// Println ...
func Println(v ...interface{}) {
	log.Println(v...)
}

var (

	// ConfigError ...
	ConfigError = "%v type=config.error"

	// HTTPError ...
	HTTPError = "%v type=http.error"

	// HTTPWarn ...
	HTTPWarn = "%v type=http.warn"

	// HTTPInfo ...
	HTTPInfo = "%v type=http.info"
)
