// Contains constants, describing the type of a connection message.
package connection

const (
	ServerStarted, ServerStopped = -1, -2
	MeJoined, MeDropped          = -3, -4
	ClientJoined, ClientLeft     = -5, -6
	ClientDropped                = -7 // losing connection or no longer being able to reach server
)
