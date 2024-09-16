package routes

import (
	"fmt"
	"net/http"
)

func handleHeartbeatPing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, fmt.Sprintf("%s method not alllowed on this api", r.Method))

			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{"ping": "pong"})
	}
}
