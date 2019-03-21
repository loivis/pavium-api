package service

import (
	"testing"

	"github.com/loivis/prunusavium-go/pavium"

	"github.com/loivis/prunusavium-go/mock"
)

func TestService_New(t *testing.T) {
	l1 := &mock.Site{}
	l2 := &mock.Site{}
	lefts := map[pavium.SiteName]pavium.Left{
		"l1": l1,
		"l2": l2,
	}

	r1 := &mock.Site{}
	r2 := &mock.Site{}
	rights := map[pavium.SiteName]pavium.Right{
		"r1": r1,
		"r2": r2,
	}

	svc := New(
		WithLefts(lefts),
		WithRights(rights),
	)

	if got, want := len(svc.lefts), len(lefts); got != want {
		t.Errorf("len(svc.lefts) = %d, want %d", got, want)
	}

	if got, want := len(svc.rights), len(rights); got != want {
		t.Errorf("len(svc.rights) = %d, want %d", got, want)
	}

	if got, want := len(svc.sites), len(lefts)+len(rights); got != want {
		t.Errorf("len(svc.sites) = %d, want %d", got, want)
	}
}

func TestService_WithLefts(t *testing.T) {
	l1 := &mock.Site{}
	l2 := &mock.Site{}
	lefts := map[pavium.SiteName]pavium.Left{
		"l1": l1,
		"l2": l2,
	}
	svc := New(
		WithLefts(lefts),
	)

	if got, want := len(svc.lefts), len(lefts); got != want {
		t.Errorf("len(svc.lefts) = %d, want %d", got, want)
	}

	if got, want := len(svc.sites), len(lefts); got != want {
		t.Errorf("len(svc.sites) = %d, want %d", got, want)
	}
}

func TestService_WithRights(t *testing.T) {
	r1 := &mock.Site{}
	r2 := &mock.Site{}
	rights := map[pavium.SiteName]pavium.Right{
		"r1": r1,
		"r2": r2,
	}
	svc := New(
		WithRights(rights),
	)

	if got, want := len(svc.rights), len(rights); got != want {
		t.Errorf("len(svc.rights) = %d, want %d", got, want)
	}

	if got, want := len(svc.sites), len(rights); got != want {
		t.Errorf("len(svc.sites) = %d, want %d", got, want)
	}
}

func TestService_WithFavStore(t *testing.T) {
	store := &mock.Store{}
	svc := New(
		WithFavStore(store),
	)

	if got, want := svc.favStore, store; got != want {
		t.Errorf("service.store = %v, want %v", got, want)
	}
}
