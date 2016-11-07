package main

import "sync"

type containerClient struct {
	dict map[string]*client
	rw   *sync.RWMutex
}

func newContainerClient() *containerClient {
	return &containerClient{
		dict: make(map[string]*client),
		rw:   new(sync.RWMutex),
	}
}

func (ctc *containerClient) add(key string, c *client) {
	if c == nil {
		panic("client is nil, you deserve a nasty panic")
	}

	ctc.rw.Lock()
	ctc.dict[key] = c
	ctc.rw.Unlock()
}

func (ctc *containerClient) remove(key string) {
	ctc.rw.RLock()
	c, contains := ctc.dict[key]
	ctc.rw.RUnlock()

	if contains {
		ctc.rw.Lock()
		delete(ctc.dict, key)
		close(c.send)
		ctc.rw.Unlock()
	}
}
