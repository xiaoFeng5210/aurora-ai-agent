import { forwardRef, type ButtonHTMLAttributes } from 'react'
import { cn } from '@/lib/cn'

export type ButtonVariant = 'primary' | 'ghost' | 'danger' | 'subtle'
export type ButtonSize = 'sm' | 'md' | 'lg'

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
  size?: ButtonSize
  loading?: boolean
}

const variantClass: Record<ButtonVariant, string> = {
  primary:
    'bg-accent-vermilion text-paper-50 hover:bg-accent-vermilion/90 active:bg-accent-vermilion/80 shadow-sm',
  ghost:
    'border border-ink-300 bg-transparent text-ink-900 hover:bg-paper-200/60 active:bg-paper-200',
  danger:
    'bg-danger text-paper-50 hover:bg-danger/90 active:bg-danger/80 shadow-sm',
  subtle:
    'bg-paper-200/60 text-ink-900 hover:bg-paper-200 active:bg-paper-300/60',
}

const sizeClass: Record<ButtonSize, string> = {
  sm: 'px-3 py-1.5 text-sm rounded-md',
  md: 'px-4 py-2 text-sm rounded-md',
  lg: 'px-6 py-3 text-base rounded-lg',
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  { variant = 'primary', size = 'md', loading, className, children, disabled, ...rest },
  ref,
) {
  return (
    <button
      ref={ref}
      disabled={disabled || loading}
      className={cn(
        'inline-flex items-center justify-center gap-2 font-medium transition disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus-visible:ring-2 focus-visible:ring-accent-vermilion/40',
        variantClass[variant],
        sizeClass[size],
        className,
      )}
      {...rest}
    >
      {loading ? <span className="h-4 w-4 animate-spin rounded-full border-2 border-paper-50 border-t-transparent" /> : null}
      {children}
    </button>
  )
})
