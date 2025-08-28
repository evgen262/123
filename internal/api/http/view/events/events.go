package events

type CalendarEvent struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Time  string `json:"time"`
	Date  string `json:"date"`
}

type CalendarEventsList struct {
	Items []*CalendarEvent `json:"events"`
	Count int              `json:"totalCount"`
}

type CalendarEventLink struct {
	ID          string `json:"id"`
	RedirectURL string `json:"redirectUrl"`
}
