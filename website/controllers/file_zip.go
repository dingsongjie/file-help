package controllers

import (
	"net/http"

	"github.com/dingsongjie/file-help/configs"
	"github.com/dingsongjie/file-help/website/models/tar"
	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// PackByGavingKey
// @Summary PackByGavingKey
// @Description 根据文件key或者文件url归档，或者归档压缩生成文件,目前只支持归档成zip或者zip.gz文件
// @Tags PackByGavingKey
// @Accept json
// @Produce json
// @param request body tar.PackRequest  true "request"
// @Success 200  {object} models.CommandResponse
// @Failure 400  {object} models.CommonErrorResponse
// @Router /Zip/Pack [post]
// @Router /Tar/Pack [post]
func Pack(c *gin.Context) {
	var request = tar.PackRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	handler, err := tar.NewPackHandler(configs.S3Endpoint, configs.S3AccessKey, configs.S3SecretKey, configs.S3BacketName)
	if err != nil {
		c.AbortWithError(500, err)
	}
	response := handler.Handle(&request)
	c.JSON(http.StatusOK, response)
}
