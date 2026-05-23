package logic

import (
	"context"
	"strings"

	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/internal/svc"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/models"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"
	"github.com/peninsula12/easy-im/go-im/pkg/encrypy"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(in *user.UpdateUserPasswordReq) (*user.UpdateUserPasswordResp, error) {
	oldPassword := strings.TrimSpace(in.OldPassword)
	newPassword := strings.TrimSpace(in.NewPassword)

	if oldPassword == "" {
		return nil, errors.WithStack(xerr.EmptyOldPassword)
	}
	if newPassword == "" {
		return nil, errors.WithStack(xerr.EmptyNewPassword)
	}
	if oldPassword == newPassword {
		return nil, errors.WithStack(xerr.PasswordUnchanged)
	}
	if len(newPassword) < 6 {
		return nil, errors.WithStack(xerr.PasswordTooShort)
	}

	currentUser := models.User{}
	if err := l.svcCtx.CSvc.DB.Where("id = ?", in.UserId).First(&currentUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithStack(xerr.IdNotFound)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by id failed, userId=%s err=%v", in.UserId, err)
	}

	if !encrypy.ValidatePasswordHash([]byte(currentUser.Password), []byte(oldPassword)) {
		return nil, errors.WithStack(xerr.UserPwdErr)
	}

	hashed, err := encrypy.GenPasswordHash([]byte(newPassword))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewServerCommonErr(), "passwordHash gen err %v", err)
	}

	if err := l.svcCtx.CSvc.DB.Model(&models.User{}).Where("id = ?", in.UserId).Update("password", string(hashed)).Error; err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "update user password failed, userId=%s err=%v", in.UserId, err)
	}

	return &user.UpdateUserPasswordResp{Success: true}, nil
}
