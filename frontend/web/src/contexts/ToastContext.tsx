import { createContext, useCallback, useState, type ReactNode } from 'react'

export type ToastVariant = 'info' | 'success' | 'error'

export interface ToastItem {
  id: number
  message: string
  variant: ToastVariant
}

export interface ToastContextValue {
  toasts: ToastItem[]
  show: (message: string, variant?: ToastVariant) => void
  dismiss: (id: number) => void
}

export const ToastContext = createContext<ToastContextValue | null>(null)

let nextId = 1

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastItem[]>([])

  const dismiss = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id))
  }, [])

  const show = useCallback(
    (message: string, variant: ToastVariant = 'info') => {
      const id = nextId++
      setToasts((prev) => [...prev, { id, message, variant }])
      window.setTimeout(() => dismiss(id), 3500)
    },
    [dismiss],
  )

  return (
    <ToastContext.Provider value={{ toasts, show, dismiss }}>
      {children}
      <ToastViewport toasts={toasts} onDismiss={dismiss} />
    </ToastContext.Provider>
  )
}

function ToastViewport({
  toasts,
  onDismiss,
}: {
  toasts: ToastItem[]
  onDismiss: (id: number) => void
}) {
  return (
    <div className="pointer-events-none fixed top-4 right-4 z-50 flex flex-col gap-2">
      {toasts.map((t) => (
        <button
          key={t.id}
          type="button"
          onClick={() => onDismiss(t.id)}
          className={`pointer-events-auto rounded-md border px-4 py-2 text-sm shadow-md transition ${
            t.variant === 'error'
              ? 'border-danger/40 bg-paper-50 text-danger'
              : t.variant === 'success'
                ? 'border-accent-olive/40 bg-paper-50 text-accent-olive'
                : 'border-ink-200 bg-paper-50 text-ink-800'
          }`}
        >
          {t.message}
        </button>
      ))}
    </div>
  )
}
