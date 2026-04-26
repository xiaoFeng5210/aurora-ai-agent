export class HttpError extends Error {
  status: number
  info: unknown

  constructor(message: string, status: number, info: unknown) {
    super(message)
    this.status = status
    this.info = info
  }
}

export const fetcher = async <T = unknown>(input: RequestInfo | URL, init?: RequestInit): Promise<T> => {
  const headers = new Headers(init?.headers)
  if (!headers.has('Content-Type') && init?.body && !(init.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }

  const res = await fetch(input, {
    credentials: 'include',
    ...init,
    headers,
  })

  if (!res.ok) {
    let info: unknown = null
    try {
      info = await res.json()
    } catch {
      info = await res.text()
    }
    throw new HttpError(`Request failed with status ${res.status}`, res.status, info)
  }

  if (res.status === 204) return undefined as T
  return res.json() as Promise<T>
}
