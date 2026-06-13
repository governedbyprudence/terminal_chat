package messageStore

import (
	"bufio"
	"encoding/json"
	"net"
	"sync"
)

type Subscription struct {
	ch      chan string
	broker  *Broker
	channel string
}

func (s *Subscription) Channel() chan string {
	return s.ch
}

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

type wireMsg struct {
	Channel string `json:"c"`
	Payload string `json:"p"`
}

type Broker struct {
	subscribers map[string][]chan string
	mu          sync.RWMutex
	conn        net.Conn
	writer      *bufio.Writer
}

var (
	globalBroker *Broker
	once         sync.Once
)

func DefaultBroker() *Broker {
	once.Do(func() {
		globalBroker = &Broker{
			subscribers: make(map[string][]chan string),
		}
	})
	return globalBroker
}

func (b *Broker) Subscribe(channel string) *Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 100)
	b.subscribers[channel] = append(b.subscribers[channel], ch)
	return &Subscription{ch: ch, broker: b, channel: channel}
}

func (b *Broker) deliverLocal(channel string, payload string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subscribers[channel] {
		select {
		case ch <- payload:
		default:
		}
	}
}

func (b *Broker) Publish(channel string, payload string) {
	b.mu.RLock()
	writer := b.writer
	b.mu.RUnlock()
	if writer != nil {
		msg := wireMsg{Channel: channel, Payload: payload}
		data, _ := json.Marshal(msg)
		_, err := writer.WriteString(string(data) + "\n")
		if err != nil {
			return
		}
		writer.Flush()
		return
	}
	b.deliverLocal(channel, payload)
}

func (b *Broker) readLoop(reader *bufio.Reader) {
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

func (b *Broker) Connect(addr string) error {
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

var (
	serverClients   = make(map[net.Conn]*bufio.Writer)
	serverClientsMu sync.RWMutex
)

func StartServer(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleServerClient(conn)
		}
	}()
	return nil
}

func handleServerClient(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	serverClientsMu.Lock()
	serverClients[conn] = writer
	serverClientsMu.Unlock()

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		serverClientsMu.RLock()
		for _, w := range serverClients {
			if _, err := w.WriteString(line); err != nil {
				continue
			}
			w.Flush()
		}
		serverClientsMu.RUnlock()
	}

	serverClientsMu.Lock()
	delete(serverClients, conn)
	serverClientsMu.Unlock()
	conn.Close()
}
