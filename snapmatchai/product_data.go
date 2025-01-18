package snapmatchai

import "cloud.google.com/go/firestore"

type ProductData struct {
	ID         string
	Data       map[string]string
	VectorData firestore.Vector32
	AssetLinks []string
}

func (b *ProductData) GetID() string {
	return b.ID
}

func (b *ProductData) SetID(id string) {
	b.ID = id
}
