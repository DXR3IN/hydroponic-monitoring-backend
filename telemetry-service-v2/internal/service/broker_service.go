package service

import (
	"log" // Tambahkan log untuk debugging

	models "github.com/DXR3IN/telemetry-service-v2/internal/domain"
)

type ClientSubscription struct {
	Channel  chan *models.Telemetry
	DeviceID string
}

type Broker struct {
	Notifier chan *models.Telemetry

	Clients map[string][]chan *models.Telemetry

	NewClients     chan ClientSubscription
	ClosingClients chan ClientSubscription
}

func NewBroker() *Broker {
	b := &Broker{
		Notifier:       make(chan *models.Telemetry, 1),
		Clients:        make(map[string][]chan *models.Telemetry),
		NewClients:     make(chan ClientSubscription),
		ClosingClients: make(chan ClientSubscription),
	}
	go b.listen()
	return b
}

func (b *Broker) removeClient(deviceID string, clientChan chan *models.Telemetry) {
	clientChans, ok := b.Clients[deviceID]
	if !ok {
		return
	}

	for i, ch := range clientChans {
		if ch == clientChan {
			b.Clients[deviceID] = append(clientChans[:i], clientChans[i+1:]...)
			if len(b.Clients[deviceID]) == 0 {
				delete(b.Clients, deviceID)
			}
			return
		}
	}
}

func (b *Broker) listen() {
	for {
		select {
		case newClient := <-b.NewClients:
			b.Clients[newClient.DeviceID] = append(b.Clients[newClient.DeviceID], newClient.Channel)
			log.Printf("Broker: New client subscribed to DeviceID: %s. Total subs: %d", newClient.DeviceID, len(b.Clients[newClient.DeviceID]))

		case closingClient := <-b.ClosingClients:
			b.removeClient(closingClient.DeviceID, closingClient.Channel)
			log.Printf("Broker: Client unsubscribed from DeviceID: %s", closingClient.DeviceID)

		case data := <-b.Notifier:
			if clientChans, ok := b.Clients[data.DeviceID]; ok {
				for _, clientChan := range clientChans {
					select {
					case clientChan <- data:
					default:
						log.Println("Broker: Client channel full, dropping message.")
					}
				}
			}
		}
	}
}
