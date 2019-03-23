package qidian

import "testing"

func TestSite_ParseChapterLink(t *testing.T) {
	s := &Site{home: "https://foo.bar"}

	for _, tc := range []struct {
		desc string
		in   string
		out  string
	}{
		{
			desc: "Success",
			in:   "https://foo.bar/info/12345",
			out:  "https://foo.bar/info/12345#Catalog",
		},
		{
			desc: "Success",
			in:   "https://foo.bar/baz",
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
