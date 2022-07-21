package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	stdlog "log"
)

var (
	appname = "kubebe"
)

func main() {
	// setup structured log
	logsink1 := zerolog.New(os.Stdout).With().Str("app", appname).Timestamp().Logger()
	logsink2 := zerolog.New(os.Stdout).With().Str("level", "debug").Str("app", appname).Timestamp().Logger()
	stdlog.SetFlags(0)
	stdlog.SetOutput(logsink2)
	log.Logger = logsink1

	// log something
	stdlog.Println("something happened!")
	log.Error().Caller().Msg("friggin error here man!")
}
