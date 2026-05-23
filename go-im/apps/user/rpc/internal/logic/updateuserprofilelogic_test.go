package logic

import (
	stderrors "errors"
	"testing"

	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/user"
	"github.com/peninsula12/easy-im/go-im/pkg/encrypy"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
	zerrors "github.com/zeromicro/x/errors"
)

func TestNormalizeUpdateUserProfile_TrimmedValues(t *testing.T) {
	normalized, err := normalizeUpdateUserProfile(&user.UpdateUserProfileReq{
		Nickname: "  valid-name  ",
		Sex:      2,
		Email:    "  valid@example.com  ",
		Avatar:   "  https://example.com/avatar.png  ",
	})
	if err != nil {
		t.Fatalf("normalizeUpdateUserProfile() error = %v", err)
	}
	if normalized.nickname != "valid-name" {
		t.Fatalf("normalizeUpdateUserProfile() nickname = %q, want %q", normalized.nickname, "valid-name")
	}
	if normalized.sex != 2 {
		t.Fatalf("normalizeUpdateUserProfile() sex = %d, want 2", normalized.sex)
	}
	if normalized.email != "valid@example.com" {
		t.Fatalf("normalizeUpdateUserProfile() email = %q, want %q", normalized.email, "valid@example.com")
	}
	if normalized.avatar != "https://example.com/avatar.png" {
		t.Fatalf("normalizeUpdateUserProfile() avatar = %q, want %q", normalized.avatar, "https://example.com/avatar.png")
	}
}

func TestNormalizeUpdateUserProfile_RejectsEmptyNickname(t *testing.T) {
	_, err := normalizeUpdateUserProfile(&user.UpdateUserProfileReq{
		Nickname: "   ",
		Sex:      1,
		Email:    "valid@example.com",
	})
	assertCodeError(t, err, xerr.EmptyNickname)
}

func TestNormalizeUpdateUserProfile_RejectsInvalidSex(t *testing.T) {
	_, err := normalizeUpdateUserProfile(&user.UpdateUserProfileReq{
		Nickname: "valid-name",
		Sex:      3,
		Email:    "valid@example.com",
	})
	assertCodeError(t, err, xerr.InvalidSex)
}

func TestNormalizeUpdateUserProfile_RejectsInvalidEmail(t *testing.T) {
	_, err := normalizeUpdateUserProfile(&user.UpdateUserProfileReq{
		Nickname: "valid-name",
		Sex:      1,
		Email:    "not-an-email",
	})
	assertCodeError(t, err, xerr.InvalidEmail)
}

func TestNormalizeUpdateUserProfile_AllowsProfileOnlyUpdate(t *testing.T) {
	normalized, err := normalizeUpdateUserProfile(&user.UpdateUserProfileReq{
		Nickname: "valid-name",
		Sex:      1,
		Email:    "valid@example.com",
		Avatar:   "https://example.com/avatar.png",
	})
	if err != nil {
		t.Fatalf("normalizeUpdateUserProfile() error = %v", err)
	}
	if normalized.nickname != "valid-name" {
		t.Fatalf("normalizeUpdateUserProfile() nickname = %q, want %q", normalized.nickname, "valid-name")
	}
}

func TestValidateAndHashNewPassword_ReturnsErrorForWrongOldPassword(t *testing.T) {
	currentHash, err := encrypy.GenPasswordHash([]byte("current-password"))
	if err != nil {
		t.Fatalf("GenPasswordHash() error = %v", err)
	}

	hash, err := validateAndHashNewPassword(string(currentHash), "wrong-password", "new-password")
	if err == nil {
		t.Fatalf("validateAndHashNewPassword() error = nil, want %v", xerr.UserPwdErr)
	}
	if hash != nil {
		t.Fatalf("validateAndHashNewPassword() hash = %v, want nil", hash)
	}
	assertCodeError(t, err, xerr.UserPwdErr)
}

func assertCodeError(t *testing.T, err error, want error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error %v, got nil", want)
	}

	var codeErr *zerrors.CodeMsg
	if !stderrors.As(err, &codeErr) {
		t.Fatalf("expected code error, got %T: %v", err, err)
	}

	var wantCodeErr *zerrors.CodeMsg
	if !stderrors.As(want, &wantCodeErr) {
		t.Fatalf("expected code error for want, got %T: %v", want, want)
	}

	if codeErr.Code != wantCodeErr.Code {
		t.Fatalf("error code = %d, want %d (err=%v)", codeErr.Code, wantCodeErr.Code, err)
	}
}
