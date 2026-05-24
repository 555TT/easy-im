package push

import (
	"github.com/mitchellh/mapstructure"
	"github.com/peninsula12/easy-im/go-im/apps/im/ws/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/im/ws/websocket"
	"github.com/peninsula12/easy-im/go-im/apps/im/ws/ws"
	"github.com/peninsula12/easy-im/go-im/pkg/status"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			_ = srv.Send(websocket.NewErrMessage(err))
			return
		}
		switch data.ChatType {
		case status.SingleChatType:
			for _, uid := range []string{data.SendId, data.RecvId} {
				err := single(srv, &data, uid)
				if err != nil {
					srv.Error(err)
				}
			}
		case status.GroupChatType:
			err := group(srv, &data)
			if err != nil {
				srv.Error(err)
			}
		default:
		}
	}
}

func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	recvConn := srv.GetConn(recvId)
	if len(recvConn) == 0 || recvConn[0] == nil {
		return nil
	}
	srv.Infof("push msg %v", data)

	return srv.Send(&websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    data.SendId,
		Data: map[string]any{
			"conversationId": data.ConversationId,
			"chatType":       data.ChatType,
			"sendId":         data.SendId,
			"recvId":         recvId,
			"recvIds":        data.RecvIds,
			"sendTime":       data.SendTime,
			"msgId":          data.MsgId,
			"readRecords":    data.ReadRecords,
			"contentType":    data.ContentType,
			"mType":          data.MType,
			"content":        data.Content,
		},
	}, recvConn[0])
}

func group(srv *websocket.Server, data *ws.Push) (err error) {
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				err = single(srv, data, id)
			})
		}(id)
	}
	return
}
