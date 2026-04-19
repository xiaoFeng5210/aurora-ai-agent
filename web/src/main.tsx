import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import './index.css'
import { AuthProvider } from '@/contexts/AuthContext'
import { ToastProvider } from '@/contexts/ToastContext'
import { router } from '@/router'

createRoot(document.getElementById('root')!).render(
    <AuthProvider>
      <ToastProvider>
      <RouterProvider router={router} />
    </ToastProvider>
  </AuthProvider>,
)
