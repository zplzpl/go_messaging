package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go_messaging/pkg/logger"
)

func serveHome(ctx *gin.Context) {

	if ctx.Request.URL.Path != "/" {
		http.Error(ctx.Writer, "Not found", http.StatusNotFound)
		return
	}
	if ctx.Request.Method != "GET" {
		http.Error(ctx.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx.HTML(http.StatusOK, "home.html", gin.H{
		"ws_server": cfg.WebsocketServer.WsUrl,
	})
}

var cfg *ConfigParser

func main() {

	// config init
	upCfg := make(chan struct{})
	err := WatchConfig(upCfg)
	if err != nil {
		panic(err)
	}
	defer close(upCfg)

	// load config
	cfg, err = LoadConfig()
	if err != nil {
		panic(err)
	}

	// init logger
	if err := logger.InitDefaultLogger(true); err != nil {
		panic(err)
	}
	defer logger.DeaultLoggerSync()

	// print
	logger.GetLogger().Info("load config file", logger.Any("cfg", cfg))

	// setup gin
	router := gin.Default()
	router.Use(cors.Default())
	router.LoadHTMLFiles("home.html")
	router.GET("/", serveHome)

	// run gin server
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}

}
