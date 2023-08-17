package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"www.github.com/dingsongjie/file-help/website/models/converter"
)

// @BasePath /converter
// ConverteFirstAndReturnS3KeyRequest
// @Summary ConverteFirstAndReturnS3KeyRequest
// @Description ConverteFirstAndReturnS3KeyRequest
// @Tags ConverteFirstAndReturnS3KeyRequest
// @Accept json
// @Produce json
// @param request body converter.ConverteFirstAndReturnS3KeyRequest true "request"
// @Success 200  {object} model.CommandSuccessResponse
// @Failure 400  {object} model.CommonErrorResponse
// @Router /ConverteFisrtPageToGavingKey [post]
func ConverteFisrtPageToGavingKey(c *gin.Context) {
	var request = converter.ConverteFirstAndReturnS3KeyRequest{}
	if err := c.BindJSON(&request); err != nil {
		return
	}
	info, err := getNewPathInfo(requestUser.Paths)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatus(400)
	}
	err = setUserInfo(requestUser.Paths, info)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatus(400)
	} else {
		c.JSON(http.StatusOK, info)
	}
}

func convertCore(request converter.ConverteFirstAndReturnS3KeyRequest) {
	for i := range request.Items {

	}
}
