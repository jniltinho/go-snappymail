import axios from 'axios'

axios.defaults.withCredentials = true
axios.defaults.xsrfCookieName = 'csrf_token'
axios.defaults.xsrfHeaderName = 'X-CSRF-Token'

axios.interceptors.request.use((cfg) => {
  const match = document.cookie.match(/(?:^|;\s*)csrf_token=([^;]+)/)
  if (match) {
    cfg.headers['X-CSRF-Token'] = decodeURIComponent(match[1])
  }
  return cfg
})

/** Prime CSRF cookie before mutating requests. */
export async function primeCsrf(): Promise<void> {
  await axios.get('/', { baseURL: '' })
}

export { axios }
