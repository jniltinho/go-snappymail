import { defineStore } from 'pinia'
import { AdminAPI } from '../api/admin'
import { clearToken, getToken, setToken } from '../api/client'

interface AuthState {
  token: string | null
  username: string
  superadmin: boolean
  domains: string[]
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: getToken(),
    username: '',
    superadmin: false,
    domains: [],
  }),
  getters: {
    isAuthenticated: (s) => !!s.token,
  },
  actions: {
    async login(username: string, password: string) {
      const res = await AdminAPI.login(username, password)
      setToken(res.token)
      this.token = res.token
      this.username = res.username
      this.superadmin = res.superadmin
      this.domains = res.domains ?? []
    },
    // refresh re-hydrates identity from the token on a hard reload.
    async refresh() {
      if (!this.token) return
      const me = await AdminAPI.me()
      this.username = me.username
      this.superadmin = me.superadmin
      this.domains = me.domains ?? []
    },
    logout() {
      clearToken()
      this.token = null
      this.username = ''
      this.superadmin = false
      this.domains = []
    },
  },
})
