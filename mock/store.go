package mock

import (
	"context"

	"github.com/loivis/pavium-api/pavium"
)

type Store struct {
	PutFunc    func(ctx context.Context, fav *pavium.Favorite) error
	DeleteFunc func(ctx context.Context, fav *pavium.Favorite) error
	ListFunc   func(ctx context.Context) []pavium.Favorite
}

func (s *Store) Put(ctx context.Context, fav *pavium.Favorite) error {
	return s.PutFunc(ctx, fav)
}

func (s *Store) Delete(ctx context.Context, fav *pavium.Favorite) error {
	return s.DeleteFunc(ctx, fav)
}

func (s *Store) List(ctx context.Context) []pavium.Favorite {
	return s.ListFunc(ctx)
}
