package models

import "time"

type Asset struct {
	Name     string
	Size     int
	Type     string
	Date     time.Time
	Category string
	URI      string
}
type Assets []Asset
