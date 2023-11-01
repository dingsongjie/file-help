package tar

import "time"

type PackRequest struct {
	FileKey  string
	IsGziped bool
	Items    *[]PackRequestItem
}

type PackRequestItem struct {
	// @description 如果此参数是一个合法的url则会根据url获取文件，否则视为minio文件key
	FileKey string
	// 文件夹层级用"/"表示
	FileName string
	// ISO 8601 时间
	LastModifyTime time.Time
}
