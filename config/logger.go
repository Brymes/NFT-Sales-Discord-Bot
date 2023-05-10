package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func CreateServiceLogger(fileName string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		log.Println("Failed to log to file, using default stderr")
	}
	return logger
}

func InitRequestLogger(service string) (*bytes.Buffer, *log.Logger) {
	b := &bytes.Buffer{}
	prefix := fmt.Sprintf("%v: ", strings.ToUpper(service))
	reqLogger := log.New(b, prefix, log.Ltime|log.Lshortfile)

	return b, reqLogger
}
