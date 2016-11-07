package main

type HubClients struct {
	clients    *containerClient
	register   chan *client // Register requests from the clients.
	unregister chan *client // Unregister requests from clients.
}

func newHub() *HubClients {
	return &HubClients{
		register:   make(chan *client, 1),
		unregister: make(chan *client, 1),
		clients:    newContainerClient(),
	}
}

func (h *HubClients) run() {
	for {
		select {
		case <-h.register:
			// h.clients.add(c)
		case <-h.unregister:
			// h.clients.remove(c)
		}
	}
}
