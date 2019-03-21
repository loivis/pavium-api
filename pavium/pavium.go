package pavium

type Chapter struct {
	Name string `json:"name,omitempty"`
	Link string `json:"link,omitempty"`
}

type SiteName string

const (
	Piaotian SiteName = "飘天文学网"
)

type Site interface {
	Chapters(link string) []Chapter
}
