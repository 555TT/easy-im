package friend

import (
	"context"
	"regexp"

	"github.com/peninsula12/easy-im/go-im/apps/social/api/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/social/api/internal/types"
	"github.com/peninsula12/easy-im/go-im/apps/social/rpc/social"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"
	"github.com/peninsula12/easy-im/go-im/pkg/ctxdata"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

// 中国大陆手机号：11 位，1 开头，第二位 3-9
var phoneRegexp = regexp.MustCompile(`^1[3-9]\d{9}$`)

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFriendPutInLogic 好友申请
func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {
	uid := ctxdata.GetUId(l.ctx)

	// 1. 手机号基本格式校验
	if !phoneRegexp.MatchString(req.Phone) {
		return nil, errors.WithStack(xerr.ParamError)
	}
	logx.Infof("问题排查phone: %s", req.Phone)
	// 2. 通过手机号查目标用户
	findResp, err := l.svcCtx.User.FindUser(l.ctx, &user.FindUserReq{
		Phone: req.Phone,
	})
	if err != nil {
		return nil, err
	}
	logx.Infof("问题排查findResp: %+v", findResp)
	if len(findResp.Users) == 0 || findResp.Users[0].Id == "" {
		return nil, errors.WithStack(xerr.FriendPhoneNotRegistered)
	}
	targetUid := findResp.Users[0].Id

	// 3. 不能添加自己
	if targetUid == uid {
		return nil, errors.WithStack(xerr.CannotAddSelf)
	}

	// 4. 走原有 social-rpc 流程
	_, err = l.svcCtx.Social.FriendPutIn(l.ctx, &social.FriendPutInReq{
		UserId:  uid,
		ReqUid:  targetUid,
		ReqMsg:  req.ReqMsg,
		ReqTime: req.ReqTime,
	})
	if err != nil {
		return nil, err
	}
	return
}
