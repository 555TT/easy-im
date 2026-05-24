import test from 'node:test'
import assert from 'node:assert/strict'

import { buildOnlineMap } from './online.ts'

test('buildOnlineMap marks only connected friends as online', () => {
  const onlineMap = buildOnlineMap(['u1', 'u2', 'u3'], ['u2', 'u4'])

  assert.deepEqual(onlineMap, {
    u1: false,
    u2: true,
    u3: false,
  })
})
