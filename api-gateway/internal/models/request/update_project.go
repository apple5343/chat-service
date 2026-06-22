package request

type UpdateProject struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AdminID     string `json:"admin_id"`
}
