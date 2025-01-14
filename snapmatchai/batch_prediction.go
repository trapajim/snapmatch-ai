package snapmatchai

type BatchPrediction struct {
	ID              string
	JobType         string
	Status          string
	InternalName    string
	JobName         string
	ModelName       string
	ModelParameters map[string]any
	InputPath       string
	OutputPath      string
}

func NewBatchPrediction(jobName, modelName, inputPath, outputPath string, modelParameters map[string]any) *BatchPrediction {
	return &BatchPrediction{
		JobName:         jobName,
		ModelName:       modelName,
		ModelParameters: modelParameters,
		InputPath:       inputPath,
		OutputPath:      outputPath,
	}
}

func (b *BatchPrediction) GetID() string {
	return b.ID
}

func (b *BatchPrediction) SetID(id string) {
	b.ID = id
}
