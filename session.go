package gott

import (
	"time"
)

type Session struct {
	client       *Client
	clean        bool
	Id           string
	MessageStore *MessageStore
}

func NewSession(client *Client, cleanFlag string) *Session {
	return &Session{
		client:       client,
		clean:        cleanFlag == "1",
		MessageStore: NewMessageStore(),
		Id:           client.ClientId,
	}
}

func (s *Session) Load() error {
	start := time.Now()
	err := GOTT.SessionStore.Get(s.Id, s)
	end := time.Since(start)
	Log("[BENCHMARK]", "session load took:", end)
	return err
}

func (s *Session) Update(value map[string]interface{}) error {
	return GOTT.SessionStore.Set(s.Id, value)
}

func (s *Session) Put() error {
	start := time.Now()
	err := GOTT.SessionStore.Set(s.Id, s)
	end := time.Since(start)
	Log("[BENCHMARK]", "session put took:", end)
	return err
}

func (s *Session) Client() *Client {
	return s.client
}

func (s *Session) Clean() bool {
	return s.clean
}

func (s *Session) StoreMessage(packetId uint16, msg *ClientMessage) {
	s.MessageStore.Store(packetId, msg)
	_ = s.Put()
}