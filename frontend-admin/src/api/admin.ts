// Typed endpoint functions for the admin API. Views never touch axios directly —
// they call these, keeping the JSON contract in one place.
import { api } from './client'
import type { Admin, Alias, Domain, LoginResponse, Mailbox, Me, Overview } from './types'

export const AdminAPI = {
  login: (username: string, password: string) =>
    api.post<LoginResponse>('/auth/login', { username, password }),
  me: () => api.get<Me>('/me'),
  overview: () => api.get<Overview>('/overview'),

  // Domains
  listDomains: () => api.get<Domain[]>('/domains'),
  createDomain: (body: Partial<Domain>) => api.post<Domain>('/domains', body),
  updateDomain: (domain: string, body: Partial<Domain>) =>
    api.put<Domain>(`/domains/${encodeURIComponent(domain)}`, body),
  deleteDomain: (domain: string) => api.del<unknown>(`/domains/${encodeURIComponent(domain)}`),

  // Mailboxes (accounts)
  listMailboxes: () => api.get<Mailbox[]>('/mailboxes'),
  createMailbox: (body: Record<string, unknown>) => api.post<Mailbox>('/mailboxes', body),
  updateMailbox: (username: string, body: Record<string, unknown>) =>
    api.put<Mailbox>(`/mailboxes/${encodeURIComponent(username)}`, body),
  deleteMailbox: (username: string) =>
    api.del<unknown>(`/mailboxes/${encodeURIComponent(username)}`),

  // Aliases
  listAliases: () => api.get<Alias[]>('/aliases'),
  createAlias: (body: Record<string, unknown>) => api.post<Alias>('/aliases', body),
  updateAlias: (address: string, body: Record<string, unknown>) =>
    api.put<Alias>(`/aliases/${encodeURIComponent(address)}`, body),
  deleteAlias: (address: string) => api.del<unknown>(`/aliases/${encodeURIComponent(address)}`),

  // Admins
  listAdmins: () => api.get<Admin[]>('/admins'),
  createAdmin: (body: Record<string, unknown>) => api.post<Admin>('/admins', body),
  updateAdmin: (username: string, body: Record<string, unknown>) =>
    api.put<Admin>(`/admins/${encodeURIComponent(username)}`, body),
  deleteAdmin: (username: string) => api.del<unknown>(`/admins/${encodeURIComponent(username)}`),
}
