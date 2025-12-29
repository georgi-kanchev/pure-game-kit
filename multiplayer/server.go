package multiplayer

import (
	"bufio"
	"fmt"
	"net"
	con "pure-game-kit/multiplayer/connection"
	"pure-game-kit/utility/text"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	nextId   int64
	clients  map[int]*net.Conn
	mu       sync.Mutex
	listener net.Listener
}

func NewServer(onMessage func(fromId, toId, tag int, message string)) *Server {
	var listener, err = net.Listen("tcp", text.New(":", port))
	if err != nil {
		return nil
	}

	var server = &Server{clients: make(map[int]*net.Conn), listener: listener}
	onMessage(0, 0, con.ServerStarted, port)

	go func() {
		for {
			var conn, err = listener.Accept()
			if err != nil {
				continue
			}

			var id = int(atomic.AddInt64(&server.nextId, 1))
			if id == 0 {
				id = int(atomic.AddInt64(&server.nextId, 1))
			}

			server.mu.Lock()
			server.clients[id] = &conn
			server.mu.Unlock()

			onMessage(id, 0, con.ClientJoined, "")                     // notify self
			server.sendToClient(true, 0, id, pong, text.New(id))       // notify new client of their id
			server.sendToAllButOne(true, id, id, con.ClientJoined, "") // notify all clients but the joiner
			go server.handleClient(id, conn, onMessage)
		}
	}()

	return server
}

//=================================================================

func (s *Server) SendToClient(clientId, tag int, message string) {
	s.sendToClient(false, 0, clientId, tag, message)
}
func (s *Server) SendToAll(tag int, message string) {
	s.sendToAll(false, 0, tag, message)
}

func (s *Server) Stop() {
	// notify all clients we are stopping on our own terms (expectedly & willingly)
	s.sendToAll(true, 0, con.ServerStopped, "")

	if s.listener != nil {
		s.listener.Close() // no more joiners, we are stopping
		s.listener = nil
	}

	time.Sleep(time.Second) // give some time for all clients to receive the messages that we are stopping

	s.mu.Lock()
	for _, conn := range s.clients {
		(*conn).Close()
	}
	s.clients = make(map[int]*net.Conn)
	s.mu.Unlock()
}

// =================================================================
// private

const port, divider, ping, pong = "9000", "│", -8, -9

func (s *Server) handleClient(id int, conn net.Conn, onMessage func(int, int, int, string)) {
	defer func() {
		if s.listener == nil {
			return // we already stopped
		}

		conn.Close()
		s.mu.Lock()
		delete(s.clients, id)
		s.mu.Unlock()
		onMessage(id, 0, con.ClientLeft, "")      // notify self
		s.sendToAll(true, id, con.ClientLeft, "") // notify all clients
	}()

	var reader = bufio.NewScanner(conn)
	for reader.Scan() {
		var raw = reader.Text()
		var parts = strings.SplitN(raw, divider, 3)
		if len(parts) != 3 {
			continue
		}

		var toId, _ = strconv.ParseInt(parts[0], 10, 64)
		var tag, _ = strconv.ParseInt(parts[1], 10, 64)
		var msg = parts[2]

		if tag == ping { // client pinged us to see if they are connected, so pong back
			s.sendToClient(true, 0, id, pong, text.New(id)) // along with their id
			continue
		}

		onMessage(id, int(toId), int(tag), msg)

		if tag < 0 {
			continue // skipping relay of internal messages
		}

		if toId == -1 { // relaying non-internal messages
			s.sendToAllButOne(true, id, id, int(tag), msg) // do not sent back to sender
		} else if toId > 0 {
			s.sendToClient(true, id, int(toId), int(tag), msg)
		}
	}
}

func (s *Server) sendToAllButOne(internally bool, fromId, butId, tag int, message string) {
	if (!internally && tag < 0) || s.listener == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for id, conn := range s.clients {
		if id != butId {
			fmt.Fprintln(*conn, text.New(fromId, divider, tag, divider, message))
		}
	}
}
func (s *Server) sendToClient(internally bool, fromId, toId, tag int, message string) {
	if (!internally && tag < 0) || s.listener == nil {
		return
	}

	s.mu.Lock()
	var conn = s.clients[toId]
	s.mu.Unlock()
	if conn != nil {
		fmt.Fprintln(*conn, text.New(fromId, divider, tag, divider, message))
	}
}
func (s *Server) sendToAll(internally bool, fromId, tag int, message string) {
	s.sendToAllButOne(internally, fromId, -1, tag, message) // no one has id -1 so it means all
}
