package chat

import (
	ws "api-gateway/internal/server/chat/ws"
	chatpb "api-gateway/pkg/api/chat_v1"
)

type Handler struct {
	hub           *ws.Hub
	chatAPIClient chatpb.ChatServiceClient
}

func NewHandler(chatAPIClient chatpb.ChatServiceClient, chatService ws.ChatService, hub *ws.Hub) *Handler {
	return &Handler{
		hub:           hub,
		chatAPIClient: chatAPIClient,
	}
}
