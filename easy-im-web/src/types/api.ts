// Wire DTOs — must match go-im JSON tags after Task 1.

// user-api
export interface LoginReq {
  phone: string
  password: string
}
export interface LoginResp {
  token: string
  expire: number
}
export interface RegisterReq {
  phone: string
  password: string
  nickname: string
  sex: number
  avatar: string
}
export interface UserInfoDTO {
  id: string
  mobile: string
  nickname: string
  sex: number
  avatar: string
  email: string
}
export interface UserInfoResp {
  user: UserInfoDTO
}
export interface UpdateProfileReq {
  nickname: string
  sex: number
  email: string
  avatar: string
}
export interface UpdateProfileResp {
  success: boolean
}
export interface UpdatePasswordReq {
  old_password: string
  new_password: string
}
export interface UpdatePasswordResp {
  success: boolean
}

// social-api
export interface FriendDTO {
  id?: string
  friend_uid?: string
  nickname?: string
  avatar?: string
  remark?: string
}
export interface FriendRequestDTO {
  id?: string
  user_id?: string
  req_uid?: string
  req_msg?: string
  req_time?: number
  handle_result?: number
  handle_msg?: string
}
export interface FriendListResp {
  list: FriendDTO[]
}
export interface FriendPutInListResp {
  list: FriendRequestDTO[]
}
export interface FriendsOnlineResp {
  onlineList: Record<string, boolean>
}
export interface FriendPutInReq {
  phone: string
  req_msg: string
  req_time: number
}
export interface FriendPutInHandleReq {
  friend_req_id: string
  handle_result: number
}

// im-api
export interface ChatLogDTO {
  id?: string
  conversationId?: string
  sendId?: string
  recvId?: string
  msgType?: number
  msgContent?: string
  chatType?: number
  sendTime?: number
}
export interface ChatLogListResp {
  list: ChatLogDTO[]
}
export interface ConversationDTO {
  conversationId?: string
  chatType?: number
  isShow?: boolean
  seq?: number
  read?: number
}
export interface GetConversationsResp {
  conversationList: Record<string, ConversationDTO>
}
export interface SetUpUserConversationReq {
  sendId: string
  recvId: string
  chatType: number
}
export interface PutConversationsReq {
  conversationList: Record<string, ConversationDTO>
}
export interface ChatLogReq {
  msgId?: string
  conversationId: string
  startSendTime?: number
  endSendTime?: number
  count?: number
}
