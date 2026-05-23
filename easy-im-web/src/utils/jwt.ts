// Minimal JWT payload decoder. No signature verification — that is the server's job.

export interface JwtPayload {
  exp: number
  iat: number
  peninsula: string // user id (project-specific claim key)
}

function base64UrlDecode(input: string): string {
  const pad = input.length % 4 === 0 ? '' : '='.repeat(4 - (input.length % 4))
  const b64 = (input + pad).replace(/-/g, '+').replace(/_/g, '/')
  return decodeURIComponent(
    atob(b64)
      .split('')
      .map((c) => '%' + c.charCodeAt(0).toString(16).padStart(2, '0'))
      .join(''),
  )
}

export function decodeJwt(token: string): JwtPayload | null {
  const parts = token.split('.')
  if (parts.length !== 3) return null
  try {
    return JSON.parse(base64UrlDecode(parts[1])) as JwtPayload
  } catch {
    return null
  }
}

export function isExpired(payload: JwtPayload | null, nowSec = Math.floor(Date.now() / 1000)): boolean {
  if (!payload) return true
  return payload.exp <= nowSec
}
