package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/peninsula12/easy-im/go-im/apps/im/ws/internal/config"
	"github.com/peninsula12/easy-im/go-im/apps/im/ws/internal/handler"
	"github.com/peninsula12/easy-im/go-im/apps/im/ws/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/im/ws/websocket"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithOnlineTracker(ctx),
		//websocket.WithServerMaxConnectionIdle(180*time.Second),
		websocket.WithServerMaxConnectionIdle(180*time.Minute),
		websocket.WithServerAck(websocket.OnlyAck),
	)
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Printf("starting server at %v ...\n", c.ListenOn)
	srv.Start()
}
