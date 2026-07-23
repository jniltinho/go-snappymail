export interface Folder {
  name: string
  label: string
  iconType: string
  unseen: number
  messages: number
  depth: number
}

export interface MailMessage {
  uid: number
  subject: string
  from: string
  fromEmail: string
  date: string
  seen: boolean
  flagged: boolean
  size: number
  to?: string
  htmlBody?: string
  plainBody?: string
  attachments?: AttachmentMeta[]
}

export interface AttachmentMeta {
  filename: string
  part: number
  sizeLabel: string
  contentType: string
}
