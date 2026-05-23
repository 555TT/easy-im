import { socialHttp } from './http'
import type {
  FriendListResp,
  FriendPutInHandleReq,
  FriendPutInListResp,
  FriendPutInReq,
  FriendsOnlineResp,
} from '@/types/api'

export function listFriends(): Promise<FriendListResp> {
  return socialHttp.get('/friends')
}

export function listOnlineFriends(): Promise<FriendsOnlineResp> {
  return socialHttp.get('/friends/online')
}

export function listFriendRequests(): Promise<FriendPutInListResp> {
  return socialHttp.get('/friend/putIns')
}

export function sendFriendRequest(body: FriendPutInReq): Promise<unknown> {
  return socialHttp.post('/friend/putIn', body)
}

export function handleFriendRequest(body: FriendPutInHandleReq): Promise<unknown> {
  return socialHttp.put('/friend/putIn', body)
}
