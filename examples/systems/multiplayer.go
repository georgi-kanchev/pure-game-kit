package example

import (
	"fmt"
	"pure-kit/engine/multiplayer"
	"pure-kit/engine/utility/text"
	"time"
)

func Multiplayer() {
	server := multiplayer.NewServer(func(fromId, toId int64, message string) {
		println(text.New("server [", fromId, " -> ", toId, "]", message))
	})

	client1 := multiplayer.NewClient("127.0.0.1", func(fromId int64, message string) {
		println(text.New("client1 [", fromId, "]", message))
	})

	time.Sleep(time.Second)
	fmt.Printf("client1.Id(): %v\n", client1.Id())
	client1.SendToAll("Hello everyone!")
	time.Sleep(time.Second)
	client2 := multiplayer.NewClient("127.0.0.1", func(fromId int64, message string) {
		println(text.New("client2 [", fromId, "]", message))
	})

	client2.SendToClient(1, "Hey, I just joined!")
	time.Sleep(time.Second)
	client1.Leave()
	time.Sleep(time.Second)
	fmt.Printf("client1.Id(): %v\n", client1.Id())
	client2.SendToAll(":(")
	time.Sleep(time.Second)
	server.Stop()

	select {}
}
