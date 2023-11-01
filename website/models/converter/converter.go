package converter

type ConvertByGavingKeyRequest struct {
	Items []ConvertByGavingKeyRequestItem `validate:"required,max=10,min=1"`
}

type ConvertByGavingKeyRequestItem struct {
	// @description 如果这个值为合法url就会根据url地址获取源文件，如果不是则默认为文件key，通过minio获取文件
	SourceKey     string `validate:"required"`
	TargetKey     string `validate:"required"`
	TargetFileDpi int    `validate:"default=0"`
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
