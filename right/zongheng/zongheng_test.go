package zongheng

import "testing"

func TestSite_ParseChapterLink(t *testing.T) {
	s := &Site{chapterURL: "https://foo.bar/baz/%v.html"}

	for _, tc := range []struct {
		desc string
		in   string
		out  string
	}{
		{
			desc: "Success",
			in:   "http://book.zongheng.com/book/12345.html?fr=pc_alading",
			out:  "https://foo.bar/baz/12345.html",
		},
		{
			desc: "WrongPrefix",
			in:   "http://foo.bar/book/12345.html",
			out:  "",
		},
		{
			desc: "NoDotSplit",
			in:   "http://book.zongheng.com/book/12345",
			out:  "",
		},
		{
			desc: "NoNumericID",
			in:   "http://book.zongheng.com/book/abc.html",
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
