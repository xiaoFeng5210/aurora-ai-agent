import { apiDelete, apiGet, apiPost, apiPut, type ApiEnvelope } from './client'

export interface Document {
  id: number
  user_id: number
  display_name: string
  file_name: string | null
  created_at: string
  updated_at: string
}

export interface CreateDocumentRequest {
  display_name: string
  file_name?: string | null
}

export interface UpdateDocumentRequest {
  display_name?: string
  file_name?: string | null
}

export const listDocuments = () => apiGet<ApiEnvelope<Document[]>>('/documents')

export const getDocument = (id: number) =>
  apiGet<ApiEnvelope<Document>>(`/documents/${id}`)

export const createDocument = (body: CreateDocumentRequest) =>
  apiPost<ApiEnvelope<Document>>('/documents', body)

export const updateDocument = (id: number, body: UpdateDocumentRequest) =>
  apiPut<ApiEnvelope<Document>>(`/documents/${id}`, body)

export const deleteDocument = (id: number) =>
  apiDelete<ApiEnvelope<null>>(`/documents/${id}`)
