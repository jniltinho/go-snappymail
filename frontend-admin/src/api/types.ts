// Typed mirrors of the /api/v1/admin/* JSON contracts (internal/admin).
// Every response is wrapped in { data?, error? }.

export interface Envelope<T> {
  data?: T
  error?: string
}

export interface LoginResponse {
  token: string
  username: string
  superadmin: boolean
  domains?: string[]
}

export interface Me {
  username: string
  superadmin: boolean
  domains: string[] | null
  role: string
}

export interface Overview {
  accounts: number
  domains: number
  aliases: number
  admins: number
  version: string | null
  servers: number | null
  queue: number | null
  active_sessions: number | null
}

export interface Domain {
  domain: string
  description: string
  aliases: number
  mailboxes: number
  maxquota: number
  transport: string
  backupmx: boolean
  active: boolean
  created: string
  modified: string
}

export interface Mailbox {
  username: string
  name: string
  maildir: string
  quota: number
  local_part: string
  domain: string
  active: boolean
  created: string
  modified: string
}

export interface Alias {
  address: string
  goto: string
  domain: string
  active: boolean
  created: string
  modified: string
}

export interface Admin {
  username: string
  superadmin: boolean
  active: boolean
  domains?: string[]
}
