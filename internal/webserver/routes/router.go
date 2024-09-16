package routes

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/GaikwadPratik/signoztest/internal/interfaces/appsvc"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
	middleware "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

var serviceName = os.Getenv("SERVICE_NAME")

type WebServerRoutesInput struct {
	AppCtx         context.Context
	AppService     appsvc.AppServicer
	LogLevel       *slog.LevelVar
	AllowedOrigins []string
}

func NewWebserverRoutes(input WebServerRoutesInput) *mux.Router {
	router := mux.NewRouter()
	//router.SkipClean(true)

	router.Use(middleware.Middleware(serviceName))
	router.HandleFunc("/api/ping", handleHeartbeatPing()).Methods(http.MethodGet)
	router.HandleFunc("/api/version", handleVersion(input.AppService)).Methods(http.MethodGet)

	router.HandleFunc("/api/loglevel", handleLogLevelChange(input.LogLevel)).Methods(http.MethodPut)
	router.Handle("/api/metrics", promhttp.Handler()).Methods(http.MethodGet)

	router.HandleFunc("/ws", handleWebSocket()).Methods(http.MethodGet)

	return router
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	var response []byte
	var err error

	switch x := payload.(type) {
	case []byte:
		response = x
	default:
		response, err = json.Marshal(payload)
		if err != nil {
			slog.Error(
				"error",
				slog.Any("error", err),
				slog.Any("payload", payload),
			)
			respondWithError(w, http.StatusInternalServerError, err.Error())

			return
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		slog.Error(
			"While writing the data into http writer",
			slog.Any("error", err),
			slog.Any("payload", payload),
		)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
