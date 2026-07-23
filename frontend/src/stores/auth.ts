import { defineStore } from 'pinia'
import { ref } from 'vue'
import { axios, primeCsrf } from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const isAuthenticated = ref(false)
  const username = ref('')
  const datetimeFormat = ref('02/01/2006 15:04')
  const appVersion = ref('')
  const quotaUsed = ref(0)
  const quotaLimit = ref(0)

  const loginUser = ref('')
  const loginPwd = ref('')
  const loginBusy = ref(false)
  const loginErr = ref<string | null>(null)

  async function fetchQuota(): Promise<void> {
    try {
      const res = await axios.get(`${API_BASE}/auth/quota`)
      quotaUsed.value = res.data.used || 0
      quotaLimit.value = res.data.limit || 0
    } catch {
      /* quota optional */
    }
  }

  async function checkSession(): Promise<void> {
    await primeCsrf()
    try {
      const [meRes, versionRes] = await Promise.allSettled([
        axios.get(`${API_BASE}/auth/me`),
        axios.get(`${API_BASE}/version`),
      ])
      if (versionRes.status === 'fulfilled') {
        appVersion.value = versionRes.value.data.version || ''
      }
      if (meRes.status === 'fulfilled') {
        username.value = meRes.value.data.username || ''
        if (meRes.value.data.datetime_format) {
          datetimeFormat.value = meRes.value.data.datetime_format
        }
        isAuthenticated.value = true
        await fetchQuota()
      }
    } catch {
      isAuthenticated.value = false
    }
  }

  async function login(): Promise<boolean> {
    loginErr.value = null
    const email = loginUser.value.trim()
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      loginErr.value = 'Please enter a valid email address.'
      return false
    }
    if (loginPwd.value.length < 1) {
      loginErr.value = 'Please enter your password.'
      return false
    }

    loginBusy.value = true
    try {
      await primeCsrf()
      const body = new URLSearchParams()
      body.set('username', email)
      body.set('password', loginPwd.value)
      await axios.post(`${API_BASE}/auth/login`, body, {
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
      })
      username.value = email
      isAuthenticated.value = true
      const me = await axios.get(`${API_BASE}/auth/me`)
      if (me.data.datetime_format) datetimeFormat.value = me.data.datetime_format
      await fetchQuota()
      return true
    } catch {
      loginErr.value = 'Invalid credentials or server unreachable.'
      return false
    } finally {
      loginBusy.value = false
    }
  }

  async function logout(): Promise<void> {
    try {
      await axios.post(`${API_BASE}/auth/logout`)
    } catch {
      /* ignore */
    }
    isAuthenticated.value = false
    username.value = ''
  }

  return {
    isAuthenticated,
    username,
    datetimeFormat,
    appVersion,
    quotaUsed,
    quotaLimit,
    loginUser,
    loginPwd,
    loginBusy,
    loginErr,
    checkSession,
    login,
    logout,
    fetchQuota,
  }
})
