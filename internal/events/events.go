package events

import (
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

type channelWithID struct {
	ID      int
	Channel chan *Message
}

func NewReceiver(topics ...Topic) chan *Message {
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

	return ch
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
