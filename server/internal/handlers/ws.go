package handlers

import (
	"net/http"
	"drawing-api/internal/api"
	"drawing-api/internal/ws"
	"drawing-api/internal/service"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
	"fmt"
)

type WsHandler struct {
	Hub *ws.Hub	
}

func NewWsHandler(hub *ws.Hub) *WsHandler {
	return &WsHandler{
		Hub: hub,
	}
}

type CreateRoomReq struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (h* WsHandler) CreateRoom(w http.ResponseWriter, r *http.Request) error {
	var req CreateRoomReq
	if err := api.ParseJSON(r, &req); err != nil {
		return err
	}

	h.Hub.Rooms[req.ID] = &ws.Room{
		ID: req.ID,
		Name: req.Name,
		Clients: make(map[int64]*ws.Client),
	}

	return api.WriteJSON(w, http.StatusOK, h.Hub.Rooms[req.ID])
}

var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WsHandler) JoinRoom(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	claims := ctx.Value(ClaimsKey).(*service.JwtClaims)
	fmt.Printf("id %v\nusername: %v\n", claims.ID, claims.Username)

	roomUUID, err := uuid.Parse(r.PathValue("roomID"))
	if err != nil {
		return api.NewAPIError(err.Error(), http.StatusBadRequest)
	}
	
	cl := &ws.Client {
		Conn: conn,
		MessageCh: make(chan *ws.Message, 16),
		ID: claims.ID,
		RoomID: roomUUID,
		Username: claims.Username,
	}

	h.Hub.RegisterCh <- cl

	go cl.WriteMessage()
	cl.ReadMessage(h.Hub)
 	return nil
}
