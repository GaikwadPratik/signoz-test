package routes

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/GaikwadPratik/signoztest/internal/entity"
)

func handleLogLevelChange(currentLogLevel *slog.LevelVar) http.HandlerFunc {
	//Write a response when we are done
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			respondWithError(w, http.StatusMethodNotAllowed, fmt.Sprintf("%s method not alllowed on this api", r.Method))

			return
		}

		// Parse request body
		requestBody := entity.RouteHandleLogLevel{}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		if err := decoder.Decode(&requestBody); err != nil {
			slog.Error(
				"While unmarshalling request body for changing log level",
				slog.Any("error", err),
				slog.Any("currentLogLevel", currentLogLevel),
			)

			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Invalid json received. Error: %v", err))

			return
		}

		// Change log level
		if strings.EqualFold(requestBody.LogLevel, "debug") {
			currentLogLevel.Set(slog.LevelDebug)
			slog.Debug("Log level is set to Debug")
		} else if strings.EqualFold(requestBody.LogLevel, "info") {
			currentLogLevel.Set(slog.LevelInfo)
			slog.Debug("Log level is set to Info")
		} else if strings.EqualFold(requestBody.LogLevel, "warn") {
			currentLogLevel.Set(slog.LevelWarn)
			slog.Debug("Log level is set to Warn")
		} else if strings.EqualFold(requestBody.LogLevel, "error") {
			currentLogLevel.Set(slog.LevelError)
			slog.Debug("Log level is set to Error")
		} else {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Unknown log level received: %v", requestBody.LogLevel))

			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"result": fmt.Sprintf("Updated log level successfully to: %v", requestBody.LogLevel)})
	}
}
