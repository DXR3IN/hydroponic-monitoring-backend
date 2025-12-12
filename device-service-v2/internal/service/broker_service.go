package service

import models "github.com/DXR3IN/device-service-v2/internal/domain"

type Broker struct {
	Notifier       chan *models.Telemetry
	Clients        map[chan *models.Telemetry]bool
	NewClients     chan chan *models.Telemetry
	ClosingClients chan chan *models.Telemetry
}

func NewBroker() *Broker {
	b := &Broker{
		Notifier:       make(chan *models.Telemetry, 1),
		Clients:        make(map[chan *models.Telemetry]bool),
		NewClients:     make(chan chan *models.Telemetry),
		ClosingClients: make(chan chan *models.Telemetry),
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
