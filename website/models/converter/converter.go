package converter

type ConvertByGavingKeyRequest struct {
	Items []ConvertByGavingKeyRequestItem `validate:"required,max=10,min=1"`
}

type ConvertByGavingKeyRequestItem struct {
	SourceKey string `validate:"required"`
	TargetKey string `validate:"required"`
}

type ConvertByGavingKeyResponseItem struct {
	SourceKey      string `validate:"required"`
	IsSucceed      bool   `validate:"required"`
	TargetFileSize int64
	Message        string
}

type ConvertByGavingKeyResponse struct {
	Items        []ConvertByGavingKeyResponseItem `validate:"required"`
	IsAllSucceed bool                             `validate:"required"`
}

func NewGetFisrtImageByGavingKeyResponse() *ConvertByGavingKeyResponse {
	return &ConvertByGavingKeyResponse{IsAllSucceed: true}
}

func (r *ConvertByGavingKeyResponse) AddItem(item *ConvertByGavingKeyResponseItem) {
	r.Items = append(r.Items, *item)
	if !item.IsSucceed {
		r.IsAllSucceed = false
	}
}
