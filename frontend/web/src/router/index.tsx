import { createHashRouter, Navigate } from 'react-router-dom'
import { Home } from '@/pages/home/Home'
import { Login } from '@/pages/auth/Login'
import { Register } from '@/pages/auth/Register'
import { Cards } from '@/pages/cards/Cards'
import { Chat } from '@/pages/chat/Chat'
import { Profile } from '@/pages/profile/Profile'
import { Knowledge } from '@/pages/knowledge/Knowledge'
import { AuthGuard } from './AuthGuard'
import { GuestGuard } from './GuestGuard'

export const router = createHashRouter([
  { path: '/', element: <Home /> },
  {
    path: '/login',
    element: (
      <GuestGuard>
        <Login />
      </GuestGuard>
    ),
  },
  { path: '/register', element: <Register /> },
  {
    path: '/cards',
    element: (
      <AuthGuard>
        <Cards />
      </AuthGuard>
    ),
  },
  {
    path: '/chat',
    element: (
      <AuthGuard>
        <Chat />
      </AuthGuard>
    ),
  },
  {
    path: '/chat/:documentId',
    element: (
      <AuthGuard>
        <Chat />
      </AuthGuard>
    ),
  },
  {
    path: '/profile',
    element: (
      <AuthGuard>
        <Profile />
      </AuthGuard>
    ),
  },
  {
    path: '/knowledge',
    element: (
      <AuthGuard>
        <Knowledge />
      </AuthGuard>
    ),
  },
  { path: '*', element: <Navigate to="/" replace /> },
])
