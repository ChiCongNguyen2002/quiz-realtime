package websocket

type Hub struct {
	Clients          map[*Client]bool
	Register         chan *Client
	Unregister       chan *Client
	Sessions         map[string]map[*Client]bool
	SessionBroadcast chan *SessionMessage
}

type SessionMessage struct {
	SessionID string
	Data      []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:          make(map[*Client]bool),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		Sessions:         make(map[string]map[*Client]bool),
		SessionBroadcast: make(chan *SessionMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			if client.SessionID != "" {
				if h.Sessions[client.SessionID] == nil {
					h.Sessions[client.SessionID] = make(map[*Client]bool)
				}
				h.Sessions[client.SessionID][client] = true
			}

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				if client.SessionID != "" {
					if _, ok := h.Sessions[client.SessionID]; ok {
						delete(h.Sessions[client.SessionID], client)
						if len(h.Sessions[client.SessionID]) == 0 {
							delete(h.Sessions, client.SessionID)
						}
					}
				}
				close(client.Send)
			}

		case sm := <-h.SessionBroadcast:
			if clients, ok := h.Sessions[sm.SessionID]; ok {
				for client := range clients {
					select {
					case client.Send <- sm.Data:
					default:
						close(client.Send)
						delete(h.Clients, client)
						delete(clients, client)
					}
				}
			}
		}
	}
}

func (h *Hub) BroadcastToSession(sessionID string, data []byte) {
	h.SessionBroadcast <- &SessionMessage{
		SessionID: sessionID,
		Data:      data,
	}
}
