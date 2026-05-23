import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const userId = ref('')
  const token = ref('')
  const isLoggedIn = ref(false)

  function login(uid: string, t = '') {
    userId.value = uid
    token.value = t
    isLoggedIn.value = true
  }

  function logout() {
    userId.value = ''
    token.value = ''
    isLoggedIn.value = false
  }

  return {
    userId,
    token,
    isLoggedIn,
    login,
    logout,
  }
})
