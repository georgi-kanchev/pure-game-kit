package multiplayer

import (
	"bufio"
	"fmt"
	"net"
	con "pure-game-kit/multiplayer/connection"
	"pure-game-kit/utility/text"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	myId         int
	conn         *net.Conn
	pingResponse chan struct{} // channel to receive the server connection pong to our ping
}

func NewClient(ip string, onMessage func(fromId, tag int, message string)) *Client {
	var conn, err = net.Dial("tcp", ip+":"+port)
	if err != nil {
		return nil
	}
	var client = &Client{conn: &conn, myId: -1, pingResponse: make(chan struct{}, 1)}
	onMessage(0, con.MeJoined, ip+":"+port)

	go func() {
		var scanner = bufio.NewScanner(conn)
		for scanner.Scan() {
			if client.conn == nil {
				return // we already left
			}

			var parts = strings.SplitN(scanner.Text(), divider, 3)
			if len(parts) != 3 {
				continue
			}

			var fromId, _ = strconv.ParseInt(parts[0], 10, 64)
			var tag, _ = strconv.ParseInt(parts[1], 10, 64)
			var msg = parts[2]

			if fromId == 0 && tag == pong {
				var myNewId, _ = strconv.ParseInt(msg, 10, 64)
				client.myId = int(myNewId)
				client.pingResponse <- struct{}{}
				continue // server responded to our ping with a pong, we are still connected (all good)
			} else if fromId == 0 && tag == con.ServerStopped {
				client.Leave() // server stopped, so we leave (to not drop)
			}

			onMessage(int(fromId), int(tag), msg)
		}
	}()

	go func() {
		for range time.Tick(time.Second) {
			if client.conn == nil {
				return // we already left
			}

			client.send(true, 0, ping, "")

			select {
			case <-client.pingResponse: // pong received, we are still connected
			case <-time.After(time.Second): // no pong received within 1 second, dropping
				if client.conn == nil {
					return // we already left, do not send ~dropped
				}

				onMessage(0, con.MeDropped, "")
				client.Leave()
				return
			}
		}
	}()

	return client
}

//=================================================================

func (client *Client) Id() int {
	return client.myId
}

//=================================================================

func (client *Client) SendToServer(tag int, message string) {
	client.send(false, 0, tag, message)
}
func (client *Client) SendToAll(tag int, message string) {
	client.send(false, -1, tag, message)
}
func (client *Client) SendToClient(clientId, tag int, message string) {
	client.send(false, clientId, tag, message)
}

func (client *Client) Leave() {
	if client.conn != nil {
		(*client.conn).Close()
		client.conn = nil
		client.myId = -1
	}
}

//=================================================================
// private

func (client *Client) send(internally bool, toId, tag int, message string) {
	if (!internally && tag < 0) || client.conn == nil {
		return
	}

	fmt.Fprintln(*client.conn, text.New(toId, divider, tag, divider, message))
}
