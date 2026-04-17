import { createHashRouter, Navigate } from 'react-router-dom'
import { Home } from '@/pages/home/Home'
import { Login } from '@/pages/auth/Login'
import { Register } from '@/pages/auth/Register'
import { Chat } from '@/pages/chat/Chat'
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
  { path: '*', element: <Navigate to="/" replace /> },
])
