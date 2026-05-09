import {
  apiDeleteJson,
  apiGet,
  apiPostForm,
  postForm,
  type ApiEnvelope,
} from './client'

export interface NetworkdiskFile {
  category: number
  fs_id: number
  path: string
  server_filename: string
  isdir: number
  local_ctime: number
  local_mtime: number
  server_ctime: number
  server_mtime: number
  size: number
}

export interface NetworkdiskFileListResponse {
  list: NetworkdiskFile[]
  errno: number
}

export interface DeleteKnowledgeFileResponse {
  baidu_networkdisk: unknown
  qdrant: unknown
}

export type RagVectorizeApiVersion = 'v1' | 'v2'

const DEFAULT_RAG_VECTORIZE_API_VERSION: RagVectorizeApiVersion = 'v2'

export const RAG_VECTORIZE_API_VERSION: RagVectorizeApiVersion =
  import.meta.env.VITE_RAG_VECTORIZE_API_VERSION === 'v1'
    ? 'v1'
    : DEFAULT_RAG_VECTORIZE_API_VERSION

export const RAG_VECTORIZE_MODE_LABEL =
  RAG_VECTORIZE_API_VERSION === 'v2' ? 'v2 消息队列' : 'v1 同步向量化'

const RAG_VECTORIZE_ENDPOINTS: Record<RagVectorizeApiVersion, string> = {
  // v1 是旧的同步向量化接口，保留在这里方便回切。
  v1: '/api/v1/rag',
  // v2 只投递向量化任务到消息队列，由消费端异步处理。
  v2: '/api/v2/rag',
}

/** Same branch as service PrecreateUpload: super user uses /apps…/super/, others /apps…/users-data/{username}/ */
export const BAIDU_KNOWLEDGE_SUPER_USERNAME = 'zhangqingfeng'

export function baiduKnowledgeListDir(username: string | null | undefined): string {
  const u = username?.trim()
  if (!u) return 'super/'
  if (u === BAIDU_KNOWLEDGE_SUPER_USERNAME) return 'super/'
  return `users-data/${u}/`
}

export const listKnowledgeFiles = (dir: string) =>
  apiGet<ApiEnvelope<NetworkdiskFileListResponse>>(
    `/file/baidu_networkdisk/file_list?dir=${encodeURIComponent(dir)}`,
  )

export const uploadKnowledgeFileToNetworkdisk = (file: File) => {
  const form = new FormData()
  form.append('filename', file.name)
  form.append('isdir', '0')
  form.append('file', file)
  return apiPostForm<ApiEnvelope<null>>('/file/baidu_networkdisk/upload', form)
}

export const createRagFromFile = (file: File) => {
  const form = new FormData()
  form.append('file', file)
  const endpoint = RAG_VECTORIZE_ENDPOINTS[RAG_VECTORIZE_API_VERSION]
  return postForm<ApiEnvelope<unknown>>(endpoint, form)
}

export const deleteKnowledgeFile = (path: string) =>
  apiDeleteJson<ApiEnvelope<DeleteKnowledgeFileResponse>>('/file/baidu_networkdisk/file', {
    path,
    async: 0,
  })
