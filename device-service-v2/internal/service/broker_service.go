package service

type Broker struct {
	Notifier       chan interface{}
	Clients        map[chan interface{}]bool
	NewClients     chan chan interface{}
	ClosingClients chan chan interface{}
}

func NewBroker() *Broker {
	b := &Broker{
		Notifier:       make(chan interface{}, 1),
		Clients:        make(map[chan interface{}]bool),
		NewClients:     make(chan chan interface{}),
		ClosingClients: make(chan chan interface{}),
	}
	go b.listen()
	return b
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.NewClients:
			b.Clients[s] = true
		case s := <-b.ClosingClients:
			delete(b.Clients, s)
		case data := <-b.Notifier:
			for clientChan := range b.Clients {
				clientChan <- data
			}
		}
	}
}
