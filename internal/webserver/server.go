package webserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/GaikwadPratik/signoztest/internal/interfaces/appsvc"
	"github.com/GaikwadPratik/signoztest/internal/webserver/routes"
)

type WebserverConf struct {
	Port         int
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
	Graceful     time.Duration
	StartWait    time.Duration
}

type WebserverDependencies struct {
	Conf       *WebserverConf
	AppService appsvc.AppServicer
	LogLevel   *slog.LevelVar
}

func Initiate(appCtx context.Context, input WebserverDependencies, errChan chan<- error) {
	srvRouteInput := routes.WebServerRoutesInput{
		AppCtx:     appCtx,
		AppService: input.AppService,
		LogLevel:   input.LogLevel,
	}

	webserverRouter := routes.NewWebserverRoutes(srvRouteInput)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", input.Conf.Port),
		WriteTimeout: input.Conf.WriteTimeout,
		ReadTimeout:  input.Conf.ReadTimeout,
		IdleTimeout:  input.Conf.IdleTimeout,
		Handler:      webserverRouter,
	}

	go func() {
		shutdownCtx, cancelShutdownCtxFn := context.WithTimeout(appCtx, input.Conf.Graceful)
		defer cancelShutdownCtxFn()
		<-appCtx.Done()
		slog.Debug("Received done on context, shutting down http server")
		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error(
				"While shutting down server",
				slog.Any("error", err),
			)

			return
		}

		slog.Info("http server shutdown correctly")
	}()

	var srvStartError error
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(
				"While initiating the web server",
				slog.Any("error", err),
				slog.Any("config", input.Conf),
			)

			srvStartError = fmt.Errorf("unable to initialize web server for app")
			errChan <- srvStartError
		}
	}()

	// For initiating broadcast to all channels
	go routes.HandleBroadcast()

	<-time.After(input.Conf.StartWait)
	if srvStartError == nil {
		//server is running, so everything went smmoth,
		// this might cause panic if an error is being sent later stage but there is no alternative for now
		// if srvStartError is not nil, then the error has already been sent by another go routine, just close the channel
		slog.Info(
			"Started web server",
			slog.Any("webserverConfig", input.Conf),
		)

		errChan <- nil
	}

	close(errChan)
}
