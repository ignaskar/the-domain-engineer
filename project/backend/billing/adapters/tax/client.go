package tax

import (
	"github.com/ThreeDotsLabs/the-domain-engineer/clients"
)

type Client struct {
	clients *clients.Clients
}

func NewClient(clients *clients.Clients) *Client {
	if clients == nil {
		panic("nil clients")
	}
	return &Client{clients: clients}
}
