package logic

import (
	"context"
	"fmt"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/models"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	var users []models.User

	if in.Phone != "" {
		var u models.User
		err := l.svcCtx.CSvc.GetUserByPhone(&u, in.Phone)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find api by phone: %s", in.Phone)
		}
		if u.ID == "" {
			// 查无此手机号，返回空列表，由调用方决定如何处理
			return &user.FindUserResp{Users: []*user.UserEntity{}}, nil
		}
		users = append(users, u)
	} else if len(in.Ids) > 0 {
		err := l.svcCtx.CSvc.GetUserByIds(&users, in.Ids)
		if err != nil {
			fmt.Printf("\n\n\n %v \n\n\n", err)
			return nil, errors.Wrapf(err, "failed to find users by IDs: %v", in.Ids)
		}
	} else if in.Name != "" {
		var u models.User
		err := l.svcCtx.CSvc.GetUserByName(&u, in.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find users by name: %s", in.Name)
		}
		if u.ID == "" {
			return &user.FindUserResp{Users: []*user.UserEntity{}}, nil
		}
		users = append(users, u)
	} else {
		return nil, errors.WithStack(xerr.ParamError)
	}

	userEntities := make([]*user.UserEntity, 0, len(users))

	for _, u := range users {
		var status, sex int32
		if u.Status != nil {
			status = int32(*u.Status)
		}
		if u.Sex != nil {
			sex = int32(*u.Sex)
		}
		userEntities = append(userEntities, &user.UserEntity{
			Id:       u.ID,
			Avatar:   u.Avatar,
			Nickname: u.Nickname,
			Phone:    u.Phone,
			Status:   status,
			Sex:      sex,
		})
	}

	return &user.FindUserResp{Users: userEntities}, nil
}
