package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"www.github.com/dingsongjie/file-help/configs"
	"www.github.com/dingsongjie/file-help/website/models/converter"
)

// @Security BearerAuth
// @BasePath /
// GetFisrtImageByGavingKey
// @Summary GetFisrtImageByGavingKey
// @Description GetFisrtImageByGavingKey
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

// @BasePath /
// GetPdfByGavingKey
// @Summary GetPdfByGavingKey
// @Description GetPdfByGavingKey
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
