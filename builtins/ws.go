package builtins

import (
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/lissein/dsptch/backends"
)

type WebSocketBackend struct {
	config  *backends.Config
	clients map[int]*net.Conn

	// Store the latest client ID
	clientID int

	mutex sync.Mutex
}

type WebSocketListenPayload struct {
	FromClient int
	Content    string
}

type WebSocketHandlePayload struct {
	Content string
	Targets []int
}

func NewWebSocketBackend(config *backends.Config) (backends.Backend, error) {
	backend := &WebSocketBackend{
		config:   config,
		clients:  make(map[int]*net.Conn),
		clientID: 0,
	}

	return backend, nil
}

func (backend *WebSocketBackend) Listen(messages chan backends.Message) {
	http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			backend.config.Logger.Error(err)
		}

		var clientID int

		backend.mutex.Lock()
		backend.clientID++
		backend.clients[backend.clientID] = &conn
		clientID = backend.clientID
		backend.mutex.Unlock()

		go func() {
			defer conn.Close()

			for {
				msg, err := wsutil.ReadClientText(conn)
				if err != nil {
					backend.config.Logger.Error(err)
					break
				}

				// TODO Add clientID to the message
				messages <- backends.Message{
					Source: "websocket",
					Payload: &WebSocketListenPayload{
						Content:    string(msg),
						FromClient: clientID,
					},
				}
			}
		}()
	}))
}

func (backend *WebSocketBackend) Handle(message backends.Message) error {
	payload := message.Payload.(WebSocketHandlePayload)

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
