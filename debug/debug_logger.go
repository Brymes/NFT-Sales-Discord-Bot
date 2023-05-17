package debug

import (
	"io"
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

func DbgInit() {
	f, err := os.OpenFile("application.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("error opening file: %v", err)
	}

	Log = logrus.New()

	//log.Formatter = &logrus.JSONFormatter{}

	Log.SetReportCaller(true)

	mw := io.MultiWriter(os.Stdout, f)

	Log.SetOutput(mw)
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
