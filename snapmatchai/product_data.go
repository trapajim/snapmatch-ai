package snapmatchai

import "cloud.google.com/go/firestore"

type ProductData struct {
	ID         string
	Data       map[string]string
	VectorData firestore.Vector32
	Owner      string
	AssetLinks []string
}

func (b *ProductData) GetID() string {
	return b.ID
}

func (b *ProductData) SetID(id string) {
	b.ID = id
}
func (b *ProductData) GetVectorData() firestore.Vector32 {
	return b.VectorData
}

func (b *ProductData) SetOwner(s string) {
	b.Owner = s
}
