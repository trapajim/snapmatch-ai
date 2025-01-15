package asset

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/trapajim/snapmatch-ai/pipeline"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"io"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Similarity int

const (
	High Similarity = iota
	Medium
	Low
)

type Service struct {
	appContext snapmatchai.Context
}

// NewService creates a new instance of the asset service
// the asset service is responsible for uploading files to storage
func NewService(appContext snapmatchai.Context) *Service {
	return &Service{
		appContext: appContext,
	}
}

// Upload uploads a file to storage
func (s *Service) Upload(ctx context.Context, file io.Reader, fileName string) error {
	err := s.appContext.Storage.Upload(ctx, file, fileName)
	if err != nil {
		var e *snapmatchai.Error
		if ok := errors.As(err, &e); ok {
			s.appContext.Logger.ErrorContext(ctx, "Google API error occurred, during file upload",
				slog.Int("status_code", e.Code),
				slog.String("error", e.Error()),
				slog.String("message", e.Message),
			)
			return NewUploadError(e, e.Error(), e.Code)
		}
		return snapmatchai.NewError(err, fmt.Errorf("could not upload file %s with error %w", fileName, err).Error(), 500)
	}
	return nil
}

type BatchUploadRequest struct {
	File io.Reader
	Name string
}

// BatchUpload uploads multiple files concurrently to storage
func (s *Service) BatchUpload(ctx context.Context, files chan BatchUploadRequest) error {
	var uploadErrs []error
	var wg sync.WaitGroup
	var mu sync.Mutex
	go func() {
		for file := range files {
			wg.Add(1)
			go func(f BatchUploadRequest) {
				defer wg.Done()
				optimizedFile, err := optimizeImage(f.File, f.Name)
				if err != nil {
					mu.Lock()
					uploadErrs = append(uploadErrs, err)
					mu.Unlock()
					log.Println(err)
					return
				}
				err = s.Upload(ctx, optimizedFile, f.Name)
				if err != nil {
					log.Println(err)
					mu.Lock()
					uploadErrs = append(uploadErrs, err)
					mu.Unlock()
				}
			}(file)
		}
	}()
	wg.Wait()
	err := s.index(ctx, s.appContext)
	if len(uploadErrs) > 0 {
		_ = s.index(ctx, s.appContext)
		return fmt.Errorf("batch upload failed: %v", uploadErrs)
	}
	err = s.index(ctx, s.appContext)
	if errors.Is(err, tableCreatedError) {
		return nil
	}
	return err
}

// Search searches for files in storage
func (s *Service) Search(ctx context.Context, query string, similarity Similarity, page snapmatchai.Pagination) ([]snapmatchai.FileRecord, snapmatchai.Pagination, error) {
	query, prms, err := buildQuery(s.appContext, query, similarity, page)
	if err != nil {
		return nil, snapmatchai.Pagination{}, err
	}
	var rows []snapmatchai.FileRecord
	err = s.appContext.DB.Query(ctx, query, prms, &rows)
	if err != nil {
		return nil, page, err
	}
	if len(rows) == 0 {
		return rows, snapmatchai.Pagination{
			NextToken: "",
			Per:       page.Per,
		}, nil
	}
	nextToken := ""
	if len(rows) == 50 {
		nextToken = rows[len(rows)-1].Updated.Format(time.RFC3339)
	}

	s.SignURLs(ctx, rows)
	pagination := snapmatchai.NewPagination(nextToken, 50)

	return rows, *pagination, nil
}

func (s *Service) List(ctx context.Context, page snapmatchai.Pagination) ([]snapmatchai.FileRecord, snapmatchai.Pagination, error) {
	var where string
	if page.NextToken != "" {
		where += "WHERE updated > @updated"
	}
	query := fmt.Sprintf(`
SELECT *, 
  COALESCE(
  (SELECT value 
   FROM UNNEST(metadata) 
   WHERE name = 'category'), 
  ''
) AS category
FROM EXTERNAL_OBJECT_TRANSFORM(TABLE %s.%s, ['SIGNED_URL'])
%s
ORDER BY category, updated DESC
LIMIT 50
`, s.appContext.Config.DatasetID, s.appContext.Config.TableID, where)
	params := make(map[string]any)
	if page.NextToken != "" {
		t, err := page.DecodeNextToken()
		if err != nil {
			return nil, page, errors.New("invalid next token")
		}
		params["updated"] = t
	}
	var rows []snapmatchai.FileRecord
	err := s.appContext.DB.Query(ctx, query, params, &rows)
	if err != nil {
		var errAs *snapmatchai.Error
		if errors.As(err, &errAs) {
			log.Println(errAs.Message)
			log.Println(errAs.Unwrap().Error())

		}
		return nil, page, err
	}
	if len(rows) == 0 {
		return rows, snapmatchai.Pagination{
			NextToken: "",
			Per:       page.Per,
		}, nil
	}
	nextToken := ""
	if len(rows) == 50 {
		nextToken = rows[len(rows)-1].Updated.Format(time.RFC3339)
	}

	pagination := snapmatchai.NewPagination(nextToken, 50)

	return rows, *pagination, nil
}

func (s *Service) SignURLs(ctx context.Context, records []snapmatchai.FileRecord) {
	for i := range records {
		signedUrl, err := s.appContext.Storage.SignUrl(ctx, records[i].ObjName, time.Hour)
		if err != nil {
			s.appContext.Logger.ErrorContext(ctx, "could not sign url", slog.String("error", err.Error()))
			continue
		}
		records[i].SignedURL = signedUrl
	}
}
func buildQuery(appContext snapmatchai.Context, searchTerm string, similarity Similarity, page snapmatchai.Pagination) (string, map[string]any, error) {
	table := fmt.Sprintf("%s_embeddings", appContext.Config.TableID)
	parameters := make(map[string]any)
	pageQuery := ""
	if page.NextToken != "" {
		pageQuery = fmt.Sprintf(" AND (distance > @last_distance OR (distance = @last_distance AND updated > @last_updated)) ")
	}
	distance := getDistance(similarity)
	query := fmt.Sprintf(`
WITH search_results AS ( 
  SELECT base.*, distance 
  FROM VECTOR_SEARCH(
    (
      SELECT *
      FROM %s.%s
    ), 'ml_generate_embedding_result',
    (
      SELECT ml_generate_embedding_result, content AS query
      FROM ML.GENERATE_EMBEDDING(
        MODEL %s,
        (SELECT '%s' AS content))
    ),
    top_k => -1, options => '{"fraction_lists_to_search": 0.01}')),
highest_distance AS (
  SELECT MIN(distance) AS best_distance
  FROM search_results
)
SELECT *
FROM search_results, highest_distance
WHERE distance <= best_distance * %f
%s
ORDER BY distance ASC, updated ASC
LIMIT 50;
`, appContext.Config.DatasetID, table, appContext.Config.BQMultiModalModel, searchTerm, distance, pageQuery)
	if page.NextToken != "" {
		t, err := page.DecodeNextToken()
		if err != nil {
			return "", nil, NewPaginationError(err, "invalid next token")
		}
		updated, dist, err := parsePaginationToken(t)
		if err != nil {
			return "", nil, NewPaginationError(err, "invalid next token")
		}
		parameters["last_updated"] = updated
		parameters["last_distance"] = dist
	}
	return query, parameters, nil
}

func getDistance(similarity Similarity) float64 {
	switch similarity {
	case High:
		return 1.02
	case Medium:
		return 1.03
	case Low:
		return 1.1
	default:
		return 1.02
	}
}
func parsePaginationToken(token string) (string, float64, error) {
	splitString := strings.Split(token, "##")
	if len(splitString) != 2 {
		return "", 0, errors.New("invalid next token")
	}
	f, err := strconv.ParseFloat(splitString[1], 64)
	if err != nil {
		return "", 0, errors.New("invalid next token")
	}
	return splitString[0], f, nil
}

type indexPipelineState struct {
	appContext snapmatchai.Context
}

func (s *Service) index(ctx context.Context, appCtx snapmatchai.Context) error {
	state := indexPipelineState{appContext: appCtx}
	p := pipeline.New(state, pipeline.WithLogger(appCtx.Logger)).Then(
		pipeline.NamedStep[indexPipelineState]{StepName: "Exists Asset Table", ExecuteFn: func(state indexPipelineState) error {
			return TableExists(ctx, state.appContext, state.appContext.Config.TableID)
		}}).OnError(func(state indexPipelineState, err error) error {
		var snapMatchErr *snapmatchai.Error
		if errors.As(err, &snapMatchErr) {
			log.Println(snapMatchErr.Code)
			if snapMatchErr.Code == http.StatusNotFound {
				return CreateAssetTable(ctx, state.appContext)
			}
		}
		return err
	}).Then(pipeline.NamedStep[indexPipelineState]{StepName: "Exists Embeddings Table", ExecuteFn: func(state indexPipelineState) error {
		return TableExists(ctx, state.appContext, embeddingsTableName(state.appContext.Config.TableID))
	}}).OnError(func(state indexPipelineState, err error) error {
		var snapMatchErr *snapmatchai.Error
		if errors.As(err, &snapMatchErr) {
			if snapMatchErr.Code == http.StatusNotFound {
				err2 := CreateEmbeddingTable(ctx, state.appContext)
				if err2 == nil {
					return tableCreatedError
				}
				return err2
			}
			return err
		}
		return err
	}).Then(pipeline.NamedStep[indexPipelineState]{StepName: "Update index", ExecuteFn: func(state indexPipelineState) error {
		return UpdateIndex(ctx, state.appContext, time.Now().Add(-time.Minute))
	}}).Execute()
	return p
}

func optimizeImage(file io.Reader, fileName string) (io.Reader, error) {
	img, err := imaging.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}
	resizedImg := imaging.Resize(img, 1024, 0, imaging.Lanczos)

	var buf bytes.Buffer
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg":
		err = imaging.Encode(&buf, resizedImg, imaging.JPEG, imaging.JPEGQuality(80))
	case "png":
		err = imaging.Encode(&buf, resizedImg, imaging.PNG)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	return &buf, nil
}
