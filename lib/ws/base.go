package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

var WebsocketServer = websocketServer{
	Clients:    make([]*WebsocketConnection, 0),
	Register:   make(chan *WebsocketConnection, 128),
	UnRegister: make(chan *WebsocketConnection, 128),
}

type websocketServer struct {
	Clients              []*WebsocketConnection
	Register, UnRegister chan *WebsocketConnection
	Locker               sync.Mutex
}

func (s *websocketServer) Start() {
	for {
		select {
		case client := <-s.Register:
			s.Locker.Lock()
			s.Clients = append(s.Clients, client)
			s.Locker.Unlock()
		case client := <-s.UnRegister:
			s.Locker.Lock()
			for i := 0; i < len(s.Clients); i++ {
				if s.Clients[i].Key == client.Key {
					s.Clients = append(s.Clients[0:i], s.Clients[i+1:]...)
				}
			}
			s.Locker.Unlock()
		}
	}
}

func (s *websocketServer) RegisterConn(key string, unsafeConn *websocket.Conn) {
	conn := &threadSafeConn{unsafeConn, sync.Mutex{}}
	c := &WebsocketConnection{
		Key: key, Conn: conn,
	}
	s.Register <- c
}

func (s *websocketServer) FindClient(key string) *WebsocketConnection {
	for _, c := range s.Clients {
		if c.Key == key {
			return c
		}
	}
	return nil
}

func (s *websocketServer) Send(msg interface{}) {
	for _, c := range s.Clients {
		c.Conn.WriteJSON(msg)
	}
}

func (s *websocketServer) SendTo(msg interface{}, key string) {
	for _, c := range s.Clients {
		if c.Key == key {
			c.Conn.WriteJSON(msg)
		}
	}
}

func (s *websocketServer) UnRegisterConn(key string) {
	c := s.FindClient(key)
	s.UnRegister <- c
}

type WebsocketConnection struct {
	Key  string
	Conn *threadSafeConn
}

type threadSafeConn struct {
	*websocket.Conn
	sync.Mutex
}

func (t *threadSafeConn) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()
	return t.Conn.WriteJSON(v)
}
