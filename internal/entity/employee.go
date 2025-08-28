package entity

type EmployeeGetKey int

const (
	EmployeeGetKeyInvalid EmployeeGetKey = iota
	EmployeeGetKeyCloudID
	EmployeeGetKeyEmail
)

func (ek EmployeeGetKey) String() string {
	switch ek {
	case EmployeeGetKeyCloudID:
		return "cloud-id"
	case EmployeeGetKeyEmail:
		return "email"
	default:
		return ""
	}
}

type EmployeeGetParams struct {
	Key     EmployeeGetKey
	CloudID string
	Email   string
}

func (ep EmployeeGetParams) ParamByKey() string {
	switch ep.Key {
	case EmployeeGetKeyCloudID:
		return ep.CloudID
	case EmployeeGetKeyEmail:
		return ep.Email
	default:
		return ""
	}
}

type EmployeeInfo struct {
	CloudID string `json:"cloud_id,omitempty"`
	Inn     string `json:"inn,omitempty"`
	OrgID   string `json:"org_id,omitempty"`
	FIO     string `json:"fio,omitempty"`
	SNILS   string `json:"snils,omitempty"`
}
