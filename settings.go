package main

import "time"

const (
	// ClearPer - Сhecking for messages to send
	// using in clients/CleanOffConn()
	ClearPer = 2000

	// TimeForAuth - Time  For Auth
	// using in connections/func (c *Connections) Add
	TimeForAuth = 10000


	// MaxConnections - max connections clients
	// using in:
	// clients/CleanOffConn()
	// server/Connections
	MaxConnections = 500


	// Maximum message size allowed from peer.
	// using in:
	// clients/Add
	maxMessageSize = 512

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)
