package store

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/loivis/prunusavium-go/pavium"
	"google.golang.org/api/iterator"
)

type Store struct {
	project    string
	collection string
}

func New(project, collection string) *Store {
	return &Store{
		project:    project,
		collection: collection,
	}
}

func (s *Store) Put(ctx context.Context, fav *pavium.Favorite) error {
	id := fav.ID()

	client, err := firestore.NewClient(ctx, s.project)
	if err != nil {
		log.Printf("failed to create firestore client: %v", err)
		return err
	}

	res, err := client.Collection(s.collection).Doc(id).Set(ctx, fav)
	if err != nil {
		log.Printf("failed to add favorite(%+v): %v", fav, err)
		return err
	}

	log.Printf("favorite %s(%+v) added at %v", id, fav, res.UpdateTime)

	return nil
}

func (s *Store) Delete(ctx context.Context, fav *pavium.Favorite) error {
	id := fav.ID()

	client, err := firestore.NewClient(ctx, s.project)
	if err != nil {
		log.Printf("failed to create firestore client: %v", err)
		return err
	}

	res, err := client.Collection(s.collection).Doc(id).Delete(ctx)
	if err != nil {
		log.Printf("failed to delete favorite(%+v): %v", fav, err)
		return err
	}

	// TODO: res.Update seems to be null time.Time
	log.Printf("favorite %s(%+v) deleted at %v", id, fav, res.UpdateTime)

	return nil
}

func (s *Store) List(ctx context.Context) []pavium.Favorite {
	var favs []pavium.Favorite

	client, err := firestore.NewClient(ctx, s.project)
	if err != nil {
		log.Printf("failed to create firestore client: %v", err)
		return favs
	}

	iter := client.Collection(s.collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("failed to iterate: %v", err)
		}

		var fav pavium.Favorite
		if err := doc.DataTo(&fav); err != nil {
			log.Printf("failed to get favorite: %v", err)
			favs = []pavium.Favorite{}
			break
		}

		favs = append(favs, fav)
	}

	return favs
}
