package logic

import (
	"context"
	"regexp"
	"strings"

	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

var emailRegexp = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type UpdateUserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileLogic {
	return &UpdateUserProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserProfileLogic) UpdateUserProfile(in *user.UpdateUserProfileReq) (*user.UpdateUserProfileResp, error) {
	normalized, err := normalizeUpdateUserProfile(in)
	if err != nil {
		return nil, err
	}

	if err := l.svcCtx.CSvc.UpdateUserProfile(in.UserId, normalized.nickname, normalized.sex, normalized.email); err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "update user profile failed, userId=%s err=%v", in.UserId, err)
	}

	return &user.UpdateUserProfileResp{Success: true}, nil
}

type normalizedUpdateUserProfile struct {
	nickname string
	sex      int8
	email    string
}

func normalizeUpdateUserProfile(in *user.UpdateUserProfileReq) (*normalizedUpdateUserProfile, error) {
	nickname := strings.TrimSpace(in.Nickname)
	if nickname == "" {
		return nil, errors.WithStack(xerr.EmptyNickname)
	}
	if in.Sex != 0 && in.Sex != 1 && in.Sex != 2 {
		return nil, errors.WithStack(xerr.InvalidSex)
	}

	email := strings.TrimSpace(in.Email)
	if email != "" && !emailRegexp.MatchString(email) {
		return nil, errors.WithStack(xerr.InvalidEmail)
	}

	return &normalizedUpdateUserProfile{
		nickname: nickname,
		sex:      int8(in.Sex),
		email:    email,
	}, nil
}
