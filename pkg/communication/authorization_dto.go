package communication

type AuthorizationResponseDto struct {
	Acls   []AclResponseDto   `json:"acls"`
	Limits []LimitDtoResponse `json:"limits"`
}
