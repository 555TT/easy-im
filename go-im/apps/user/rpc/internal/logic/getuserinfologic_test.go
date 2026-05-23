package logic

import (
	stderrors "errors"
	"testing"

	"gorm.io/gorm"
	"github.com/peninsula12/easy-im/go-im/apps/user/rpc/models"
	"github.com/peninsula12/easy-im/go-im/pkg/xerr"
)

func TestGetUserInfoResult_ReturnsNotFoundForEmptyRows(t *testing.T) {
	_, err := getUserInfoResult(nil, nil, "missing-user-id")
	assertCodeError(t, err, xerr.IdNotFound)
}

func TestGetUserInfoResult_ReturnsNotFoundWhenQueryReturnsRecordNotFound(t *testing.T) {
	_, err := getUserInfoResult(nil, gorm.ErrRecordNotFound, "missing-user-id")
	assertCodeError(t, err, xerr.IdNotFound)
}

func TestGetUserInfoResult_ReturnsDBErrorWhenQueryFails(t *testing.T) {
	_, err := getUserInfoResult(nil, stderrors.New("db down"), "user-id")
	assertCodeError(t, err, xerr.NewDBErr())
}

func TestGetUserInfoResult_ReturnsFirstUserWhenQuerySucceeds(t *testing.T) {
	want := models.User{ID: "user-id", Email: "user@example.com"}

	got, err := getUserInfoResult([]models.User{want}, nil, want.ID)
	if err != nil {
		t.Fatalf("getUserInfoResult() error = %v", err)
	}
	if got.ID != want.ID {
		t.Fatalf("getUserInfoResult() id = %q, want %q", got.ID, want.ID)
	}
	if got.Email != want.Email {
		t.Fatalf("getUserInfoResult() email = %q, want %q", got.Email, want.Email)
	}
}

func TestBuildUserEntity_UsesZeroValuesForNilStatusAndSex(t *testing.T) {
	got := buildUserEntity(models.User{ID: "user-id", Email: "user@example.com"})

	if got.Id != "user-id" {
		t.Fatalf("buildUserEntity() id = %q, want %q", got.Id, "user-id")
	}
	if got.Email != "user@example.com" {
		t.Fatalf("buildUserEntity() email = %q, want %q", got.Email, "user@example.com")
	}
	if got.Status != 0 {
		t.Fatalf("buildUserEntity() status = %d, want 0", got.Status)
	}
	if got.Sex != 0 {
		t.Fatalf("buildUserEntity() sex = %d, want 0", got.Sex)
	}
}

func TestBuildUserEntity_MapsStatusAndSexValues(t *testing.T) {
	status := int8(1)
	sex := int8(2)

	got := buildUserEntity(models.User{ID: "user-id", Status: &status, Sex: &sex})

	if got.Status != int32(status) {
		t.Fatalf("buildUserEntity() status = %d, want %d", got.Status, status)
	}
	if got.Sex != int32(sex) {
		t.Fatalf("buildUserEntity() sex = %d, want %d", got.Sex, sex)
	}
}

