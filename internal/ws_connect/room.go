package ws_connect

type Room struct {
	RoomManager *RoomManager

	// Room Id
	RoomId string

	// Registered Clients.
	Clients map[*Client]bool

	// Inbound messages from the Clients.
	broadcast chan []byte

	// Register requests from the Clients.
	Register chan *Client

	// Unregister requests from Clients.
	Unregister chan *Client

	online int64
}

func NewRoom(roomManager *RoomManager, roomId string) *Room {
	return &Room{
		RoomManager: roomManager,
		RoomId:      roomId,
		broadcast:   make(chan []byte),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Clients:     make(map[*Client]bool),
	}
}

func (h *Room) Broadcast(buf []byte) {
	h.broadcast <- buf
}

func (h *Room) Run() {

BREAK_LOOP:
	for {
		select {

		case client := <-h.Register:
			h.Clients[client] = true
			h.online++

		case client := <-h.Unregister:

			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
				h.online--
			}

			if h.online == 0 {
				h.RoomManager.RemoveRoom(h.RoomId)
				break BREAK_LOOP
			}

		case message := <-h.broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}

	}

	return
}
