package event

type CalendarEventRequest struct {
	SessionID string
	PortalURL string
	Year      int
	Month     int
}

type CalendarEvent struct {
	ID    string
	Title string
	Time  string
	Date  string
}

type CalendarEventsList struct {
	Items []*CalendarEvent
	Count int
}

type CalendarEventLinksRequest struct {
	SessionID string
	PortalURL string
	EventIDs  []string
}

type CalendarEventLink struct {
	ID          string
	RedirectURL string
}
