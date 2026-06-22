package request

type CreateProject struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Users       []string `json:"users"`
}
