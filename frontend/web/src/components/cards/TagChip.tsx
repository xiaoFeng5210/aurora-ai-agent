import type { ButtonHTMLAttributes } from 'react'
import { cn } from '@/lib/cn'

export interface TagChipProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  active?: boolean
}

export function TagChip({ active, className, children, ...rest }: TagChipProps) {
  return (
    <button
      type="button"
      className={cn(
        'inline-flex shrink-0 cursor-pointer items-center gap-1 whitespace-nowrap rounded-full border px-3 py-1.5 text-xs font-medium transition',
        active
          ? 'border-accent-vermilion bg-accent-vermilion/10 text-accent-vermilion'
          : 'border-ink-200 bg-paper-50 text-ink-600 hover:border-ink-300 hover:text-ink-900',
        className,
      )}
      {...rest}
    >
      {children}
    </button>
  )
}
