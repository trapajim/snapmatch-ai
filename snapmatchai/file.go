package snapmatchai

import "time"

type Metadata struct {
	Key   string `bigquery:"key" json:"key"`
	Value string `bigquery:"value" json:"value"`
}
type FileRecord struct {
	URI         string     `bigquery:"uri" json:"uri"`
	SignedURL   string     `bigquery:"signed_url" json:"signedURL"`
	ContentType string     `bigquery:"content_type" json:"contentType"`
	Size        int        `bigquery:"size" json:"size"`
	Updated     time.Time  `bigquery:"updated" json:"updated"`
	Metadata    []Metadata `bigquery:"metadata" json:"metadata"`
	Category    string     `bigquery:"category" json:"category"`
	ObjPath     string     `bigquery:"obj_path" json:"objPath"`
	ObjName     string     `bigquery:"obj_name" json:"objName"`
	FileID      string     `bigquery:"file_id" json:"fileID"`
	Distance    float64    `bigquery:"distance" json:"distance"`
}
