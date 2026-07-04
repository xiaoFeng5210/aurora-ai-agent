import { apiDelete, apiGet, apiPost, apiPut, type ApiEnvelope } from './client'

export interface Tag {
  id: number
  user_id: number
  name: string
  created_at: string
  updated_at: string
}

export interface CreateTagRequest {
  name: string
}

export interface UpdateTagRequest {
  name?: string
}

export interface QueryTagRequest {
  name?: string
  page?: number
  page_size?: number
}

export const createTag = (body: CreateTagRequest) => apiPost<ApiEnvelope<Tag>>('/tags', body)

export const getTag = (id: number) => apiGet<ApiEnvelope<Tag>>(`/tags/${id}`)

export const queryTags = (body: QueryTagRequest = {}) =>
  apiPost<ApiEnvelope<Tag[]>>('/tags/query', body)

export const updateTag = (id: number, body: UpdateTagRequest) =>
  apiPut<ApiEnvelope<Tag>>(`/tags/${id}`, body)

export const deleteTag = (id: number) => apiDelete<ApiEnvelope<null>>(`/tags/${id}`)
