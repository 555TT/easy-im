package user

import (
	"context"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/peninsula12/easy-im/go-im/apps/user/api/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/user/api/internal/types"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"
	"github.com/peninsula12/easy-im/go-im/pkg/status"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLoginLogic 用户登入
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	phone := strings.TrimSpace(req.Phone)
	password := strings.TrimSpace(req.Password)
	if phone == "" || password == "" {
		return nil, errors.WithStack(xerr.ParamError)
	}

	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    phone,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	var res types.LoginResp
	err = copier.Copy(&res, loginResp)
	if err != nil {
		return nil, err
	}
	resp = &res

	// 处理登入的业务
	err = l.svcCtx.Redis.HsetCtx(l.ctx, status.REDIS_ONLINE_USER, loginResp.Id, "1")

	// 为每个用户的在线状态单独设置过期时间
	//expireKey := fmt.Sprintf("%s:%s", status.REDIS_ONLINE_USER, loginResp.Id)
	//expireErr := l.svcCtx.Redis.ExpireCtx(l.ctx, expireKey, 30*60)
	//if expireErr != nil {
	//	return nil, expireErr
	//}

	return
}
