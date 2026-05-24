import test from 'node:test'
import assert from 'node:assert/strict'

import { mapFriendDTO } from './contact.ts'

test('mapFriendDTO keeps backend nickname', () => {
  const friend = mapFriendDTO(
    {
      id: 'f1',
      friend_uid: 'u2',
      nickname: '小明',
      avatar: 'a.png',
      remark: '',
    },
    {},
  )

  assert.equal(friend.nickname, '小明')
  assert.equal(friend.userId, 'u2')
})

test('mapFriendDTO does not fall back to user id when nickname is missing', () => {
  const friend = mapFriendDTO(
    {
      id: 'f1',
      friend_uid: 'u2',
      nickname: '',
      avatar: '',
      remark: '',
    },
    {},
  )

  assert.equal(friend.nickname, '未知用户')
  assert.notEqual(friend.nickname, friend.userId)
})
