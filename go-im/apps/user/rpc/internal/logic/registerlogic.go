package logic

import (
	"context"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/models"
	"github.com/peninsula12/easy-im/go-im/pkg/ctxdata"
	"github.com/peninsula12/easy-im/go-im/pkg/encrypy"
	"github.com/peninsula12/easy-im/go-im/pkg/suid"
	"github.com/peninsula12/easy-im/go-im/pkg/utils"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
	"github.com/pkg/errors"
	"time"

	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	var err error
	// 1.检查手机号是否已被注册
	existing := models.User{}
	err = l.svcCtx.CSvc.GetUserByPhone(&existing, in.Phone)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find api by phone err %v req %v", err, in.Phone)
	}
	if existing.ID != "" {
		return nil, errors.WithStack(xerr.PhoneAlreadyRegistered)
	}

	// 2.定义新增用户
	U := &models.User{
		ID:       suid.GenerateID(),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Status:   utils.ConvertToInt8(0),
		Sex:      utils.ConvertToInt8(in.Sex),
	}

	if in.Password != "" {
		pass, err := encrypy.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, errors.Wrapf(xerr.NewServerCommonErr(), "passwordHash gen err %v", err)
		}
		U.Password = string(pass)
	}
	// 3.保存用户
	err = l.svcCtx.CSvc.CreateUser(U)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "save api %v failed ,err %v", in, err)
	}

	// 4. 生成token，使用新建用户的 ID
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, U.ID)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "extdata get jwt token"+
			" err %v", in.Phone)
	}
	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil

}
