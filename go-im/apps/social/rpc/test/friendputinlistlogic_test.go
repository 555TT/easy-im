package test

import (
	"context"
	"github.com/peninsula12/easy-im/go-im/apps/social/rpc/internal/logic"
	"github.com/peninsula12/easy-im/go-im/apps/social/rpc/social"
	"testing"
)

func TestFriendPutInListLogic_FriendPutInList(t *testing.T) {
	t.Run("returns empty list when user has no friend requests", func(t *testing.T) {
		l := logic.NewFriendPutInListLogic(context.Background(), svcCtx)
		got, err := l.FriendPutInList(&social.FriendPutInListReq{UserId: "user-with-no-friend-requests"})
		if err != nil {
			t.Fatalf("FriendPutInList() error = %v, want nil", err)
		}
		if got == nil {
			t.Fatal("FriendPutInList() got nil response")
		}
		if len(got.List) != 0 {
			t.Fatalf("FriendPutInList() list length = %d, want 0", len(got.List))
		}
	})
}
