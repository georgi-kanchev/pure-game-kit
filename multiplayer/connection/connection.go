package connection

const (
	ServerStarted = -1
	ServerStopped = -2

	MeJoined  = -3
	MeDropped = -4

	ClientJoined  = -5
	ClientLeft    = -6
	ClientDropped = -7 // losing connection or no longer being able to reach server
)
