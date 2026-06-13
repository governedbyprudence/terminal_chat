package messageStore

import (
	"testing"
	"time"
)

func TestBroker_PublishSubscribe(t *testing.T) {
	b := New()
	sub := b.Subscribe("chat")
	defer sub.Close()

	b.Publish("chat", "hello")
	select {
	case msg := <-sub.Channel():
		if msg != "hello" {
			t.Errorf("expected 'hello', got %q", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestBroker_MultipleSubscribers(t *testing.T) {
	b := New()
	sub1 := b.Subscribe("chat")
	sub2 := b.Subscribe("chat")
	defer sub1.Close()
	defer sub2.Close()

	b.Publish("chat", "broadcast")

	for i, sub := range []*Subscription{sub1, sub2} {
		select {
		case <-sub.Channel():
		case <-time.After(time.Second):
			t.Fatalf("subscriber %d did not receive message", i+1)
		}
	}
}

func TestBroker_ChannelIsolation(t *testing.T) {
	b := New()
	chatSub := b.Subscribe("chat")
	presSub := b.Subscribe("presence")
	defer chatSub.Close()
	defer presSub.Close()

	b.Publish("chat", "chat msg")
	b.Publish("presence", "presence msg")

	select {
	case msg := <-chatSub.Channel():
		if msg != "chat msg" {
			t.Errorf("expected 'chat msg', got %q", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("chat subscriber timeout")
	}

	select {
	case msg := <-presSub.Channel():
		if msg != "presence msg" {
			t.Errorf("expected 'presence msg', got %q", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("presence subscriber timeout")
	}
}

func TestBroker_SubscribeClose(t *testing.T) {
	b := New()
	sub := b.Subscribe("chat")
	sub.Close()

	b.Publish("chat", "after close")
	select {
	case _, ok := <-sub.Channel():
		if ok {
			t.Error("channel should be closed after Close()")
		}
	default:
	}
}

func TestBroker_BufferedChannelNonBlocking(t *testing.T) {
	b := New()
	sub := b.Subscribe("chat")
	defer sub.Close()

	for i := 0; i < 200; i++ {
		b.Publish("chat", "msg")
	}

	ch := sub.Channel()
	count := 0
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				t.Fatal("channel closed unexpectedly")
			}
			count++
		default:
			goto done
		}
	}
done:
	if count < 1 {
		t.Error("should have received some messages")
	}
}
