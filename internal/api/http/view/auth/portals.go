package auth

/*
type User struct {
	// Фамилия пользователя
	LastName string `json:"lastName"`
	// Имя пользователя
	FirstName string `json:"firstName"`
	// Отчество пользователя
	MiddleName string `json:"middleName"`
}
*/

type Portal struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Image    string `json:"image"`
	IsActive bool   `json:"isActive"`
}
