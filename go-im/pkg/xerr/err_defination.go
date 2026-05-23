package xerr

// User err
var (
	PhoneNotFound          = New(ServerCommonError, "api not found")
	PhoneAlreadyRegistered = New(ServerCommonError, "手机号已注册")
	IdNotFound             = New(ServerCommonError, "api id not found")
	UserPwdErr             = New(ServerCommonError, "password is wrong")
	ParamError             = New(RequestParamError, "params error")
	InvalidEmail           = New(ServerCommonError, "邮箱格式不正确")
	EmptyNickname          = New(ServerCommonError, "用户名不能为空")
	InvalidSex             = New(ServerCommonError, "性别参数不合法")
	PasswordRequired       = New(ServerCommonError, "old_password and new_password must be provided together")
	PasswordUnchanged      = New(ServerCommonError, "new password must be different from old password")
	PasswordTooShort       = New(ServerCommonError, "new password must be at least 6 characters")
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
