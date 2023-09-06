package tar

import "time"

type PackRequest struct {
	FileKey  string
	IsGziped bool
	Items    *[]PackRequestItem
}

type PackRequestItem struct {
	FileKey, FileName string
	LastModifyTime    time.Time
}
