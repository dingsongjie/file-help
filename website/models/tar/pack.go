package tar

import "time"

type PackRequest struct {
	FileKey  string
	IsGziped bool
	Items    *[]PackRequestItem
}

type PackRequestItem struct {
	FileKey string
	// 文件夹层级用"/"表示
	FileName string
	// ISO 8601 时间
	LastModifyTime time.Time
}
