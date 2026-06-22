package task

import taskpb "api-gateway/pkg/api/task_v1"

type Handler struct {
	taskAPIClient taskpb.TaskServiceClient
}

func NewHandler(taskAPIClient taskpb.TaskServiceClient) *Handler {
	return &Handler{
		taskAPIClient: taskAPIClient,
	}
}
