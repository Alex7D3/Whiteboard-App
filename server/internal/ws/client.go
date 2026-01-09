package ws

import (
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
	"log"
)

type Message struct {
	ID       uuid.UUID `json:"id"`
	RoomID   uuid.UUID `json:"room_id"`
	Content  string    `json:"content"`
	Username string    `json:"username"`
}

type Client struct {
	ID        int64     `json:"id"`
	RoomID    uuid.UUID `json:"room_id"`
	Username  string    `json:"username"`
	Conn      *websocket.Conn
	MessageCh chan *Message
}

type MessageType int

const (
	Chat = iota
	Info
	Blob
)

func (cl *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.UnregisterCh <- cl
		cl.Conn.Close()
	}()

	for {
		mType, m, err := cl.Conn.ReadMessage()
		if err != nil && websocket.IsUnexpectedCloseError(
			err,
			websocket.CloseGoingAway,
			websocket.CloseAbnormalClosure,
		) {
			log.Printf("error %v", err)
			break	
		}
		switch (mType) {
		case websocket.TextMessage:
			hub.BroadcastCh <- &Message{
				Content: string(m),
				RoomID: cl.RoomID,
				Username: cl.Username,
			}
		case websocket.BinaryMessage: 		
		}
	}
}

func (cl *Client) WriteMessage() {
	defer cl.Conn.Close()
	for {
		msg, ok := <- cl.MessageCh
		if !ok {
			return
		}
		cl.Conn.WriteJSON(msg)
	}
}
