import { apiDelete, apiGet, apiPost, apiPut, type ApiEnvelope } from './client'
import type { Tag } from './tag'

export interface Card {
  id: number
  user_id: number
  content: string
  tags: string[]
  tag_ids?: number[]
  external_links: string[]
  internal_links: string[]
  created_at: string
  updated_at: string
}

export interface CreateCardRequest {
  content: string
  tags?: string[]
  tag_ids?: number[]
  external_links?: string[]
  internal_links?: string[]
}

export interface UpdateCardRequest {
  content?: string
  tags?: string[]
  tag_ids?: number[]
  external_links?: string[]
  internal_links?: string[]
}

export interface QueryCardRequest {
  content?: string
  tags?: string[]
  tag_ids?: number[]
  page?: number
  page_size?: number
}

export const createCard = (body: CreateCardRequest) => apiPost<ApiEnvelope<Card>>('/cards', body)

export const getCard = (id: number) => apiGet<ApiEnvelope<Card>>(`/cards/${id}`)

export const queryCards = (body: QueryCardRequest = {}) =>
  apiPost<ApiEnvelope<Card[]>>('/cards/query', body)

export const updateCard = (id: number, body: UpdateCardRequest) =>
  apiPut<ApiEnvelope<Card>>(`/cards/${id}`, body)

export const deleteCard = (id: number) => apiDelete<ApiEnvelope<null>>(`/cards/${id}`)

export interface CardTag {
  id: number
  user_id: number
  card_id: number
  tag_id: number
  created_at: string
  updated_at: string
}

export const listTagsByCard = (cardId: number) =>
  apiGet<ApiEnvelope<Tag[]>>(`/cards/${cardId}/tags`)

export const addTagToCard = (cardId: number, tagId: number) =>
  apiPost<ApiEnvelope<CardTag>>(`/cards/${cardId}/tags`, { tag_id: tagId })

export const replaceCardTags = (cardId: number, tagIds: number[]) =>
  apiPut<ApiEnvelope<CardTag[]>>(`/cards/${cardId}/tags`, { tag_ids: tagIds })

export const deleteTagFromCard = (cardId: number, tagId: number) =>
  apiDelete<ApiEnvelope<null>>(`/cards/${cardId}/tags/${tagId}`)
