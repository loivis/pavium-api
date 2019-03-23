package piaotian

import "testing"

func TestSite_ParseChapterLink(t *testing.T) {
	s := &Site{chapterURL: "https://foo.bar/html/"}

	for _, tc := range []struct {
		desc string
		in   string
		out  string
	}{
		{
			desc: "bookinfo",
			in:   "https://foo.bar/bookinfo/123/456.html",
			out:  "https://foo.bar/html/123/456/",
		},
		{
			desc: "html",
			in:   "https://foo.bar/html/123/456/",
			out:  "https://foo.bar/html/123/456/",
		},
		{
			desc: "empty",
			in:   "",
			out:  "",
		},
		{
			desc: "invalid1",
			in:   "https://foo.bar/html/123/",
			out:  "",
		},
		{
			desc: "invalid2",
			in:   "https://foo.bar/html/123",
			out:  "",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			if got, want := s.parseChapterLink(tc.in), tc.out; got != want {
				t.Fatalf("parsed link = %q, want %q", got, want)
			}
		})
	}
}
