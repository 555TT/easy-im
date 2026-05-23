import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { decodeJwt, isExpired } from '@/utils/jwt'

export const useAuthStore = defineStore(
  'auth',
  () => {
    const token = ref<string>('')
    const expire = ref<number>(0)
    const userId = ref<string>('')

    const isLoggedIn = computed(() => {
      if (!token.value) return false
      const p = decodeJwt(token.value)
      return !isExpired(p)
    })

    function setSession(t: string, exp: number): void {
      token.value = t
      expire.value = exp
      const p = decodeJwt(t)
      userId.value = p?.peninsula ?? ''
    }

    function clear(): void {
      token.value = ''
      expire.value = 0
      userId.value = ''
    }

    return { token, expire, userId, isLoggedIn, setSession, clear }
  },
  { persist: { key: 'easy-im.auth' } },
)
