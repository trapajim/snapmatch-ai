package asset

import (
	"context"
	"fmt"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"strings"
	"time"
)

func embeddingsTableName(tableID string) string {
	return fmt.Sprintf("%s_embeddings", tableID)
}

func TableExists(ctx context.Context, appContext snapmatchai.Context, table string) error {
	return appContext.DB.TableExists(ctx, appContext.Config.DatasetID, table)
}

func CreateAssetTable(ctx context.Context, appContext snapmatchai.Context) error {
	connection := fmt.Sprintf("`projects/%s/locations/%s/connections/%s`", appContext.Config.ProjectID, appContext.Config.Location, appContext.Config.BQVertexConn)
	tableName := fmt.Sprintf("`%s.%s.%s`", appContext.Config.ProjectID, appContext.Config.DatasetID, appContext.Config.TableID)
	q := fmt.Sprintf(`
	CREATE EXTERNAL TABLE %s 
	WITH CONNECTION %s
	OPTIONS(
		object_metadata = 'SIMPLE',
		uris = ['gs://%s/*']
);
	`, tableName, connection, appContext.Config.StorageBucket)
	return appContext.DB.Query(ctx, q, nil, snapmatchai.Void{})
}

func CreateEmbeddingTable(ctx context.Context, appContext snapmatchai.Context) error {
	sourceTableName := fmt.Sprintf("`%s.%s.%s`", appContext.Config.ProjectID, appContext.Config.DatasetID, appContext.Config.TableID)
	tableName := fmt.Sprintf("`%s.%s.%s`", appContext.Config.ProjectID, appContext.Config.DatasetID, embeddingsTableName(appContext.Config.TableID))
	q := fmt.Sprintf(`CREATE OR REPLACE TABLE %s AS (
  SELECT *,
    REGEXP_EXTRACT(uri, r'[^/]+$') AS obj_name,
    REGEXP_REPLACE(REGEXP_EXTRACT(uri, r'[^/]+$'), r'\.png$', '') AS file_id,
  FROM ML.GENERATE_EMBEDDING(
    MODEL %s,
    TABLE %s 
  )
);`, tableName, appContext.Config.BQMultiModalModel, sourceTableName)
	return appContext.DB.Query(ctx, q, nil, snapmatchai.Void{})
}

func UpdateIndex(ctx context.Context, appContext snapmatchai.Context, updated time.Time) error {
	schema, err := appContext.DB.Schema(ctx, appContext.Config.DatasetID, embeddingsTableName(appContext.Config.TableID))
	if err != nil {
		panic(err)
	}

	targetAlias := "target"
	sourceAlias := "source"
	updateClause := generateUpdateClause(schema, targetAlias, sourceAlias)

	insertColumns, insertValues := generateInsertClause(schema, sourceAlias)

	embeddingsTable := fmt.Sprintf("`%s.%s`", appContext.Config.DatasetID, embeddingsTableName(appContext.Config.TableID))
	assetTable := fmt.Sprintf("`%s.%s`", appContext.Config.DatasetID, appContext.Config.TableID)
	sql := fmt.Sprintf(`
MERGE %s AS target
USING (
  SELECT *, 
      REGEXP_EXTRACT(uri, r'[^/]+$') AS obj_name,
      REGEXP_REPLACE(REGEXP_EXTRACT(uri, r'[^/]+$'), r'\.png$', '') AS file_id,
  FROM ML.GENERATE_EMBEDDING(
    MODEL %s,
    (
      SELECT *
      FROM %s 
      WHERE updated >= @updated 
    )
  )
) AS source
ON target.uri = source.uri
WHEN MATCHED THEN
  UPDATE SET %s
WHEN NOT MATCHED THEN
  INSERT (%s)
  VALUES (%s);
`, embeddingsTable, appContext.Config.BQMultiModalModel, assetTable, updateClause, insertColumns, insertValues)
	return appContext.DB.Query(ctx, sql, map[string]any{"updated": updated.UTC()}, snapmatchai.Void{})
}

func generateUpdateClause(schema []snapmatchai.DBSchema, targetAlias, sourceAlias string) string {
	var updateClauses []string
	for _, field := range schema {
		updateClauses = append(updateClauses, fmt.Sprintf("%s.%s = %s.%s", targetAlias, field.Name, sourceAlias, field.Name))
	}
	return strings.Join(updateClauses, ", ")
}

func generateInsertClause(schema []snapmatchai.DBSchema, sourceAlias string) (string, string) {
	var columns []string
	var values []string
	for _, field := range schema {
		columns = append(columns, field.Name)
		values = append(values, fmt.Sprintf("%s.%s", sourceAlias, field.Name))
	}
	return strings.Join(columns, ", "), strings.Join(values, ", ")
}
