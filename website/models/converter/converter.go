package converter

type ConvertByGavingKeyRequest struct {
	Items []ConvertByGavingKeyRequestItem `validate:"max=10,min=1"`
}

type ConvertByGavingKeyRequestItem struct {
	sourceKey string `validate:"required"`
	targetKey string `validate:"required"`
}

type ConvertByGavingKeyResponseItem struct {
	SourceKey string `validate:"required"`
	IsSucceed bool   `validate:"required"`
	Message   string
}

type ConvertByGavingKeyResponse struct {
	items        []ConvertByGavingKeyResponseItem `validate:"required"`
	isAllSucceed bool                             `validate:"required"`
}

func NewGetFisrtImageByGavingKeyResponse() *ConvertByGavingKeyResponse {
	return &ConvertByGavingKeyResponse{isAllSucceed: true}
}

func (r *ConvertByGavingKeyResponse) AddItem(item *ConvertByGavingKeyResponseItem) {
	r.items = append(r.items, *item)
	if !item.IsSucceed {
		r.isAllSucceed = false
	}
}
