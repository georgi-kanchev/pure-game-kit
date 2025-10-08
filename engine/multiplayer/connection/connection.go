package connection

const (
	Tag = "▶"

	ServerStarted = Tag + "started"
	ServerStopped = Tag + "stopped"

	ClientJoined  = Tag + "joined"
	ClientLeft    = Tag + "left"
	ClientDropped = Tag + "dropped"
)
