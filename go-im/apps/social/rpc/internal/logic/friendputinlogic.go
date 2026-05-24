package logic

import (
	"context"
	"time"

	"github.com/peninsula12/easy-im/go-im/apps/social/rpc/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/social/rpc/models"
	"github.com/peninsula12/easy-im/go-im/pkg/status"
	"github.com/peninsula12/easy-im/go-im/pkg/suid"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/peninsula12/easy-im/go-im/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	var friend models.Friend
	err := l.svcCtx.CSvc.DB.Where("user_id = ? and friend_uid = ?", in.UserId, in.ReqUid).First(&friend).Error
	if err == nil {
		return nil, errors.WithStack(xerr.FriendAlreadyExists)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend by user_id %v and friend_uid %v err %v", in.UserId, in.ReqUid, err)
	}

	var friendReq models.FriendRequest
	err = l.svcCtx.CSvc.DB.Where("user_id = ? and req_uid = ? and handle_result = ?", in.ReqUid, in.UserId, status.PendingHandlerResult).First(&friendReq).Error
	if err == nil {
		return nil, errors.WithStack(xerr.FriendRequestOnPending)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friend request by user_id %v and req_uid %v err %v", in.ReqUid, in.UserId, err)
	}

	err = l.svcCtx.CSvc.DB.Create(&models.FriendRequest{
		ID:           suid.GenerateID(),
		UserID:       in.ReqUid,
		ReqUID:       in.UserId,
		ReqMsg:       in.ReqMsg,
		ReqTime:      time.UnixMilli(in.ReqTime),
		HandleResult: status.PendingHandlerResult,
	}).Error
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "create friend request by user_id %v and req_uid %v, rawErr: %v", in.UserId, in.ReqUid, err)
	}
	return &social.FriendPutInResp{}, nil
}
