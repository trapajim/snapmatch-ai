package models

type BarChart struct {
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
	Colors []string `json:"colors"`
}
