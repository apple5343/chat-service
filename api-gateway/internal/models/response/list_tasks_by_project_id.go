package response

type ListTasksByProjectID struct {
	Tasks []GetTask `json:"tasks"`
}
