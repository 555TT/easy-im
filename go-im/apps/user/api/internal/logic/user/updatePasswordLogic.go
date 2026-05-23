package user

import (
	"context"

	"github.com/peninsula12/easy-im/go-im/apps/user/api/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/user/api/internal/types"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"
	"github.com/peninsula12/easy-im/go-im/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改密码
func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(req *types.UpdatePasswordReq) (resp *types.UpdatePasswordResp, err error) {
	uid := ctxdata.GetUId(l.ctx)
	_, err = l.svcCtx.User.UpdateUserPassword(l.ctx, &user.UpdateUserPasswordReq{
		UserId:      uid,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdatePasswordResp{Success: true}, nil
}
