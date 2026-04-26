import { forwardRef, type InputHTMLAttributes } from 'react'
import { cn } from '@/lib/cn'

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  invalid?: boolean
}

export const Input = forwardRef<HTMLInputElement, InputProps>(function Input(
  { className, invalid, ...rest },
  ref,
) {
  return (
    <input
      ref={ref}
      className={cn(
        'w-full rounded-md border bg-paper-50 px-3 py-2 text-sm text-ink-900 placeholder:text-ink-300 transition focus:outline-none focus:ring-2 focus:ring-accent-vermilion/30',
        invalid ? 'border-danger/60 focus:ring-danger/30' : 'border-ink-200 focus:border-accent-vermilion/60',
        className,
      )}
      {...rest}
    />
  )
})
