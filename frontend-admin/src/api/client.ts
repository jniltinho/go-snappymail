import axios, { AxiosError } from 'axios'
import type { Envelope } from './types'

declare const API_BASE: string

const TOKEN_KEY = 'admin.token'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

const http = axios.create({ baseURL: API_BASE })

// Attach the bearer token to every request.
http.interceptors.request.use((cfg) => {
  const token = getToken()
  if (token) {
    cfg.headers = cfg.headers ?? {}
    cfg.headers.Authorization = `Bearer ${token}`
  }
  return cfg
})

// A rejected admin token means the session is over — drop it so the router
// guard can bounce the user back to the login screen.
http.interceptors.response.use(
  (res) => res,
  (err: AxiosError) => {
    if (err.response?.status === 401) clearToken()
    return Promise.reject(err)
  },
)

/** apiError extracts the server's { error } message, falling back to the HTTP text. */
export function apiError(err: unknown): string {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as Envelope<unknown> | undefined
    if (data?.error) return data.error
    return err.message
  }
  return String(err)
}

/** unwrap returns response.data.data, the payload inside the envelope. */
async function unwrap<T>(p: Promise<{ data: Envelope<T> }>): Promise<T> {
  const res = await p
  return res.data.data as T
}

export const api = {
  get: <T>(url: string) => unwrap<T>(http.get(url)),
  post: <T>(url: string, body?: unknown) => unwrap<T>(http.post(url, body)),
  put: <T>(url: string, body?: unknown) => unwrap<T>(http.put(url, body)),
  del: <T>(url: string) => unwrap<T>(http.delete(url)),
}
