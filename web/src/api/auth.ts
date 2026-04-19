import { apiGet, apiPost } from './client'

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

export const register = (body: RegisterRequest) => apiPost<void>('/register', body)

export const login = (body: LoginRequest) => apiPost<{code: number, message: string}>('/login', body)

export const getMe = () => apiGet<User>('/users/me')
