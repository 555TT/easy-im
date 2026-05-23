package logic

import (
	"context"
	stderrors "errors"
	"regexp"
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

	var passwordHash *string
	if normalized.newPassword != "" {
		currentUser := models.User{}
		if err := l.svcCtx.CSvc.DB.Where("id = ?", in.UserId).First(&currentUser).Error; err != nil {
			if stderrors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.WithStack(xerr.IdNotFound)
			}
			return nil, errors.Wrapf(xerr.NewDBErr(), "find user by id failed, userId=%s err=%v", in.UserId, err)
		}

		passwordHash, err = validateAndHashNewPassword(currentUser.Password, normalized.oldPassword, normalized.newPassword)
		if err != nil {
			return nil, err
		}
	}

	if err := l.svcCtx.CSvc.UpdateUserProfile(in.UserId, normalized.nickname, normalized.sex, normalized.email, normalized.avatar, passwordHash); err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "update user profile failed, userId=%s err=%v", in.UserId, err)
	}

	return &user.UpdateUserProfileResp{Success: true}, nil
}

type normalizedUpdateUserProfile struct {
	nickname    string
	sex         int8
	email       string
	avatar      string
	oldPassword string
	newPassword string
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

	avatar := strings.TrimSpace(in.Avatar)
	oldPassword := strings.TrimSpace(in.OldPassword)
	newPassword := strings.TrimSpace(in.NewPassword)
	if (oldPassword == "") != (newPassword == "") {
		return nil, errors.WithStack(xerr.PasswordRequired)
	}
	if oldPassword != "" {
		if oldPassword == newPassword {
			return nil, errors.WithStack(xerr.PasswordUnchanged)
		}
		if len(newPassword) < 6 {
			return nil, errors.WithStack(xerr.PasswordTooShort)
		}
	}

	return &normalizedUpdateUserProfile{
		nickname:    nickname,
		sex:         int8(in.Sex),
		email:       email,
		avatar:      avatar,
		oldPassword: oldPassword,
		newPassword: newPassword,
	}, nil
}

func validateAndHashNewPassword(currentPasswordHash, oldPassword, newPassword string) (*string, error) {
	if !encrypy.ValidatePasswordHash([]byte(currentPasswordHash), []byte(oldPassword)) {
		return nil, errors.WithStack(xerr.UserPwdErr)
	}

	hashed, err := encrypy.GenPasswordHash([]byte(newPassword))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewServerCommonErr(), "passwordHash gen err %v", err)
	}
	hashString := string(hashed)
	return &hashString, nil
}
