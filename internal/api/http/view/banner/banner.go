package banners

type Banner struct {
	ImageUrl string `json:"imageUrl"`
	Url      string `json:"url"`
	Order    int    `json:"order"`
}

type BannersList struct {
	Items    []*Banner `json:"items"`
	Interval int       `json:"interval"`
}
