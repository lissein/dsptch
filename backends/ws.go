package backends

import (
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type WebSocketBackend struct {
	config  *Config
	clients map[int]*net.Conn

	// Store the latest client ID
	clientId int

	mutex sync.Mutex
}

type WebSocketPayload struct {
	Targets []int
	Content string
}

func NewWebSocketBackend(config *Config) (Backend, error) {
	backend := &WebSocketBackend{
		config:   config,
		clients:  make(map[int]*net.Conn),
		clientId: 0,
	}

	return backend, nil
}

func (backend *WebSocketBackend) Listen(messages chan Message) {
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
				messages <- Message{
					Source:  "websocket",
					Payload: string(msg),
				}
			}
		}()
	}))
}

func (backend *WebSocketBackend) Handle(message Message) error {
	payload := message.Payload.(WebSocketPayload)

	conns := make([]*net.Conn, 0)
	for _, target := range payload.Targets {
		client, found := backend.clients[target]
		if !found {
			continue
		}
		conns = append(conns, client)
	}

	for _, client := range conns {
		err := wsutil.WriteServerText(*client, []byte(message.Payload.(string)))
		if err != nil {
			backend.config.Logger.Error("Failed to send websocket message")
		}
	}

	return nil
}
