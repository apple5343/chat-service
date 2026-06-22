package user

import (
	userpb "api-gateway/pkg/api/user_v1"
)

type Handler struct {
	userAPIClient userpb.UserV1Client
}

func NewHandler(userAPIClient userpb.UserV1Client) *Handler {
	return &Handler{userAPIClient: userAPIClient}
}
