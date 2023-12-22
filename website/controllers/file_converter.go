package controllers

import (
	"net/http"

	"github.com/dingsongjie/file-help/configs"
	"github.com/dingsongjie/file-help/website/models/converter"
	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// GetFisrtImageByGavingKey
// @Summary GetFisrtImageByGavingKey
// @Description 根据文件key或者文件url获取文件并转成相应的目标图片，只转第一个图层或者第一页，目前支持psd->jpeg;ai->jpeg;pdf->jpeg
// @Tags GetFisrtImageByGavingKey
// @Accept json
// @Produce json
// @param request body converter.ConvertByGavingKeyRequest true "request"
// @Success 200  {object} converter.ConvertByGavingKeyResponse
// @Failure 400  {object} models.CommonErrorResponse
// @Router /Converter/GetFisrtImageByGavingKey [post]
func GetFisrtImageByGavingKey(c *gin.Context) {
	var request = converter.ConvertByGavingKeyRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	handler, err := converter.NewGetFisrtImageByGavingKeyRequestHandler(configs.S3Endpoint, configs.S3AccessKey, configs.S3SecretKey, configs.S3BacketName)
	if err != nil {
		c.AbortWithError(500, err)
	}
	GetFisrtImageByGavingKeyResponse := handler.Handle(&request)
	c.JSON(http.StatusOK, GetFisrtImageByGavingKeyResponse)
}

// @Security BearerAuth
// GetPdfByGavingKey
// @Summary GetPdfByGavingKey
// @Description 根据文件key或者文件url获取文件并转成相应的目标pdf，只转第一个图层或者第一页，目前支持psd->pdf;ai->pdf;pdf->pdf
// @Tags GetPdfByGavingKey
// @Accept json
// @Produce json
// @param request body converter.ConvertByGavingKeyRequest true "request"
// @Success 200  {object} converter.ConvertByGavingKeyResponse
// @Failure 400  {object} models.CommonErrorResponse
// @Router /Converter/GetPdfByGavingKey [post]
func GetPdfByGavingKey(c *gin.Context) {
	GetFisrtImageByGavingKey(c)
}
