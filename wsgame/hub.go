package main

type hub struct {
	clients    *containerClient
	register   chan *client // Register requests from the clients.
	unregister chan *client // Unregister requests from clients.
}

func newHub() *hub {
	return &hub{
		register:   make(chan *client, 1),
		unregister: make(chan *client, 1),
		clients:    newContainerClient(),
	}
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients.add(c)
		case c := <-h.unregister:
			h.clients.remove(c)
		}
	}
}
