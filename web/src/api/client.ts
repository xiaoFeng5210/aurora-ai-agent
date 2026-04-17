import { fetcher, HttpError } from '@/lib/fetcher'

export const API_BASE = '/api/v1'

const url = (path: string) => `${API_BASE}${path.startsWith('/') ? path : `/${path}`}`

export const apiGet = <T = unknown>(path: string, init?: RequestInit) =>
  fetcher<T>(url(path), { method: 'GET', ...init })

export const apiPost = <T = unknown>(path: string, body?: unknown, init?: RequestInit) =>
  fetcher<T>(url(path), {
    method: 'POST',
    body: body === undefined ? undefined : JSON.stringify(body),
    ...init,
  })

export const apiPut = <T = unknown>(path: string, body?: unknown, init?: RequestInit) =>
  fetcher<T>(url(path), {
    method: 'PUT',
    body: body === undefined ? undefined : JSON.stringify(body),
    ...init,
  })

export const apiDelete = <T = unknown>(path: string, init?: RequestInit) =>
  fetcher<T>(url(path), { method: 'DELETE', ...init })

export const swrFetcher = <T = unknown>(path: string) => apiGet<T>(path)

export { HttpError }
