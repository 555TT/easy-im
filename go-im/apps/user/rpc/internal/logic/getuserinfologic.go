package logic

import (
	"context"
	stderrors "errors"

	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/models"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	u := make([]models.User, 0, 1)
	s := make([]string, 1)
	s[0] = in.User
	err := l.svcCtx.CSvc.GetUserByIds(&u, s)
	ur, err := getUserInfoResult(u, err, in.User)
	if err != nil {
		return nil, err
	}

	userEntity := buildUserEntity(ur)
	return &user.GetUserInfoResp{User: &userEntity}, nil
}

func getUserInfoResult(users []models.User, queryErr error, userID string) (models.User, error) {
	if queryErr != nil {
		if stderrors.Is(queryErr, gorm.ErrRecordNotFound) {
			return models.User{}, errors.WithStack(xerr.IdNotFound)
		}
		return models.User{}, errors.Wrapf(xerr.NewDBErr(), "find api by id err %v req %v", queryErr, userID)
	}
	if len(users) == 0 {
		return models.User{}, errors.WithStack(xerr.IdNotFound)
	}
	return users[0], nil
}

func buildUserEntity(ur models.User) user.UserEntity {
	return user.UserEntity{
		Id:       ur.ID,
		Avatar:   ur.Avatar,
		Nickname: ur.Nickname,
		Phone:    ur.Phone,
		Status:   int32(safeInt8(ur.Status)),
		Sex:      int32(safeInt8(ur.Sex)),
		Email:    ur.Email,
	}
}

func safeInt8(v *int8) int8 {
	if v == nil {
		return 0
	}
	return *v
}
