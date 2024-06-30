package communication

import "encoding/json"

type AuthorizationDtoResponse struct {
	Acls   []AclDtoResponse   `json:"acls"`
	Limits []LimitDtoResponse `json:"limits"`
}

func (d AuthorizationDtoResponse) MarshalJSON() ([]byte, error) {
	acls := d.Acls
	if d.Acls == nil {
		acls = make([]AclDtoResponse, 0)
	}

	userLimits := d.Limits
	if d.Limits == nil {
		userLimits = make([]LimitDtoResponse, 0)
	}

	out := struct {
		Acls   []AclDtoResponse   `json:"acls"`
		Limits []LimitDtoResponse `json:"limits"`
	}{
		Acls:   acls,
		Limits: userLimits,
	}

	return json.Marshal(out)
}
