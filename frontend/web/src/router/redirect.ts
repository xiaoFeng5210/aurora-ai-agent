const AUTH_PATHS = new Set(['/login', '/register'])

export function safeRedirectPath(raw: string | null | undefined, fallback = '/'): string {
  if (!raw) return fallback

  let path = raw
  try {
    path = decodeURIComponent(raw)
  } catch {
    return fallback
  }

  if (!path.startsWith('/') || path.startsWith('//')) return fallback

  const pathname = path.split('?')[0] ?? path
  if (AUTH_PATHS.has(pathname)) return fallback

  return path
}
