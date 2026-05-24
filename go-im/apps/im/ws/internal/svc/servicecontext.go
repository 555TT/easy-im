package svc

import (
	"github.com/peninsula12/easy-im/go-im/apps/im/model"
	"github.com/peninsula12/easy-im/go-im/apps/im/ws/internal/config"
	"github.com/peninsula12/easy-im/go-im/apps/task/mq/mqclient"
	"github.com/peninsula12/easy-im/go-im/pkg/status"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"golang.org/x/net/context"
)

type ServiceContext struct {
	Config config.Config
	mqclient.MsgChatTransferClient
	mqclient.MsgReadTransferClient
	immodels.ChatLogModel
	*redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(
			c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		MsgReadTransferClient: mqclient.NewMsgReadTransferClient(
			c.MsgReadTransfer.Addrs, c.MsgReadTransfer.Topic),
		ChatLogModel: immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		Redis:        redis.MustNewRedis(c.Redisx),
	}
}

func (s *ServiceContext) MarkOnline(ctx context.Context, uid string) error {
	return s.Redis.HsetCtx(ctx, status.REDIS_ONLINE_USER, uid, "1")
}

func (s *ServiceContext) MarkOffline(ctx context.Context, uid string) error {
	_, err := s.Redis.HdelCtx(ctx, status.REDIS_ONLINE_USER, uid)
	return err
}
