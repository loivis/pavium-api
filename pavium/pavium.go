package pavium

type Chapter struct {
	Name string `json:"name,omitempty"`
	Link string `json:"link,omitempty"`
}

type SiteName string

type Site interface {
	Chapters(link string) []Chapter
}
