import { apiGet, apiPost, apiPut, type ApiEnvelope } from './client'

export interface RegisterRequest {
  username: string
  password: string
  email: string
  phone?: string
  birthday?: string
  user_prompt?: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface User {
  id: number
  username: string
  email: string
  phone: string
  birthday: string
  user_prompt: string
  created_at: string
  updated_at: string
}

export interface UpdateMeRequest {
  username?: string
  email?: string
  phone?: string
  birthday?: string
  user_prompt?: string
}

export const register = (body: RegisterRequest) =>
  apiPost<ApiEnvelope<null>>('/register', body)

export const login = (body: LoginRequest) =>
  apiPost<ApiEnvelope<null>>('/login', body)

export const logout = () => apiPost<ApiEnvelope<null>>('/logout')

export const getMe = () => apiGet<ApiEnvelope<User>>('/users/me')

export const updateMe = (body: UpdateMeRequest) =>
  apiPut<ApiEnvelope<User>>('/users/me', body)
