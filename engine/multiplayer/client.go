package multiplayer

import (
	"bufio"
	"fmt"
	"net"
	con "pure-kit/engine/multiplayer/connection"
	"pure-kit/engine/utility/text"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	myId         int64
	conn         *net.Conn
	pingResponse chan struct{} // channel to receive the server connection pong to our ping
}

func NewClient(ip string, onMessage func(fromId int64, message string)) *Client {
	var conn, err = net.Dial("tcp", ip+":"+port)
	if err != nil {
		return nil
	}
	var client = &Client{conn: &conn, myId: -1, pingResponse: make(chan struct{}, 1)}
	onMessage(0, con.ClientJoined+" "+ip+":"+port)

	go func() {
		var scanner = bufio.NewScanner(conn)
		for scanner.Scan() {
			if client.conn == nil {
				return // we already left
			}

			var parts = strings.SplitN(scanner.Text(), divider, 2)
			if len(parts) != 2 {
				continue
			}

			if parts[0] == "0" && text.Contains(parts[1], pong) {
				client.myId, _ = strconv.ParseInt(strings.SplitN(parts[1], " ", 2)[1], 10, 64)
				client.pingResponse <- struct{}{}
				continue // server responded to our ping with a pong, we are still connected (all good)
			}
			if parts[0] == "0" && parts[1] == con.ServerStopped {
				client.Leave() // server stopped, so we leave (to not drop)
			}

			var fromId, _ = strconv.ParseInt(parts[0], 10, 64)
			onMessage(fromId, parts[1])
		}
	}()

	go func() {
		for range time.Tick(time.Second) {
			if client.conn == nil {
				return // we already left
			}

			client.SendToServer(ping)

			select {
			case <-client.pingResponse: // pong received, we are still connected
			case <-time.After(time.Second): // no pong received within 1 second, dropping
				if client.conn == nil {
					return // we already left, do not send ~dropped
				}

				onMessage(0, con.ClientDropped)
				client.Leave()
				return
			}
		}
	}()

	return client
}

func (client *Client) Id() int64 {
	return client.myId
}

func (client *Client) SendToServer(message string) {
	fmt.Fprintln(*client.conn, text.New(0, divider, message))
}
func (client *Client) SendToAll(message string) {
	fmt.Fprintln(*client.conn, text.New(-1, divider, message))
}
func (client *Client) SendToClient(clientId uint64, message string) {
	fmt.Fprintln(*client.conn, text.New(clientId, divider, message))
}

func (client *Client) Leave() {
	if client.conn != nil {
		(*client.conn).Close()
		client.conn = nil
		client.myId = -1
	}
}
