package logging

import (
	"io"
	"os"

	"github.com/charmbracelet/log"
)

var Log *log.Logger = log.New(os.Stderr)

func SetupLogger() (io.Closer, error) {
	f, err := os.OpenFile("quick.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	Log = log.NewWithOptions(io.MultiWriter(os.Stderr, f), log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Prefix:          "📁",
		Level:           log.DebugLevel,
	})
	return f, nil
}
