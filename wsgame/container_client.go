package main

import "sync"

type containerClient struct {
	dict map[*client]bool // TODO: use a Tries or similar
	rw   *sync.RWMutex
}

func newContainerClient() *containerClient {
	return &containerClient{
		dict: make(map[*client]bool),
		rw:   new(sync.RWMutex),
	}
}

func (ctc *containerClient) add(c *client) {
	ctc.dict[c] = true
}

func (ctc *containerClient) remove(c *client) {
	if _, contains := ctc.dict[c]; contains {
		delete(ctc.dict, c)
		close(c.send)
	}
}
