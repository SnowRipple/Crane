package logger

import (
	log "github.com/op/go-logging"
	stdlog "log"
	"os"
)

const (
	LOGGER_NAME = "crane"
)

var logger = log.MustGetLogger(LOGGER_NAME)

//Customize logger
func init() {

	//Customize the output format
	log.SetFormatter(log.MustStringFormatter("▶ %{level:.1s} 0x%{id:x}  %{message}"))

	// Setup one stdout and one syslog backend.
	logBackend := log.NewLogBackend(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
	logBackend.Color = true

	syslogBackend, err := log.NewSyslogBackend("")
	if err != nil {
		logger.Fatal("Failed to set up syslog backend:", err)
	}

	// Combine them both into one log backend.
	log.SetBackend(logBackend, syslogBackend)

	//Default log level
	log.SetLevel(log.NOTICE, LOGGER_NAME)
}

//Get the instance...sorry, forgot it's go...value of type log.Logger.
func GetLogger() *log.Logger {
	return logger
}
