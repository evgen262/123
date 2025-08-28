package portalsv2

type Portal struct {
	ID             int      `json:"id"`
	IconID         string   `json:"icon"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Count          Count    `json:"count"`
	SctructureType []string `json:"sctructureType"`
	Head           *Head    `json:"head"`
}

type Count struct {
	Employees int `json:"employees"`
	Podved    int `json:"podved"`
}

type Head struct {
	FirstName   string `json:"firstName"`
	MiddleName  string `json:"middleName"`
	LastName    string `json:"lastName"`
	Description string `json:"description"`
	ImageID     string `json:"image"`
}

type PortalsFilterRequest struct {
	PortalIDs []int `json:"ids"`
}
