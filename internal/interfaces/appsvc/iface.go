package appsvc

import "context"

type AppServicer interface {
	AppGetVersion() string
	AppHandleRequest(requestCtx context.Context) error
}
