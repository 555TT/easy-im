package test

import (
	"context"
	"github.com/peninsula12/easy-im/go-im/apps/social/rpc/internal/logic"
	"github.com/peninsula12/easy-im/go-im/apps/social/rpc/social"
	"testing"
)

func TestFriendListLogic_FriendList(t *testing.T) {
	t.Run("returns empty list when user has no friends", func(t *testing.T) {
		l := logic.NewFriendListLogic(context.Background(), svcCtx)
		got, err := l.FriendList(&social.FriendListReq{UserId: "user-with-no-friends"})
		if err != nil {
			t.Fatalf("FriendList() error = %v, want nil", err)
		}
		if got == nil {
			t.Fatal("FriendList() got nil response")
		}
		if len(got.List) != 0 {
			t.Fatalf("FriendList() list length = %d, want 0", len(got.List))
		}
	})
}
