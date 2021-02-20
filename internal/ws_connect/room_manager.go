package ws_connect

import (
	"sync"
)

type RoomManager struct {
	m    map[string]*Room
	lock *sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		m:    make(map[string]*Room),
		lock: new(sync.RWMutex),
	}
}

func (g *RoomManager) RemoveRoom(id string) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.m[id]; ok {
		delete(g.m, id)
	}

	return
}

func (g *RoomManager) PureFindRoom(id string) (*Room, bool) {

	g.lock.RLock()
	h, ok := g.m[id]
	g.lock.RUnlock()
	if ok {
		return h, true
	}

	return nil, false
}

func (g *RoomManager) FindRoom(id string) *Room {

	g.lock.RLock()
	h, ok := g.m[id]
	g.lock.RUnlock()
	if ok {
		return h
	}

	h = NewRoom(g, id)
	g.addRoom(id, h)

	return h
}

func (g *RoomManager) addRoom(id string, h *Room) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.m[id] = h
}
