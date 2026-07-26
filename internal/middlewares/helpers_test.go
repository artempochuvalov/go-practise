package middlewares

import (
	"bytes"
	"log/slog"
)

func newTestLogger() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer

	log := slog.New(
		slog.NewTextHandler(&buf, nil),
	)

	return log, &buf
}
