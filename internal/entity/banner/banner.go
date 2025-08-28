package banner

type Banner struct {
	ImageUrl string
	Url      string
	Order    int
}

type BannersList struct {
	Items    []*Banner
	Interval int
}
