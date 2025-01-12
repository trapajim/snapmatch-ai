package snapmatchai

type BatchPredictionRequest struct {
	JobName         string
	ModelName       string
	ModelParameters map[string]any
	InputPath       string
	OutputPath      string
}

type BatchPredictionJobConfig struct {
	Name   string
	Status string
}

func NewBatchPredictionRequest(jobName, modelName, inputPath, outputPath string, modelParameters map[string]any) BatchPredictionRequest {
	return BatchPredictionRequest{
		JobName:         jobName,
		ModelName:       modelName,
		ModelParameters: modelParameters,
		InputPath:       inputPath,
		OutputPath:      outputPath,
	}
}
