import React, { memo, useState } from 'react'
import ReactMarkdown, { type Components } from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeHighlight from 'rehype-highlight'
import { Check, Copy } from 'lucide-react'
import { cn } from '@/lib/cn'

function flattenText(node: React.ReactNode): string {
  if (node == null || typeof node === 'boolean') return ''
  if (typeof node === 'string') return node
  if (typeof node === 'number') return String(node)
  if (Array.isArray(node)) return node.map(flattenText).join('')
  if (React.isValidElement(node)) {
    return flattenText((node.props as { children?: React.ReactNode }).children)
  }
  return ''
}

function extractCodeInfo(children: React.ReactNode): { lang?: string; code: string } {
  const first = Array.isArray(children) ? children[0] : children
  if (React.isValidElement(first)) {
    const el = first as React.ReactElement<{
      className?: string
      children?: React.ReactNode
    }>
    const className = el.props.className || ''
    const lang = /language-([\w-]+)/.exec(className)?.[1]
    return { lang, code: flattenText(el.props.children) }
  }
  return { code: flattenText(children) }
}

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false)
  const onCopy = async () => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      setTimeout(() => setCopied(false), 1400)
    } catch {
      // ignore
    }
  }
  return (
    <button
      type="button"
      onClick={onCopy}
      className="inline-flex items-center gap-1 text-paper-300/70 transition hover:text-paper-50"
      aria-label="复制代码"
    >
      {copied ? <Check className="h-3 w-3" /> : <Copy className="h-3 w-3" />}
      {copied ? '已复制' : '复制'}
    </button>
  )
}

const components: Components = {
  p: ({ children }) => (
    <p className="my-3 leading-[1.8] first:mt-0 last:mb-0">{children}</p>
  ),
  h1: ({ children }) => (
    <h1 className="font-display mt-6 mb-3 text-2xl font-bold text-ink-950 first:mt-0">
      {children}
    </h1>
  ),
  h2: ({ children }) => (
    <h2 className="font-display mt-6 mb-2.5 text-xl font-bold text-ink-950 first:mt-0">
      {children}
    </h2>
  ),
  h3: ({ children }) => (
    <h3 className="font-display mt-5 mb-2 text-lg font-semibold text-ink-950 first:mt-0">
      {children}
    </h3>
  ),
  h4: ({ children }) => (
    <h4 className="mt-4 mb-2 text-base font-semibold text-ink-950 first:mt-0">
      {children}
    </h4>
  ),
  ul: ({ children }) => (
    <ul className="my-3 list-disc space-y-1.5 pl-6 marker:text-ink-500">
      {children}
    </ul>
  ),
  ol: ({ children }) => (
    <ol className="my-3 list-decimal space-y-1.5 pl-6 marker:text-ink-500">
      {children}
    </ol>
  ),
  li: ({ children }) => <li className="leading-[1.75]">{children}</li>,
  a: ({ href, children }) => (
    <a
      href={href}
      target="_blank"
      rel="noreferrer noopener"
      className="text-accent-vermilion underline decoration-accent-vermilion/40 underline-offset-2 transition hover:decoration-accent-vermilion"
    >
      {children}
    </a>
  ),
  blockquote: ({ children }) => (
    <blockquote className="my-4 border-l-2 border-accent-amber/60 bg-paper-100/70 py-1 pl-4 text-ink-700">
      {children}
    </blockquote>
  ),
  hr: () => <hr className="my-6 border-t border-ink-200" />,
  strong: ({ children }) => (
    <strong className="font-semibold text-ink-950">{children}</strong>
  ),
  em: ({ children }) => <em className="italic text-ink-800">{children}</em>,
  table: ({ children }) => (
    <div className="my-4 overflow-x-auto rounded-lg border border-ink-200">
      <table className="w-full border-collapse text-sm">{children}</table>
    </div>
  ),
  thead: ({ children }) => (
    <thead className="bg-paper-200/60 text-ink-950">{children}</thead>
  ),
  th: ({ children }) => (
    <th className="border-b border-ink-200 px-3 py-2 text-left font-semibold">
      {children}
    </th>
  ),
  td: ({ children }) => (
    <td className="border-b border-ink-200/70 px-3 py-2 align-top text-ink-800">
      {children}
    </td>
  ),
  code: ({ className, children, ...rest }) => {
    if (className && /language-/.test(className)) {
      return (
        <code className={className} {...rest}>
          {children}
        </code>
      )
    }
    return (
      <code className="rounded bg-paper-200/70 px-1.5 py-0.5 font-mono text-[0.86em] text-accent-vermilion">
        {children}
      </code>
    )
  },
  pre: ({ children }) => {
    const { lang, code } = extractCodeInfo(children)
    return (
      <div className="my-4 overflow-hidden rounded-xl border border-ink-950/20 bg-ink-950 shadow-sm">
        <div className="flex items-center justify-between border-b border-paper-50/10 px-3 py-1.5 text-[11px] text-paper-300">
          <span className="font-mono lowercase tracking-wide">
            {lang || 'code'}
          </span>
          <CopyButton text={code} />
        </div>
        <pre className="overflow-x-auto p-4 font-mono text-[13px] leading-relaxed text-paper-200">
          {children}
        </pre>
      </div>
    )
  },
}

export const Markdown = memo(function Markdown({
  content,
  className,
}: {
  content: string
  className?: string
}) {
  return (
    <div
      className={cn(
        'markdown-body text-[15px] leading-[1.75] text-ink-900',
        className,
      )}
    >
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeHighlight]}
        components={components}
      >
        {content}
      </ReactMarkdown>
    </div>
  )
})
