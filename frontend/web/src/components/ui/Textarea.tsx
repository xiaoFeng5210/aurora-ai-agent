import { forwardRef, type TextareaHTMLAttributes } from 'react'
import { cn } from '@/lib/cn'

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaHTMLAttributes<HTMLTextAreaElement>>(
  function Textarea({ className, ...rest }, ref) {
    return (
      <textarea
        ref={ref}
        className={cn(
          'w-full resize-none rounded-md border border-ink-200 bg-paper-50 px-3 py-2 text-sm text-ink-900 placeholder:text-ink-300 transition focus:border-accent-vermilion/60 focus:outline-none focus:ring-2 focus:ring-accent-vermilion/30',
          className,
        )}
        {...rest}
      />
    )
  },
)
