package kadry

type mobileAppInfo struct {
	PersonID       string `json:"PersonID,omitempty"`
	FIOPerson      string `json:"FIOPerson,omitempty"`
	SNILS          string `json:"SNILS,omitempty"`
	OrgID          string `json:"OrgID,omitempty"`
	InnOrg         string `json:"InnOrg,omitempty"`
	NameOrg        string `json:"NameOrg,omitempty"`
	SubdivID       string `json:"SubdivID,omitempty"`
	NameSubdiv     string `json:"NameSubdiv,omitempty"`
	PositionID     string `json:"PositionID,omitempty"`
	NamePosition   string `json:"NamePosition,omitempty"`
	EmploymentType string `json:"EmploymentType,omitempty"`
	DateRecept     string `json:"DateRecept,omitempty"`
}

type mobileApp []mobileAppInfo
