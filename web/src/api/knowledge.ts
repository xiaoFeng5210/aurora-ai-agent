import {
  apiDeleteJson,
  apiGet,
  apiPostForm,
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
  return apiPostForm<ApiEnvelope<unknown>>('/rag', form)
}

export const deleteKnowledgeFile = (path: string) =>
  apiDeleteJson<ApiEnvelope<DeleteKnowledgeFileResponse>>('/file/baidu_networkdisk/file', {
    path,
    async: 0,
  })
