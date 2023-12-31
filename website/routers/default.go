package routers

import (
	"net/http"
	"time"

	"github.com/dingsongjie/file-help/pkg/log"

	docs "github.com/dingsongjie/file-help/api/swag"
	"github.com/dingsongjie/file-help/configs"
	"github.com/dingsongjie/file-help/website/controllers"
	"golang.org/x/oauth2"

	"github.com/STRockefeller/go-linq"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginoauth2 "github.com/zalando/gin-oauth2"
)

func AddRouter(r *gin.Engine) *gin.Engine {
	oauth2Enpoint := oauth2.Endpoint{AuthURL: configs.OIDCIntrospectEndpoint}
	ginoauth2.VarianceTimer = 3000 * time.Millisecond // defaults to 30s
	ginoauth2.ClientId = configs.OIDCClientId
	ginoauth2.ClientSecret = configs.OIDCClientSecret
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "health",
		})
	})
	fileheler := r.Group(configs.BaseUrl)
	{
		fileheler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		fileheler.POST("/Converter/GetFisrtImageByGavingKey", ginoauth2.Auth(AudAndScopeCheck("default", configs.OIDCAudience, configs.OIDCScope), oauth2Enpoint), controllers.GetFisrtImageByGavingKey)
		fileheler.POST("/Converter/GetPdfByGavingKey", ginoauth2.Auth(AudAndScopeCheck("default", configs.OIDCAudience, configs.OIDCScope), oauth2Enpoint), controllers.GetPdfByGavingKey)
		fileheler.POST("/Zip/Pack", ginoauth2.Auth(AudAndScopeCheck("default", configs.OIDCAudience, configs.OIDCScope), oauth2Enpoint), controllers.Pack)
		fileheler.POST("/Tar/Pack", ginoauth2.Auth(AudAndScopeCheck("default", configs.OIDCAudience, configs.OIDCScope), oauth2Enpoint), controllers.Pack)
		fileheler.POST("/GetImgInfo", ginoauth2.Auth(AudAndScopeCheck("default", configs.OIDCAudience, configs.OIDCScope), oauth2Enpoint), controllers.GetImgInfo)
	}

	// swagger
	docs.SwaggerInfo.BasePath = configs.BaseUrl

	return r
}

func AudAndScopeCheck(name, aud string, scopes ...string) func(tc *ginoauth2.TokenContainer, ctx *gin.Context) bool {
	log.Logger.Sugar().Debugf("scopeCheck %s configured to grant access only if scopes: %v are present", name, scopes)
	configuredScopes := scopes
	return func(tc *ginoauth2.TokenContainer, ctx *gin.Context) bool {
		audLinq := linq.Linq[string]{}
		audLinq.AddRange(tc.Aud)
		if !audLinq.Exists(func(s string) bool {
			return s == aud
		}) {
			return false
		}
		for _, s := range configuredScopes {
			if cur, ok := tc.Scopes[s]; ok {
				log.Logger.Sugar().Debugf("Found configured scope %s", s)
				ctx.Set(s, cur) // set value from token of configured scope to the context, which you can use in your application.
			} else {
				return false
			}
		}
		//Getting the uid for identification of the service calling
		if cur, ok := tc.Scopes["uid"]; ok {
			ctx.Set("uid", cur)
		}
		return true
	}
}
