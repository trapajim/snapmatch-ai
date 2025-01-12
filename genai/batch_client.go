package genai

import (
	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"context"
	"fmt"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"google.golang.org/protobuf/types/known/structpb"
)

type BatchClient struct {
	client    *aiplatform.JobClient
	location  string
	projectID string
}

func NewBatchClient(client *aiplatform.JobClient, location, projectID string) *BatchClient {
	return &BatchClient{
		location:  location,
		projectID: projectID,
		client:    client,
	}
}

func (b *BatchClient) CreateBatchPredictionJob(ctx context.Context, config snapmatchai.BatchPredictionRequest) (snapmatchai.BatchPredictionJobConfig, error) {
	modelParameters, err := structpb.NewValue(config.ModelParameters)
	if err != nil {
		return snapmatchai.BatchPredictionJobConfig{}, snapmatchai.NewError(err, "failed to create model", 400)
	}
	req := &aiplatformpb.CreateBatchPredictionJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", b.projectID, b.location),
		BatchPredictionJob: &aiplatformpb.BatchPredictionJob{
			DisplayName:     config.JobName,
			Model:           fmt.Sprintf("publishers/google/models/%s", config.ModelName),
			ModelParameters: modelParameters,
			InputConfig: &aiplatformpb.BatchPredictionJob_InputConfig{
				Source: &aiplatformpb.BatchPredictionJob_InputConfig_GcsSource{
					GcsSource: &aiplatformpb.GcsSource{
						Uris: []string{config.InputPath},
					},
				},
				InstancesFormat: "jsonl",
			},
			OutputConfig: &aiplatformpb.BatchPredictionJob_OutputConfig{
				Destination: &aiplatformpb.BatchPredictionJob_OutputConfig_GcsDestination{
					GcsDestination: &aiplatformpb.GcsDestination{
						OutputUriPrefix: config.OutputPath,
					},
				},
				PredictionsFormat: "jsonl",
			},
		},
	}

	r, err := b.client.CreateBatchPredictionJob(ctx, req)
	if err != nil {
		return snapmatchai.BatchPredictionJobConfig{}, snapmatchai.NewError(err, "failed to create batch prediction job", 400)
	}
	return snapmatchai.BatchPredictionJobConfig{
		Name:   r.GetName(),
		Status: r.GetState().String(),
	}, nil
}
