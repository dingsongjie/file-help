package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"www.github.com/dingsongjie/file-help/configs"
	"www.github.com/dingsongjie/file-help/website/models/tar"
)

// @Security BearerAuth
// GetImgInfo
// @Summary GetImgInfo
// @Description GetImgInfo
// @Tags GetImgInfo
// @Accept json
// @Produce json
// @param request body tar.PackRequest  true "request"
// @Success 200  {object} models.CommandResponse
// @Failure 400  {object} models.CommonErrorResponse
// @Router /GetImgInfo [post]
func GetImgInfo(c *gin.Context) {
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
