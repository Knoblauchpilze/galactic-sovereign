package communication

type AclResponseDto struct {
	Resource   string `json:"resource"`
	Permission string `json:"permission"`
}
