package xerr

// User err
var (
	PhoneNotFound          = New(ServerCommonError, "api not found")
	PhoneAlreadyRegistered = New(ServerCommonError, "手机号已注册")
	IdNotFound             = New(ServerCommonError, "api id not found")
	UserPwdErr             = New(ServerCommonError, "password is wrong")
	ParamError             = New(RequestParamError, "params error")
)

// Friend Err
var (
	FriendAlreadyExists      = New(ServerCommonError, "friend already exists")
	FriendRequestOnPending   = New(ServerCommonError, "friend request on pending")
	FriendRequestRefused     = New(ServerCommonError, "friend request refused")
	FriendListNotFound       = New(ServerCommonError, "friend list not found")
	FriendReqListNotFound    = New(ServerCommonError, "friend request list not found")
	FriendPhoneNotRegistered = New(ServerCommonError, "该手机号未注册")
	CannotAddSelf            = New(ServerCommonError, "不能添加自己为好友")

	FindFriendByIdErr = New(ServerCommonError, "find friend by id error")
)

// Group Err
var (
	GroupNotFound        = New(ServerCommonError, "group not found ")
	GroupPutInNotFound   = New(ServerCommonError, "group put in request not found")
	GroupInviterNotFound = New(ServerCommonError, "group inviter not found")

	FindGroupByIdErr = New(ServerCommonError, "find group by id error, user haven't attend in any group")
)
