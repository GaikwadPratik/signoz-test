package appservice

import "context"

type app struct {
}

func New() app {
	return app{}
}

func (a app) AppHandleRequest(ctx context.Context) error {
	return nil
}

func (a app) AppGetVersion() string {
	return "v0.0.1"
}
