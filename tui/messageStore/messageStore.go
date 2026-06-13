package messageStore

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"sync"
)

const (
	ChannelChat     = "chat"
	ChannelPresence = "presence"
)

type Subscription struct {
	ch      chan string
	broker  *broker
	channel string
}

func (s *Subscription) Channel() chan string { return s.ch }

func (s *Subscription) Close() {
	s.broker.mu.Lock()
	defer s.broker.mu.Unlock()
	subs := s.broker.subscribers[s.channel]
	for i, sub := range subs {
		if sub == s.ch {
			s.broker.subscribers[s.channel] = append(subs[:i], subs[i+1:]...)
			close(s.ch)
			return
		}
	}
}

type Broker interface {
	Subscribe(channel string) *Subscription
	Publish(channel, payload string)
	Connect(addr string) error
	StartServer(addr string) error
	StopServer()
}

type wireMsg struct {
	Channel string `json:"channel"`
	Payload string `json:"payload"`
}

type broker struct {
	subscribers   map[string][]chan string
	mu            sync.RWMutex
	conn          net.Conn
	writer        *bufio.Writer
	listener      net.Listener
	serverClients map[net.Conn]*bufio.Writer
	serverMu      sync.RWMutex
}

func New() Broker {
	return &broker{
		subscribers:   make(map[string][]chan string),
		serverClients: make(map[net.Conn]*bufio.Writer),
	}
}

func (b *broker) Subscribe(channel string) *Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 100)
	b.subscribers[channel] = append(b.subscribers[channel], ch)
	return &Subscription{ch: ch, broker: b, channel: channel}
}

func (b *broker) deliverLocal(channel string, payload string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subscribers[channel] {
		select {
		case ch <- payload:
		default:
		}
	}
}

func (b *broker) Publish(channel string, payload string) {
	b.mu.Lock()
	if b.writer != nil {
		msg := wireMsg{Channel: channel, Payload: payload}
		data, err := json.Marshal(msg)
		if err != nil {
			b.mu.Unlock()
			log.Printf("json.Marshal error: %v", err)
			return
		}
		_, err = b.writer.WriteString(string(data) + "\n")
		if err != nil {
			b.mu.Unlock()
			log.Printf("write error: %v", err)
			return
		}
		b.writer.Flush()
		b.mu.Unlock()
		return
	}
	b.mu.Unlock()
	b.deliverLocal(channel, payload)
}

func (b *broker) readLoop(reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		var msg wireMsg
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}
		b.deliverLocal(msg.Channel, msg.Payload)
	}
	b.mu.Lock()
	if b.conn != nil {
		b.conn.Close()
		b.conn = nil
		b.writer = nil
	}
	b.mu.Unlock()
}

func (b *broker) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	b.mu.Lock()
	b.conn = conn
	b.writer = bufio.NewWriter(conn)
	b.mu.Unlock()
	go b.readLoop(bufio.NewReader(conn))
	return nil
}

func (b *broker) StartServer(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	b.mu.Lock()
	b.listener = listener
	b.mu.Unlock()
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go b.handleServerClient(conn)
		}
	}()
	return nil
}

func (b *broker) StopServer() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.listener != nil {
		b.listener.Close()
		b.listener = nil
	}
}

func (b *broker) handleServerClient(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	b.serverMu.Lock()
	b.serverClients[conn] = writer
	b.serverMu.Unlock()

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		b.serverMu.RLock()
		clients := make([]*bufio.Writer, 0, len(b.serverClients))
		for _, w := range b.serverClients {
			clients = append(clients, w)
		}
		b.serverMu.RUnlock()

		for _, w := range clients {
			if _, err := w.WriteString(line); err != nil {
				continue
			}
			w.Flush()
		}
	}

	b.serverMu.Lock()
	delete(b.serverClients, conn)
	b.serverMu.Unlock()
	conn.Close()
}
