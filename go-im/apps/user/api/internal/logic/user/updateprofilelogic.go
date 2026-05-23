package user

import (
	"context"

	"github.com/peninsula12/easy-im/go-im/apps/user/api/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/user/api/internal/types"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"
	"github.com/peninsula12/easy-im/go-im/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateProfileLogic 更新用户资料
func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileReq) (resp *types.UpdateProfileResp, err error) {
	uid := ctxdata.GetUId(l.ctx)
	_, err = l.svcCtx.User.UpdateUserProfile(l.ctx, &user.UpdateUserProfileReq{
		UserId:   uid,
		Nickname: req.Nickname,
		Sex:      int32(req.Sex),
		Email:    req.Email,
		Avatar:   req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateProfileResp{Success: true}, nil
}
