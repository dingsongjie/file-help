package imginfo

type GetImgInfoRequest struct {
	Items []GetImgInfoRequestItem `validate:"required,max=20,min=1"`
}
type GetImgInfoRequestItem struct {
	// @description minio文件key
	FileKey string
}
type GetImgInfoResponse struct {
	IsAllSucceed bool
	Items        []GetImgInfoItemResponse `validate:"required,max=20,min=1"`
}

type GetImgInfoItemResponse struct {
	FileKey       string `validate:"required"`
	Width, Height int
	Message       string
	IsSucceed     bool
}

func NewGetImgInfoResponse() *GetImgInfoResponse {
	return &GetImgInfoResponse{IsAllSucceed: true}
}

func (r *GetImgInfoResponse) AddItem(item *GetImgInfoItemResponse) {
	r.Items = append(r.Items, *item)
	if !item.IsSucceed {
		r.IsAllSucceed = false
	}
}
