package ai

type FileData struct {
	FileUri  string `json:"fileUri,omitempty"`
	MIMEType string `json:"mimeType,omitempty"`
}
type Parts struct {
	Text     string    `json:"text,omitempty"`
	FileData *FileData `json:"fileData,omitempty"`
}
type Content struct {
	Role  string  `json:"role"`
	Parts []Parts `json:"parts"`
}

type Request struct {
	Contents []Content `json:"contents"`
}

type PredictionRequest struct {
	Request Request `json:"request"`
}
