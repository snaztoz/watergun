package socket

func NewHub() *hub {
	return &hub{
		messageProcessor: newMessageProcessor(),
		clients:          make(map[string]*client),
		register:         make(chan *client),
		unregister:       make(chan *client),
	}
}

type hub struct {
	messageProcessor *messageProcessor
	clients          map[string]*client
	register         chan *client
	unregister       chan *client
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

			// case message := <-h.broadcast:
			// 	for client := range h.clients {
			// 		select {
			// 		case client.send <- message:
			// 		default:
			// 			close(client.send)
			// 			delete(h.clients, client)
			// 		}
			// 	}
		}
	}
}

func (h *hub) processMessage(client *client, message *ReadMessage) error {
	return h.messageProcessor.process(client, message)
}

func (h *hub) registerClient(client *client) {
	h.clients[client.id] = client
}

func (h *hub) unregisterClient(client *client) {
	if _, exist := h.clients[client.id]; !exist {
		return
	}

	delete(h.clients, client.id)
	close(client.send)
}
