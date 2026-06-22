package response

type ListUserProjects struct {
	Projects []GetProject `json:"projects"`
}
