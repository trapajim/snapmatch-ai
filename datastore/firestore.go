package datastore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"github.com/trapajim/snapmatch-ai/server/middleware"
	"github.com/trapajim/snapmatch-ai/services/ai"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"google.golang.org/api/iterator"
)

type FirestoreRepository[T snapmatchai.Entity] struct {
	client     *firestore.Client
	aiService  *ai.Service
	collection string
}

func NewFirestoreRepository[T snapmatchai.Entity](client *firestore.Client, aiService *ai.Service, collection string) *FirestoreRepository[T] {
	return &FirestoreRepository[T]{
		client:     client,
		aiService:  aiService,
		collection: collection,
	}
}

func (r *FirestoreRepository[T]) getSessionCollection(ctx context.Context) (*firestore.CollectionRef, error) {
	sess := middleware.GetSession(ctx)
	if sess == nil {
		return nil, errors.New("session cannot be empty")
	}
	return r.client.Collection(r.collection).Doc(sess.SessionID()).Collection("documents"), nil
}

func (r *FirestoreRepository[T]) Create(ctx context.Context, entity T) error {
	collection, err := r.getSessionCollection(ctx)
	if err != nil {
		return err
	}

	doc := collection.NewDoc()
	entity.SetID(doc.ID)

	_, err = doc.Set(ctx, entity)
	if err != nil {
		return nil
	}
	if indexable, ok := any(entity).(snapmatchai.IndexableEntity); ok {
		sess := middleware.GetSession(ctx)
		indexable.SetOwner(sess.SessionID())
		_, err = r.client.Collection("search").NewDoc().Set(ctx, indexable)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *FirestoreRepository[T]) Read(ctx context.Context, id string) (T, error) {
	var entity T
	collection, err := r.getSessionCollection(ctx)
	if err != nil {
		return entity, err
	}

	doc, err := collection.Doc(id).Get(ctx)
	if err != nil {
		return entity, err
	}
	if err := doc.DataTo(&entity); err != nil {
		return entity, err
	}
	entity.SetID(doc.Ref.ID)
	return entity, nil
}

func (r *FirestoreRepository[T]) Update(ctx context.Context, entity T) error {
	collection, err := r.getSessionCollection(ctx)
	if err != nil {
		return err
	}

	_, err = collection.Doc(entity.GetID()).Set(ctx, entity)
	return err
}

func (r *FirestoreRepository[T]) Delete(ctx context.Context, id string) error {
	collection, err := r.getSessionCollection(ctx)
	if err != nil {
		return err
	}

	_, err = collection.Doc(id).Delete(ctx)
	return err
}

func (r *FirestoreRepository[T]) List(ctx context.Context, filters map[string]any) ([]T, error) {
	var results []T
	collection, err := r.getSessionCollection(ctx)
	if err != nil {
		return nil, err
	}

	query := collection.Query
	for field, value := range filters {
		query = query.Where(field, "==", value)
	}

	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		var entity T
		if err := doc.DataTo(&entity); err != nil {
			return nil, err
		}
		entity.SetID(doc.Ref.ID)
		results = append(results, entity)
	}

	return results, nil
}

func (r *FirestoreRepository[T]) Search(ctx context.Context, query string) ([]T, error) {
	var results []T
	sess := middleware.GetSession(ctx)
	if sess == nil {
		return nil, errors.New("session cannot be empty")
	}
	emb, err := r.aiService.GenerateEmbeddings(ctx, query)
	if err != nil {
		return nil, err
	}

	vectorQuery := r.client.Collection("search").Where("Owner", "==", sess.SessionID()).FindNearest("VectorData", emb, 10, firestore.DistanceMeasureEuclidean, nil)
	iter := vectorQuery.Documents(ctx)
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		var entity T
		if err := doc.DataTo(&entity); err != nil {
			return nil, err
		}
		entity.SetID(doc.Ref.ID)
		results = append(results, entity)
	}

	return results, nil
}
