package events

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log/slog"
	"strings"
	"sync"
)

var (
	masterLock        = new(sync.RWMutex)
	channelsByTopic   = make(map[Topic][]*channelWithID)
	topicsByChannelID = make(map[int][]Topic)
	idCounter         = 0
)

type Topic string

type Message struct {
	Topic Topic
	Data  any
}

func (m *Message) ToSSE() ([]byte, error) {
	var sb strings.Builder
	sb.WriteString("event: ")
	sb.WriteString(string(m.Topic))
	sb.WriteRune('\n')

	jdat, err := json.Marshal(m.Data)
	if err != nil {
		return nil, err
	}

	sb.WriteString("data: ")
	sb.Write(jdat)
	sb.WriteString("\n\n")

	return []byte(sb.String()), nil
}

type channelWithID struct {
	ID      int
	Channel chan *Message
}

func NewReceiver(topics ...Topic) (int, chan *Message) {
	masterLock.Lock()
	defer masterLock.Unlock()

	ch := make(chan *Message, 32)

	id := idCounter
	idCounter += 1

	chID := &channelWithID{
		ID:      id,
		Channel: ch,
	}

	for _, topic := range topics {
		channelsByTopic[topic] = append(channelsByTopic[topic], chID)
	}

	topicsByChannelID[id] = topics

	return id, ch
}

func CloseReceiver(id int) {
	masterLock.Lock()
	defer masterLock.Unlock()

	topics, found := topicsByChannelID[id]

	if !found {
		return
	}

	delete(topicsByChannelID, id)

	var hasClosed bool

	for _, topic := range topics {
		chans := channelsByTopic[topic]
		if len(chans) == 0 {
			continue
		}
		n := 0
		for _, ch := range chans {
			if ch.ID != id {
				chans[n] = ch
				n += 1
			} else if !hasClosed {
				close(ch.Channel)
				hasClosed = true
			}
		}
		channelsByTopic[topic] = chans[:n]
	}
}

func SendEvent(topic Topic, data any) {
	masterLock.RLock()
	defer masterLock.RUnlock()

	msg := &Message{
		Topic: topic,
		Data:  data,
	}
	chans := channelsByTopic[topic]

	for _, ch := range chans {
		m := *msg
		ch.Channel <- &m
	}
}

func AsStreamWriter(id int, receiver chan *Message) fasthttp.StreamWriter {
	return func(w *bufio.Writer) {
		for msg := range receiver {
			sseData, err := msg.ToSSE()
			if err != nil {
				slog.Error("SSE error", "error", fmt.Errorf("failed to generate SSE event from message: %w", err))
				break
			}
			_, _ = w.Write(sseData)
			if err := w.Flush(); err != nil {
				// Client disconnected
				break
			}
		}
		CloseReceiver(id)
	}
}
