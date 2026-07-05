import type { ButtonHTMLAttributes } from 'react'
import { Check } from 'lucide-react'
import { cn } from '@/lib/cn'

export interface TagChipProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  active?: boolean
  showCheckbox?: boolean
}

export function TagChip({ active, showCheckbox, className, children, ...rest }: TagChipProps) {
  return (
    <button
      type="button"
      className={cn(
        'inline-flex shrink-0 cursor-pointer items-center whitespace-nowrap rounded-full border px-3 py-1.5 text-xs font-medium transition',
        showCheckbox ? 'gap-2' : 'gap-1',
        active
          ? 'border-accent-vermilion bg-accent-vermilion/10 text-accent-vermilion'
          : 'border-ink-200 bg-paper-50 text-ink-600 hover:border-ink-300 hover:text-ink-900',
        className,
      )}
      {...rest}
    >
      {showCheckbox ? (
        <span
          className={cn(
            'inline-flex h-3.5 w-3.5 shrink-0 items-center justify-center rounded-[3px] border transition',
            active
              ? 'border-accent-vermilion bg-accent-vermilion text-paper-50'
              : 'border-ink-300 bg-paper-50',
          )}
          aria-hidden
        >
          {active ? <Check className="h-2.5 w-2.5" strokeWidth={3} /> : null}
        </span>
      ) : null}
      {children}
    </button>
  )
}
