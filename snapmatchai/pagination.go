package snapmatchai

import (
	"encoding/base64"
	"errors"
)

type Pagination struct {
	NextToken string `json:"nextToken"`
	Per       int    `json:"per"`
}

func NewPagination(nextToken string, per int) *Pagination {
	encodedToken := base64.StdEncoding.EncodeToString([]byte(nextToken))
	return &Pagination{
		NextToken: encodedToken,
		Per:       per,
	}
}

func (p *Pagination) DecodeNextToken() (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(p.NextToken)
	if err != nil {
		return "", errors.New("failed to decode NextToken")
	}
	return string(decodedBytes), nil
}
