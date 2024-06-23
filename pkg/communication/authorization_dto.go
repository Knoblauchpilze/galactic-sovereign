package communication

type AuthorizationResponseDto struct {
	Acls   []AclDtoResponse   `json:"acls"`
	Limits []LimitDtoResponse `json:"limits"`
}
