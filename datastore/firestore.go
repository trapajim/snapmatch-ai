package datastore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
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

func (r *FirestoreRepository[T]) Create(ctx context.Context, entity T) error {
	doc := r.client.Collection(r.collection).NewDoc()
	entity.SetID(doc.ID)
	_, err := doc.Set(ctx, entity)
	return err
}

func (r *FirestoreRepository[T]) Read(ctx context.Context, id string) (T, error) {
	var entity T
	doc, err := r.client.Collection(r.collection).Doc(id).Get(ctx)
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
	_, err := r.client.Collection(r.collection).Doc(entity.GetID()).Set(ctx, entity)
	return err
}

func (r *FirestoreRepository[T]) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection(r.collection).Doc(id).Delete(ctx)
	return err
}

func (r *FirestoreRepository[T]) List(ctx context.Context, filters map[string]any) ([]T, error) {
	var results []T
	query := r.client.Collection(r.collection).Query
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
	collection := r.client.Collection(r.collection)
	emb, err := r.aiService.GenerateEmbeddings(ctx, query)
	if err != nil {
		return nil, err
	}
	vectorQuery := collection.FindNearest("VectorData", emb, 10, firestore.DistanceMeasureEuclidean, nil)
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
