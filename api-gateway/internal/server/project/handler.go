package project

import (
	chatpb "api-gateway/pkg/api/chat_v1"
	projectpb "api-gateway/pkg/api/project_v1"
)

type Handler struct {
	projectAPIClient projectpb.ProjectServiceClient
	chatApiClient    chatpb.ChatServiceClient
}

func NewHandler(projectAPIClient projectpb.ProjectServiceClient, chatApiClient chatpb.ChatServiceClient) *Handler {
	return &Handler{projectAPIClient: projectAPIClient, chatApiClient: chatApiClient}
}
