package utils

type Broker[T any] struct {
	done       chan struct{}
	pub        chan T
	sub, unsub chan (chan<- T)
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		done:  make(chan struct{}),
		pub:   make(chan T, 1),
		sub:   make(chan (chan<- T), 1),
		unsub: make(chan (chan<- T), 1),
	}
}

func (b *Broker[T]) Start() {
	subs := make(map[chan<- T]struct{})
	for {
		select {
		case <-b.done:
			for sub := range subs {
				close(sub)
			}
			return
		case sub := <-b.sub:
			subs[sub] = struct{}{}
		case sub := <-b.unsub:
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

func (b *Broker[T]) Close() {
	close(b.done)
}

func (b *Broker[T]) Subscribe() chan T {
	sub := make(chan T, 5)
	b.sub <- sub
	return sub
}

func (b *Broker[T]) Unsubscribe(sub chan T) {
	b.unsub <- sub
}

func (b *Broker[T]) Publish(msg T) {
	b.pub <- msg
}
