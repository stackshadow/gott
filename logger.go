package gott

import (
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger initializes a zerolog-logger
//
// if an filename is configured, we use the lumberjack-logger to write to it
//
// if the filename contains a folder, it will be created
func NewLogger(cnf loggingConfig) *zerolog.Logger {

	var writers []io.Writer
	var err error

	// per default we use consolewriter
	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})

	// if an filename provided, we use the lumberjack-logger
	if cnf.Filename != "" {

		// created needed folder
		configFolder := path.Dir(cnf.Filename)
		if err = os.MkdirAll(configFolder, 0744); err != nil {
			log.Error().Err(err).Str("path", configFolder).Msg("can't create log directory")
		}

		// add the writer
		if err == nil {
			var fileWriter io.Writer = &lumberjack.Logger{
				Filename:   cnf.Filename,
				MaxBackups: cnf.MaxBackups, // files
				MaxSize:    cnf.MaxSize,    // megabytes
				MaxAge:     cnf.MaxAge,     // days
			}

			writers = append(writers, zerolog.ConsoleWriter{Out: fileWriter, NoColor: true})
		}
	}
	mw := io.MultiWriter(writers...)

	logger := zerolog.New(mw).With().Timestamp().Logger()

	logger.Info().
		Str("fileName", cnf.Filename).
		Int("maxSizeMB", cnf.MaxSize).
		Int("maxBackups", cnf.MaxBackups).
		Int("maxAgeInDays", cnf.MaxAge).
		Msg("logging configured")

	return &logger
}
