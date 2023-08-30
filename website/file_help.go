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

func main() {
	var (
		ginMode      string
		s3ServiceUrl string
		s3AccessKey  string
		s3SecretKey  string
		s3BacketName string
	)

	log.Initialise()
	godotenv.Load(".env")
	flag.StringVar(&ginMode, "gin-mode", "release", "Gin mode")
	flag.StringVar(&s3ServiceUrl, "s3-endpoint", "", "s3 service url")
	flag.StringVar(&s3AccessKey, "s3-access-key", "", "s3 account accesskey")
	flag.StringVar(&s3SecretKey, "s3-secret-key", "", "s3 account secretkey")
	flag.StringVar(&s3BacketName, "s3-bucket-name", "", "s3 target bucket name")
	flag.Parse()
	configs.ConfigGin(flag.CommandLine)
	configs.ConfigS3(flag.CommandLine)
	Initialize.RegisterConverters()
	defer Initialize.DestoryConverters()
	logger := log.Logger
	defer logger.Sync()

	gin.SetMode(configs.GinMode)
	r := gin.Default()
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
