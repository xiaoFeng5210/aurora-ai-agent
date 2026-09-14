import { apiGet, type ApiEnvelope } from './client'

export const getMyPoints = () =>
  apiGet<ApiEnvelope<{ balance: number }>>('/users/me/points')
