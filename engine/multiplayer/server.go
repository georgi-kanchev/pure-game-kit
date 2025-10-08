package multiplayer

import (
	"bufio"
	"fmt"
	"net"
	con "pure-kit/engine/multiplayer/connection"
	"pure-kit/engine/utility/text"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	nextId   int64
	clients  map[int64]*net.Conn
	mu       sync.Mutex
	listener net.Listener
}

func NewServer(onMessage func(fromId, toId int64, message string)) *Server {
	var listener, err = net.Listen("tcp", text.New(":", port))
	if err != nil {
		return nil
	}

	var server = &Server{clients: make(map[int64]*net.Conn), listener: listener}
	onMessage(0, 0, con.ServerStarted+" "+port)

	go func() {
		for {
			var conn, err = listener.Accept()
			if err != nil {
				continue
			}

			var id = atomic.AddInt64(&server.nextId, 1)
			if id == 0 {
				id = atomic.AddInt64(&server.nextId, 1)
			}

			server.mu.Lock()
			server.clients[id] = &conn
			server.mu.Unlock()

			onMessage(id, 0, con.ClientJoined)                                  // notify self
			server.SendToClient(id, text.New("0", divider, pong, " ", id))      // notify new client of their id
			server.sendToAllButOne(id, text.New(id, divider, con.ClientJoined)) // notify all clients but the joiner
			go server.handleClient(id, conn, onMessage)
		}
	}()

	return server
}

func (server *Server) SendToClient(clientId int64, message string) {
	server.mu.Lock()
	var conn = server.clients[clientId]
	server.mu.Unlock()
	if conn != nil {
		fmt.Fprintln(*conn, message)
	}
}
func (server *Server) SendToAll(message string) {
	server.sendToAllButOne(-1, message) // no one has id -1 so it means all
}

func (server *Server) Stop() {
	if server.listener != nil {
		server.listener.Close() // no more joiners, we are stopping
		server.listener = nil
	}

	server.SendToAll("0" + divider + con.ServerStopped)
	time.Sleep(time.Second) // notify all clients we are stopping on our own terms (expectedly & willingly)

	server.mu.Lock()
	for _, conn := range server.clients {
		(*conn).Close()
	}
	server.clients = make(map[int64]*net.Conn)
	server.mu.Unlock()
}

// =================================================================
// private

const (
	port    = "9000"
	ping    = con.Tag + "ping"
	pong    = con.Tag + "pong"
	divider = "│"
)

func (server *Server) sendToAllButOne(butId int64, message string) {
	server.mu.Lock()
	defer server.mu.Unlock()
	for id, conn := range server.clients {
		if id != butId {
			fmt.Fprintln(*conn, message)
		}
	}
}

func (server *Server) handleClient(id int64, conn net.Conn, onMessage func(int64, int64, string)) {
	defer func() {
		if server.listener == nil {
			return // we already stopped
		}

		conn.Close()
		server.mu.Lock()
		delete(server.clients, id)
		server.mu.Unlock()
		onMessage(id, 0, con.ClientLeft)                        // notify self
		server.SendToAll(text.New(id, divider, con.ClientLeft)) // notify all clients
	}()

	var reader = bufio.NewScanner(conn)
	for reader.Scan() {
		var msg = reader.Text()
		var parts = strings.SplitN(msg, divider, 2)
		if len(parts) != 2 {
			continue
		}

		if parts[1] == ping { // client pinged us to see if they are connected, so pong back
			server.SendToClient(id, text.New("0", divider, pong, " ", id)) // send their id as well
			continue
		}

		switch parts[0] {
		case "0":
			onMessage(id, 0, parts[1])
		case "-1":
			onMessage(id, -1, parts[1])
			server.sendToAllButOne(id, text.New(id, divider, parts[1])) // do not sent back to sender
		default:
			var targetId, _ = strconv.ParseInt(parts[0], 10, 64)
			onMessage(id, targetId, parts[1])
			server.SendToClient(targetId, text.New(id, divider, parts[1]))
		}
	}
}
