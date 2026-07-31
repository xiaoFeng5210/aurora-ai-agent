const KEY = 'ariadne.rememberedCredentials'

export interface RememberedCredentials {
  email: string
  password: string
}

function encode(text: string): string {
  const bytes = new TextEncoder().encode(text)
  let binary = ''
  for (const b of bytes) binary += String.fromCharCode(b)
  return btoa(binary)
}

function decode(encoded: string): string {
  const binary = atob(encoded)
  const bytes = Uint8Array.from(binary, (c) => c.charCodeAt(0))
  return new TextDecoder().decode(bytes)
}

export function loadCredentials(): RememberedCredentials | null {
  try {
    const raw = localStorage.getItem(KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw) as { email?: unknown; password?: unknown }
    if (typeof parsed.email !== 'string' || typeof parsed.password !== 'string') return null
    return { email: decode(parsed.email), password: decode(parsed.password) }
  } catch {
    return null
  }
}

export function saveCredentials(credentials: RememberedCredentials): void {
  try {
    localStorage.setItem(
      KEY,
      JSON.stringify({ email: encode(credentials.email), password: encode(credentials.password) }),
    )
  } catch {
    // localStorage 不可用时静默忽略
  }
}

export function clearCredentials(): void {
  try {
    localStorage.removeItem(KEY)
  } catch {
    // 忽略
  }
}
