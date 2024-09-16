package routes

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	broadcast      = make(chan []byte)
	clientMap      = make(map[*websocket.Conn]bool)
	clientMapMutex sync.Mutex
)

func handleWebSocket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// CheckOrigin: func(r *http.Request) bool {
			// 	origin := r.Header.Get("Origin")

			// },
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error(
				"While upgrading connection to websocket",
				slog.Any("error", err),
			)

			respondWithError(w, http.StatusBadRequest, "unable to upgrade connection to websocket")
		}

		defer func() {
			err := conn.Close()
			if err != nil {
				slog.Error(
					"While closing websocket connection",
					slog.Any("error", err),
				)
			}
		}()

		// Register the client
		clientMapMutex.Lock()
		clientMap[conn] = true
		clientMapMutex.Unlock()

		slog.Debug(
			"Current client map",
			slog.Int("numberOfClients", len(clientMap)),
		)

		//Broadcast to all clients, since muliple computers can open page in multiple tabs at same time
		broadcast <- []byte("Hi socket")
	}
}

// Sending data to all clients
func HandleBroadcast() {
	for {
		for client := range clientMap {
			err := client.WriteMessage(websocket.TextMessage, []byte("hi"))
			if err != nil {
				slog.Error(
					"While sending message to client, removing this client",
					slog.Any("error", err),
				)

				clientMapMutex.Lock()
				err := client.Close()
				if err != nil {
					slog.Error(
						"While closing connection during broadcast",
						slog.Any("error", err),
					)
				}

				delete(clientMap, client)
				clientMapMutex.Unlock()
			}

			slog.Debug(
				"Current clients map",
				slog.Int("numberOfClients", len(clientMap)),
			)
		}
	}
}
