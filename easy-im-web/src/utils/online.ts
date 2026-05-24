export function buildOnlineMap(friendIds: string[], onlineUserIds: string[]): Record<string, boolean> {
  const onlineSet = new Set(onlineUserIds)
  const onlineMap: Record<string, boolean> = {}

  for (const friendId of friendIds) {
    onlineMap[friendId] = onlineSet.has(friendId)
  }

  return onlineMap
}
