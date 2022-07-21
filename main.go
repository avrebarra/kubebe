package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	stdlog "log"
)

var (
	appname       = "kubebe"
	prettylogging = true
)

func main() {
	// setup structured log
	logsink := log.Logger
	if prettylogging {
		logsink = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	logsink = logsink.With().Str("app", appname).Timestamp().Logger()

	stdlog.SetFlags(0)
	stdlog.SetOutput(logsink.With().Str("level", "debug").Logger())
	log.Logger = logsink

	// log something
	stdlog.Println("initializing server...")
}
