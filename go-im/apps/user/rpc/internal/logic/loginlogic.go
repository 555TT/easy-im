package logic

import (
	"context"
	"strings"
	"time"

	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/models"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"
	"github.com/peninsula12/easy-im/go-im/pkg/ctxdata"
	"github.com/peninsula12/easy-im/go-im/pkg/encrypy"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	phone := strings.TrimSpace(in.Phone)
	password := strings.TrimSpace(in.Password)
	if phone == "" || password == "" {
		return nil, errors.WithStack(xerr.ParamError)
	}

	u := &models.User{}
	err := l.svcCtx.CSvc.GetUserByPhone(u, phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(xerr.PhoneNotFound)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find api by phone %v err %v ", phone, err)
	}

	if !encrypy.ValidatePasswordHash([]byte(u.Password), []byte(password)) {
		return nil, errors.WithStack(xerr.UserPwdErr)
	}
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, u.ID)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "etxdata get jwt token"+
			" err %v ", in.Phone)
	}
	return &user.LoginResp{
		Id:     u.ID,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
