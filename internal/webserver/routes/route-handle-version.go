package routes

import (
	"fmt"
	"net/http"

	"github.com/GaikwadPratik/signoztest/internal/interfaces/appsvc"
)

func handleVersion(apiHandlerService appsvc.AppServicer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, fmt.Sprintf("%s method not allowed on this api", r.Method))

			return
		}
		respondWithJSON(w, http.StatusOK, apiHandlerService.AppGetVersion())
	}
}
