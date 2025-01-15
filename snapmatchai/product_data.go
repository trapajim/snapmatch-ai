package snapmatchai

type ProductData struct {
	ID         string
	Data       map[string]string
	AssetLinks []string
}

func (b *ProductData) GetID() string {
	return b.ID
}

func (b *ProductData) SetID(id string) {
	b.ID = id
}
