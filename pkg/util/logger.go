package util

import (
	"golang.org/x/exp/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	// Init the new structured logger with the default text handler
	Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}
