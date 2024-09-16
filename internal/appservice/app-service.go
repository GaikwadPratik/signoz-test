package appservice

import (
	"context"
	"log/slog"
)

type app struct {
}

func New() app {
	return app{}
}

func (a app) AppHandleRequest(ctx context.Context) error {
	return nil
}

func (a app) AppGetVersion() string {
	slog.Info("Received get version")
	return "v0.0.1"
}
