package communication

type AuthorizationResponseDto struct {
	Acls   []AclResponseDto   `json:"acls"`
	Limits []LimitResponseDto `json:"limits"`
}
