import { apiDelete, apiGet, apiPost, apiPut } from './client'

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

export const listDocuments = () => apiGet<Document[]>('/documents')

export const getDocument = (id: number) => apiGet<Document>(`/documents/${id}`)

export const createDocument = (body: CreateDocumentRequest) =>
  apiPost<Document>('/documents', body)

export const updateDocument = (id: number, body: UpdateDocumentRequest) =>
  apiPut<Document>(`/documents/${id}`, body)

export const deleteDocument = (id: number) =>
  apiDelete<{ code: number; message: string }>(`/documents/${id}`)
