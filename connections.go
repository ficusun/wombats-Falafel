package main

import (
	"sync"
)

// ClientHub Client Hub
type ClientHub struct {
	Connections map[int]*Client
	ID          ID
}

// getID () int - return last client ID
// ch.ID++
// return ch.ID
func (ch *ClientHub) getID() int {
	return ch.ID.Get()
}

// AddClient
func (ch *ClientHub) AddClient(client *Client) bool {
	ID := ch.getID()

	_, ok := ch.Connections[ID]
	if !ok {
		ch.Connections[ID] = client
		ch.Connections[ID].start()
		ch.Connections[ID].ID = ID
		ch.Connections[ID].Status = true

		return true
	}

	return false
}

// DelClientByID - DelClientByID(ID int) bool
func (ch *ClientHub) DelClientByID(ID int) bool {

	delete(ch.Connections, ID)
	_, ok := ch.Connections[ID]
	if ok == false {
		return true
	}
	return false
}

// GetClientByID(ID int) (bool, Client)
func (ch *ClientHub) GetClientByID(ID int) (bool, *Client) {

	client, ok := ch.Connections[ID]
	return ok, client
}

// struct ID ===============================================================<
// 	sync.RWMutex
//	ID int
type ID struct {
	sync.RWMutex
	ID int
}

// Get () int
// safe with Mutex
func (id *ID) Get() int {
	// block for read
	id.RLock()

	// change ID
	id.ID++

	// unblock
	defer id.RUnlock()

	// return result
	return id.ID
}

// ========= end struct ID =====================================>
