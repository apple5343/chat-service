package request

type GetAccessToken struct {
	RefreshToken string `json:"refresh_token"`
}
