package example

import (
	"pure-game-kit/multiplayer"
	"pure-game-kit/multiplayer/connection"
	"pure-game-kit/utility/text"
	"time"
)

func Multiplayer() {
	var server = multiplayer.NewServer(func(fromId, toId, tag int, message string) {
		switch tag {
		case connection.ServerStarted:
			println(text.New("server: started on port ", message))
		case connection.ClientJoined:
			println(text.New("server: client", fromId, " joined"))
		case connection.ClientLeft:
			println(text.New("server: client", fromId, " left"))
		case connection.ClientDropped:
			println(text.New("server: client", fromId, " lost connection"))
		}

		if tag == 0 {
			println(text.New("client", fromId, " says: ", message))
		}
	})

	var client1 = multiplayer.NewClient("127.0.0.1", func(fromId, tag int, message string) {
		switch tag {
		case connection.ServerStopped:
			println(text.New("client1: server stopped"))
		case connection.MeJoined:
			println(text.New("client1: i joined ", message))
		case connection.MeDropped:
			println(text.New("client1: i lost connection"))
		case connection.ClientJoined:
			println(text.New("client1: client", fromId, " joined"))
		case connection.ClientLeft:
			println(text.New("client1: client", fromId, " left"))
		case connection.ClientDropped:
			println(text.New("client1: client", fromId, " lost connection"))
		}
	})

	time.Sleep(time.Second)

	client1.SendToAll(0, "Hello everyone!")
	time.Sleep(time.Second)

	var client2 = multiplayer.NewClient("127.0.0.1", func(fromId, tag int, message string) {
		switch tag {
		case connection.ServerStopped:
			println(text.New("client2: server stopped"))
		case connection.MeJoined:
			println(text.New("client2: i joined ", message))
		case connection.MeDropped:
			println(text.New("client2: i lost connection"))
		case connection.ClientJoined:
			println(text.New("client2: client", fromId, " joined"))
		case connection.ClientLeft:
			println(text.New("client2: client", fromId, " left"))
		case connection.ClientDropped:
			println(text.New("client2: client", fromId, " lost connection"))
		}
	})

	client2.SendToClient(1, 0, "Hey, I just joined!")
	time.Sleep(time.Second)

	client1.Leave()
	time.Sleep(time.Second)

	client2.SendToAll(0, ":(")
	time.Sleep(time.Second)

	server.Stop()

	select {}
}
