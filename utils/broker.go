package utils

type subscription[MessageType any] struct {
	sub  chan<- MessageType
	done <-chan struct{}
}

type Broker[MessageType any] struct {
	sub  chan subscription[MessageType]
	pub  chan MessageType
	done chan struct{}
}

func (b *Broker[MessageType]) Subscribe(done <-chan struct{}) chan MessageType {
	sub := make(chan MessageType, 5)
	b.sub <- subscription[MessageType]{sub: sub, done: done}
	return sub
}

func (b *Broker[MessageType]) Publish(msg MessageType) {
	b.pub <- msg
}

func (b *Broker[MessageType]) Shutdown() {
	close(b.done)
}

func (b *Broker[MessageType]) handle() {
	subs := make(map[chan<- MessageType]struct{})
	unsub := make(chan chan<- MessageType)

	for {
		select {
		case <-b.done:
			for sub := range subs {
				close(sub)
			}
			return

		case s := <-b.sub:
			subs[s.sub] = struct{}{}

			go func() {
				<-s.done
				unsub <- s.sub
			}()

		case sub := <-unsub:
			delete(subs, sub)
			close(sub)

		case msg := <-b.pub:
			for sub := range subs {
				// use non-blocking send to protect the broker
				select {
				case sub <- msg:
				default:
				}
			}
		}
	}
}

func NewBroker[MessageType any]() *Broker[MessageType] {
	var b Broker[MessageType]
	b.sub = make(chan subscription[MessageType])
	b.pub = make(chan MessageType, 1)
	b.done = make(chan struct{})

	go b.handle()

	return &b
}
