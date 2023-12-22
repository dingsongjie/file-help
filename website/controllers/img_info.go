package controllers

import (
	"net/http"

	"github.com/dingsongjie/file-help/website/models/imginfo"
	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// GetImgInfo
// @Summary GetImgInfo
// @Description 批量获取图片宽高，只支持png,jpg,gif类型图片
// @Tags GetImgInfo
// @Accept json
// @Produce json
// @param request body imginfo.GetImgInfoRequest  true "request"
// @Success 200  {object} imginfo.GetImgInfoResponse
// @Failure 400  {object} models.CommonErrorResponse
// @Router /GetImgInfo [post]
func GetImgInfo(c *gin.Context) {
	var request = imginfo.GetImgInfoRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	query := imginfo.NewImgInfoQueries()
	response := query.GetImgInfo(&request)
	c.JSON(http.StatusOK, response)
}
