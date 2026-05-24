import type { FriendDTO } from '../types/api'
import type { Friend } from '../types/domain'

const UNKNOWN_FRIEND_NAME = '未知用户'

export function mapFriendDTO(friend: FriendDTO, onlineMap: Record<string, boolean>): Friend {
  const userId = friend.friend_uid ?? ''
  const nickname = friend.nickname?.trim() || UNKNOWN_FRIEND_NAME

  return {
    id: friend.id || userId,
    userId,
    nickname,
    avatar: friend.avatar ?? '',
    remark: friend.remark?.trim() ?? '',
    online: !!onlineMap[userId],
  }
}
