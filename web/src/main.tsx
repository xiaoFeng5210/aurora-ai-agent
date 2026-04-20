import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import "@radix-ui/themes/styles.css";
import './index.css'
import { AuthProvider } from '@/contexts/AuthContext'
import { ToastProvider } from '@/contexts/ToastContext'
import { router } from '@/router'
import { Theme } from '@radix-ui/themes'

createRoot(document.getElementById('root')!).render(
  <Theme
    appearance="light"
    accentColor="tomato"
    grayColor="sand"
    radius="medium"
    scaling="100%"
    panelBackground="translucent"
  >
    <AuthProvider>
      <ToastProvider>
        <RouterProvider router={router} />
      </ToastProvider>
    </AuthProvider>
  </Theme>,
)
