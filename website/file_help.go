package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/namsral/flag"
	"go.uber.org/zap"
	"www.github.com/dingsongjie/file-help/configs"
	"www.github.com/dingsongjie/file-help/pkg/converter/Initialize"
	"www.github.com/dingsongjie/file-help/pkg/log"
	"www.github.com/dingsongjie/file-help/website/routers"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	var (
		ginMode, s3ServiceUrl, s3AccessKey, s3SecretKey, s3BacketName, oidcClientId, oidcClientSecret, oidcAuthority,
		oidcScope, oidcAudience, oidcIntrospectEndpoint string
	)

	log.Initialise()
	godotenv.Load(".env")
	flag.StringVar(&ginMode, "gin-mode", "release", "Gin mode")
	flag.StringVar(&s3ServiceUrl, "s3-endpoint", "", "s3 service url")
	flag.StringVar(&s3AccessKey, "s3-access-key", "", "s3 account accesskey")
	flag.StringVar(&s3SecretKey, "s3-secret-key", "", "s3 account secretkey")
	flag.StringVar(&s3BacketName, "s3-bucket-name", "", "s3 target bucket name")

	flag.StringVar(&oidcClientId, "oidc-client-id", "", "oidc clientId")
	flag.StringVar(&oidcClientSecret, "oidc-client-secret", "", "oidc clientSecret")
	flag.StringVar(&oidcAuthority, "oidc-authority", "", "oidc authority")
	flag.StringVar(&oidcScope, "oidc-scope", "", "oidc scope")
	flag.StringVar(&oidcAudience, "oidc-audience", "", "oidc audience")
	flag.StringVar(&oidcIntrospectEndpoint, "oidc-introspect-endpoint", "", "oidc introspect enpoint")
	flag.Parse()
	configs.ConfigGin(flag.CommandLine)
	configs.ConfigS3(flag.CommandLine)
	configs.ConfigIdentityServer(flag.CommandLine)
	Initialize.RegisterConverters()
	defer Initialize.DestoryConverters()
	logger := log.Logger
	defer logger.Sync()

	gin.SetMode(configs.GinMode)
	r := gin.Default()

	// store := cookie.NewStore([]byte("secret"))
	// r.Use(sessions.Sessions("mysession", store))
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	r.Use(ginzap.RecoveryWithZap(logger, true))
	routers.AddRouter(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	Initialize.DestoryConverters()
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown:", zap.Error(err))
	}
	logger.Info("Server exiting")
}
