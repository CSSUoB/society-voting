package events

import (
	"bufio"
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
	Data  string
}

func (m *Message) ToSSE() []byte {
	var sb strings.Builder
	sb.WriteString("event: ")
	sb.WriteString(string(m.Topic))
	sb.WriteRune('\n')

	sp := strings.Split(m.Data, "\n")
	for i, x := range sp {
		sp[i] = "data: " + x
	}

	sb.WriteString(strings.Join(sp, "\n"))
	sb.WriteString("\n\n")

	return []byte(sb.String())
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

func SendEvent(topic Topic, data string) {
	masterLock.RLock()
	defer masterLock.RUnlock()

	msg := &Message{
		Topic: topic,
		Data:  data,
	}
	chans := channelsByTopic[topic]

	for _, ch := range chans {
		ch.Channel <- msg
	}
}

func AsStreamWriter(id int, receiver chan *Message) fasthttp.StreamWriter {
	slog.Debug("starting SSE streamwriter", "id", id)
	return func(w *bufio.Writer) {
		for msg := range receiver {
			_, _ = w.Write(msg.ToSSE())
			if err := w.Flush(); err != nil {
				// Client disconnected
				break
			}
		}
		CloseReceiver(id)
		slog.Debug("closing SSE connection", "id", id)
	}
}
