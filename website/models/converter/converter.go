package converter

import "strings"

type ConverteFirstAndReturnS3KeyRequest struct {
	Items []ConverteFirstAndReturnS3KeyItem `validate:"max=50"`
}

type ConverteFirstAndReturnS3KeyItem struct {
	sourceKey string `validate:"required"`
	targetKey string `validate:"required"`
}

func (r *ConverteFirstAndReturnS3KeyItem) Handle() (bool, error) {
	strings.Split(r.sourceKey, ".")
}
func ()

func (r *ConverteFirstAndReturnS3KeyItem) validate() (bool, error) {

}
