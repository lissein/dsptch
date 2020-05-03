package backends

import (
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	utils "github.com/lissein/dsptch/utils"
)

type WebSocketBackend struct {
	config  *BackendConfig
	clients map[int]*net.Conn

	// Store the latest client ID
	clientId int

	mutex sync.Mutex
}

func NewWebSocketBackend(config *BackendConfig) *WebSocketBackend {
	backend := &WebSocketBackend{
		config:   config,
		clients:  make(map[int]*net.Conn),
		clientId: 0,
	}

	return backend
}

func (backend *WebSocketBackend) Listen(messages chan BackendInputMessage) {
	http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			backend.config.Logger.Error(err)
		}

		backend.mutex.Lock()
		backend.clientId++
		backend.clients[backend.clientId] = &conn
		backend.mutex.Unlock()

		go func() {
			defer conn.Close()

			for {
				msg, err := wsutil.ReadClientText(conn)
				if err != nil {
					backend.config.Logger.Error(err)
					break
				}

				// TODO Add clientId to the message
				messages <- BackendInputMessage{
					Source:  "websocket",
					Content: string(msg),
				}
			}
		}()
	}))
}

func (backend *WebSocketBackend) HandleMessage(message BackendOutputMessage) error {
	targets := utils.ToIntSlice(message.Targets)

	conns := make([]*net.Conn, len(targets))
	for _, target := range targets {
		client, found := backend.clients[target]
		if !found {
			continue
		}
		conns = append(conns, client)
	}

	for _, client := range backend.clients {
		err := wsutil.WriteServerText(*client, []byte(message.Content.(string)))
		if err != nil {
			backend.config.Logger.Error("Failed to send websocket message")
		}
	}

	return nil
}
