package util

import (
	"golang.org/x/exp/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}
