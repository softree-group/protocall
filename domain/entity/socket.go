package entity

const (
	SocketEventLeave       = "leave"
	SocketEventConnection  = "connection"
	SocketEventConnected   = "connected"
	SocketEventStartRecord = "start_record"
)

type SocketMessage map[string]interface{}
