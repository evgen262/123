package users

type ShortUser struct {
	// TODO: добавить id пользователя когда будет сервис пользователей
	// ID         string     `json:"id"`
	LastName   string     `json:"lastName"`
	FirstName  string     `json:"firstName"`
	MiddleName string     `json:"middleName"`
	ImageID    string     `json:"photoUrl"`
	Gender     string     `json:"gender"`
	PortalData PortalData `json:"portalData"`
}

type PortalData struct {
	PersonID   string `json:"personId"`
	EmployeeID string `json:"employeeId"`
}
