package datastore

import (
	"cloud.google.com/go/bigquery"
	"context"
	"errors"
	"fmt"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"reflect"
)

type BigQuery struct {
	client *bigquery.Client
}

func NewBigQuery(client *bigquery.Client) *BigQuery {
	return &BigQuery{
		client: client,
	}
}

func (q BigQuery) Query(ctx context.Context, queryString string, parameters map[string]any, target any) error {
	query := q.client.Query(queryString)
	if err := q.addQueryParameters(query, parameters); err != nil {
		return err
	}

	if _, isVoid := target.(snapmatchai.Void); isVoid {
		return q.executeVoidQuery(ctx, query)
	}

	return q.readQueryResults(ctx, query, target)
}

// TableExists checks if a table exists, if the table does not exist an error is returned
func (q BigQuery) TableExists(ctx context.Context, dataset, tableName string) error {
	_, err := q.client.Dataset(dataset).Table(tableName).Metadata(ctx)
	return handleApiError(err)
}

func (q BigQuery) Schema(ctx context.Context, dataset, tableName string) ([]snapmatchai.DBSchema, error) {
	meta, err := q.client.Dataset(dataset).Table(tableName).Metadata(ctx)
	if err != nil {
		return nil, handleApiError(err)
	}
	schema := make([]snapmatchai.DBSchema, len(meta.Schema))
	for i, field := range meta.Schema {
		schema[i] = snapmatchai.DBSchema{
			Name: field.Name,
			Type: string(field.Type),
		}
	}
	return schema, nil
}

func (q BigQuery) addQueryParameters(query *bigquery.Query, parameters map[string]any) error {
	for param, value := range parameters {
		query.QueryConfig.Parameters = append(query.QueryConfig.Parameters, bigquery.QueryParameter{
			Name:  param,
			Value: value,
		})
	}
	return nil
}

// executeVoidQuery executes the query and waits for the result without populating the target
func (q BigQuery) executeVoidQuery(ctx context.Context, query *bigquery.Query) error {
	s, err := query.Run(ctx)
	if err != nil {
		return handleApiError(err)
	}
	res, err := s.Wait(ctx)
	if err != nil {
		return handleApiError(err)
	}
	if res.Err() != nil {
		return handleApiError(res.Err())
	}
	return nil
}

// readQueryResults reads the query results and populates the target struct
func (q BigQuery) readQueryResults(ctx context.Context, query *bigquery.Query, target any) error {
	it, err := query.Read(ctx)
	if err != nil {
		return handleApiError(err)
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("target must be a pointer to a slice")
	}

	sliceValue := v.Elem()
	elemType := sliceValue.Type().Elem()

	for {
		elemPtr := reflect.New(elemType)
		if err := it.Next(elemPtr.Interface()); errors.Is(err, iterator.Done) {
			break
		} else if err != nil {
			return handleApiError(err)
		}
		sliceValue.Set(reflect.Append(sliceValue, elemPtr.Elem()))
	}
	return nil
}

func handleApiError(err error) error {
	if err == nil {
		return nil
	}
	var e *googleapi.Error
	if ok := errors.As(err, &e); ok {
		return snapmatchai.NewError(err, fmt.Sprintf("Google API error: %s, Code: %d", e.Message, e.Code), e.Code)
	}
	return snapmatchai.NewError(err, "error occurred, during querying", 500)
}
