import { createHashRouter, Navigate } from 'react-router-dom'
import { Home } from '@/pages/home/Home'
import { Login } from '@/pages/auth/Login'
import { Register } from '@/pages/auth/Register'
import { Chat } from '@/pages/chat/Chat'
import { Profile } from '@/pages/profile/Profile'
import { AuthGuard } from './AuthGuard'

export const router = createHashRouter([
  { path: '/', element: <Home /> },
  { path: '/login', element: <Login /> },
  { path: '/register', element: <Register /> },
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
  { path: '*', element: <Navigate to="/" replace /> },
])
