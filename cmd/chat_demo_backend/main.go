package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	api "go_messaging/cmd/chat_demo_backend/app"
	"go_messaging/internal/pubsub"
	"go_messaging/internal/ws_connect"
	"go_messaging/pkg/logger"
	"go_messaging/pkg/unique_id"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var roomManager *ws_connect.RoomManager

func serveHome(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "", []byte("websocket server"))
}

// serveWs handles websocket requests from the peer.
func ServeWs(ctx *gin.Context) {

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.GetLogger().Error("ws Upgrade err", logger.Error(err))
		return
	}

	// create new room
	roomId := fmt.Sprintf("%d", unique_id.GenerateID().Int64())
	room := roomManager.FindRoom(roomId)
	go func() {
		room.Run()
	}()

	// client
	client := ws_connect.NewClient(room, conn, make(chan []byte, 256))
	client.Room.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()

}

func main() {

	// config init
	upCfg := make(chan struct{})
	err := WatchConfig(upCfg)
	if err != nil {
		panic(err)
	}
	defer close(upCfg)

	// load config
	cfg, err := LoadConfig()
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

	// ws connect
	roomManager = ws_connect.NewRoomManager()

	// rabbit mq connect
	conn, err := amqp.Dial(cfg.RabbitMQ.DialUrl)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to RabbitMQ,err:%s", err.Error()))
	}
	defer func() {
		_ = conn.Close()
	}()

	ch, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to open a channel RabbitMQ,err:%s", err.Error()))
	}
	defer func() {
		_ = ch.Close()
	}()

	consumeCh, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to open a consume channel RabbitMQ,err:%s", err.Error()))
	}
	defer func() {
		_ = consumeCh.Close()
	}()

	// mq exchange
	if err = ch.ExchangeDeclare(
		cfg.RabbitMQ.Exchange.Name,        // name
		cfg.RabbitMQ.Exchange.Type,        // type
		cfg.RabbitMQ.Exchange.Durable,     // durable
		cfg.RabbitMQ.Exchange.AutoDeleted, // auto-deleted
		cfg.RabbitMQ.Exchange.Internal,    // internal
		cfg.RabbitMQ.Exchange.NoWait,      // no-wait
		nil,                               // arguments
	); err != nil {
		panic(fmt.Sprintf("Failed to declare an exchange RabbitMQ,err:%s", err.Error()))
	}

	if err = consumeCh.ExchangeDeclare(
		cfg.RabbitMQ.Exchange.Name,        // name
		cfg.RabbitMQ.Exchange.Type,        // type
		cfg.RabbitMQ.Exchange.Durable,     // durable
		cfg.RabbitMQ.Exchange.AutoDeleted, // auto-deleted
		cfg.RabbitMQ.Exchange.Internal,    // internal
		cfg.RabbitMQ.Exchange.NoWait,      // no-wait
		nil,                               // arguments
	); err != nil {
		panic(fmt.Sprintf("Failed to declare an consume exchange RabbitMQ,err:%s", err.Error()))
	}

	// init mq publisher
	rm := NewMQPublisher(cfg, ch)
	pubsub.InitDefaultPublisher(rm)

	// new consume logic
	cl := NewConsumeLogic(roomManager)
	// run mq consumer

	consume, err := NewMQConsumer(cfg, consumeCh, cl)
	if err != nil {
		panic(fmt.Sprintf("Failed to New Consumer,err:%s", err.Error()))
	}

	if err := consume.Run(); err != nil {
		panic(fmt.Sprintf("Failed to Run Consumer,err:%s", err.Error()))
	}

	// init upgrader
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// setup gin
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/", serveHome)
	router.GET("/ws", ServeWs)

	api.RegisterRouter(router)

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
