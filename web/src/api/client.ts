import { fetcher, HttpError } from '@/lib/fetcher'

export const API_BASE = '/api/v1'

const url = (path: string) => `${API_BASE}${path.startsWith('/') ? path : `/${path}`}`

export interface ApiEnvelope<T> {
  code: number
  message: string
  data?: T
}

const isEnvelope = (v: unknown): v is ApiEnvelope<unknown> =>
  !!v && typeof v === 'object' && 'code' in v

const unwrap = async <T>(p: Promise<unknown>): Promise<T> => {
  const raw = await p

  if (isEnvelope(raw)) {
    if (raw.code !== 0) {
      throw new HttpError(raw.message || `Business error (${raw.code})`, 200, raw)
    }
    return raw as T
  }
  return raw as T
}

export const apiGet = <T = unknown>(path: string, init?: RequestInit) =>
  unwrap<T>(fetcher<unknown>(url(path), { method: 'GET', ...init }))

export const apiPost = <T = unknown>(path: string, body?: unknown, init?: RequestInit) =>
  unwrap<T>(
    fetcher<unknown>(url(path), {
      method: 'POST',
      body: body === undefined ? undefined : JSON.stringify(body),
      ...init,
    }),
  )

export const apiPut = <T = unknown>(path: string, body?: unknown, init?: RequestInit) =>
  unwrap<T>(
    fetcher<unknown>(url(path), {
      method: 'PUT',
      body: body === undefined ? undefined : JSON.stringify(body),
      ...init,
    }),
  )

export const apiDelete = <T = unknown>(path: string, init?: RequestInit) =>
  unwrap<T>(fetcher<unknown>(url(path), { method: 'DELETE', ...init }))

export const swrFetcher = <T = unknown>(path: string) => apiGet<T>(path)

export { HttpError }
