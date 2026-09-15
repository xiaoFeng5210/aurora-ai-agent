import type { MouseEvent } from 'react'
import { Eye, EyeOff } from 'lucide-react'
import { cn } from '@/lib/cn'

export function CardVisibilityToggle({
  visible,
  pending,
  onToggle,
  className,
  size = 'sm',
  withLabel = false,
}: {
  visible: boolean
  pending?: boolean
  onToggle: () => void
  className?: string
  size?: 'sm' | 'lg'
  withLabel?: boolean
}) {
  const onClick = (e: MouseEvent) => {
    e.stopPropagation()
    if (pending) return
    onToggle()
  }

  const icon = visible ? (
    <Eye className={size === 'lg' ? 'h-4 w-4' : 'h-3.5 w-3.5'} />
  ) : (
    <EyeOff className={size === 'lg' ? 'h-4 w-4' : 'h-3.5 w-3.5'} />
  )

  return (
    <button
      type="button"
      aria-label={visible ? '隐去卡片正文' : '显示卡片正文'}
      aria-pressed={!visible}
      title={visible ? '隐去正文' : '显示正文'}
      disabled={pending}
      onClick={onClick}
      className={cn(
        'inline-flex cursor-pointer items-center justify-center text-ink-400 transition hover:bg-paper-200/70 hover:text-ink-800 disabled:cursor-not-allowed disabled:opacity-50',
        withLabel
          ? 'gap-1 rounded-md px-2 py-1 text-xs text-ink-500 hover:text-ink-800'
          : size === 'lg'
            ? 'h-9 w-9 rounded-full'
            : 'h-7 w-7 rounded-md',
        !visible && !withLabel && 'text-accent-vermilion/80 hover:text-accent-vermilion',
        className,
      )}
    >
      {icon}
      {withLabel ? (visible ? '隐去正文' : '显示正文') : null}
    </button>
  )
}
