package main

import (
	"flag"
	"fmt"

	"github.com/peninsula12/easy-im/go-im/pkg/resultx"

	"github.com/peninsula12/easy-im/go-im/apps/user/api/internal/config"
	"github.com/peninsula12/easy-im/go-im/apps/user/api/internal/handler"
	"github.com/peninsula12/easy-im/go-im/apps/user/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	Run(c)
}

func Run(c config.Config) {
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	httpx.SetOkHandler(resultx.OKHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
