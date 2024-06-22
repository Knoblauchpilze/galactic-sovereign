package communication

type AclResponseDto struct {
	Resource    string   `json:"resource"`
	Permissions []string `json:"permissions"`
}
