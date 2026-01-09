package ws

import (
	"github.com/google/uuid"
	"fmt"
)

type Room struct {
	ID      uuid.UUID         `json:"id"`
	Name    string	          `json:"name"`
	Clients map[int64]*Client `json:"clients"`
}

type Hub struct {
	Rooms        map[uuid.UUID]*Room
	RegisterCh   chan *Client
	UnregisterCh chan *Client
	BroadcastCh  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[uuid.UUID]*Room),
		RegisterCh: make(chan *Client),
		UnregisterCh: make(chan *Client),
		BroadcastCh: make(chan *Message, 16),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <- h.RegisterCh:
			println("registerch")
			if r, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
					fmt.Printf("User %s registered to room %s\n", cl.Username, cl.RoomID)
					h.BroadcastCh <- &Message {
						Content: fmt.Sprintf("%s has joined.", cl.Username),
						RoomID: cl.RoomID,
						Username: cl.Username,
					}
				}
			}
		case cl := <- h.UnregisterCh:
			println("unregisterch")
			if r, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := r.Clients[cl.ID]; ok {
					h.BroadcastCh <- &Message {
						Content: fmt.Sprintf("%s has left.", cl.Username),
						RoomID: cl.RoomID,
						Username: cl.Username,
					}
					delete(r.Clients, cl.ID)
					close(cl.MessageCh)
				}
			}
		case msg := <- h.BroadcastCh:
			println("broadcastch")
			if r, ok := h.Rooms[msg.RoomID]; ok {
				for _, cl := range r.Clients {
					cl.MessageCh <- msg	
				}
			}
		}
	}
}
