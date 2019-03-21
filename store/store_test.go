package store

import "testing"

func TestNew(t *testing.T) {
	project := "foo"
	collection := "bar"

	store := New(project, collection)

	if got, want := store.project, project; got != want {
		t.Errorf("got store.project %q, want %q", got, want)
	}

	if got, want := store.collection, collection; got != want {
		t.Errorf("got store.collection %q, want %q", got, want)
	}
}
