package server

import (
	"api-gateway/internal/server/auth"
	"api-gateway/internal/server/chat"
	"api-gateway/internal/server/helper"
	"api-gateway/internal/server/project"
	"api-gateway/internal/server/task"
	"api-gateway/internal/server/user"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewServer инициализирует новый сервер, применяет middleware и устанавливает пути.
func NewServer(
	userHandler *user.Handler,
	authHandler *auth.Handler,
	projectHandler *project.Handler,
	taskHandler *task.Handler,
	chatHandler *chat.Handler,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	addRoutes(router, userHandler, authHandler, projectHandler, taskHandler, chatHandler)

	var handler http.Handler = router

	return handler
}

func addRoutes(
	router *chi.Mux,
	userHandler *user.Handler,
	authHandler *auth.Handler,
	projectHandler *project.Handler,
	taskHandler *task.Handler,
	chatHandler *chat.Handler,
) {
	router.Route("/api", func(r chi.Router) {
		r.Post("/register", helper.MakeHandler(userHandler.Register))
		r.Post("/login", helper.MakeHandler(authHandler.Login))
		r.Post("/get_refresh_token", helper.MakeHandler(authHandler.GetRefreshToken))
		r.Post("/get_access_token", helper.MakeHandler(authHandler.GetAccessToken))

		r.With(authHandler.Auth).Route("/users", func(r chi.Router) {
			r.Get("/{user_id}", helper.MakeHandler(userHandler.GetUser))
			r.Patch("/{user_id}", helper.MakeHandler(userHandler.UpdateUser))
			r.Delete("/{user_id}", helper.MakeHandler(userHandler.DeleteUser))
		})

		r.With(authHandler.Auth).Route("/projects", func(r chi.Router) {
			r.Get("/{project_id}", helper.MakeHandler(projectHandler.GetProject))
			r.Get("/", helper.MakeHandler(projectHandler.ListUserProjects))
			r.Post("/", helper.MakeHandler(projectHandler.CreateProject))
			r.Patch("/{project_id}", helper.MakeHandler(projectHandler.UpdateProject))
			r.Post("/{project_id}/add_user", helper.MakeHandler(projectHandler.AddUser))
			r.Delete("/{project_id}", helper.MakeHandler(projectHandler.DeleteProject))
			r.Post("/{project_id}/subscribe_on_notifications", helper.MakeHandler(projectHandler.SubscribeOnProjectNotifications))

			r.Route("/{project_id}/tasks", func(r chi.Router) {
				r.Post("/", helper.MakeHandler(taskHandler.CreateTask))
				r.Get("/{task_id}", helper.MakeHandler(taskHandler.GetTask))
				r.Get("/", helper.MakeHandler(taskHandler.ListTasksByProjectID))
				r.Patch("/{task_id}", helper.MakeHandler(taskHandler.UpdateTask))
				r.Delete("/{task_id}", helper.MakeHandler(taskHandler.DeleteTask))
			})
		})

		r.With(authHandler.Auth).Route("/chats", func(r chi.Router) {
			r.Get("/{project_id}/messages", helper.MakeHandler(chatHandler.GetMessages))
			r.Get("/ws", helper.MakeHandler(chatHandler.JoinRooms))
			r.Get("/", helper.MakeHandler(chatHandler.GetChats))
		})
	})
}
