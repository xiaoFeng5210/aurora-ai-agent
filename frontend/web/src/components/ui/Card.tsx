import type { HTMLAttributes } from 'react'
import { cn } from '@/lib/cn'

export function Card({ className, ...rest }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'rounded-xl border border-ink-200 bg-paper-50 shadow-[0_1px_0_0_rgba(58,46,31,0.04),0_8px_24px_-12px_rgba(58,46,31,0.15)]',
        className,
      )}
      {...rest}
    />
  )
}
