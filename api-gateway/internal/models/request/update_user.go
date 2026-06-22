package request

type UpdateUser struct {
	Grade    string `json:"grade"`
	Password string `json:"password"`
}
