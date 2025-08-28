package portalsv2

type Complex struct {
	ID              int    `json:"id"`
	Sort            int    `json:"sort"`
	Group           int    `json:"group"`
	FirstName       string `json:"firstName"`
	MiddleName      string `json:"middleName"`
	LastName        string `json:"lastName"`
	IconID          string `json:"photo"`
	HeadDescription string `json:"headDescription"`
	PortalIDs       []int  `json:"oivIds"`
}
