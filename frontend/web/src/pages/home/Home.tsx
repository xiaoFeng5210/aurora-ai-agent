import { useNavigate } from 'react-router-dom'
import { PenLine, Tags, Sparkles } from 'lucide-react'
import { SiteHeader } from '@/components/layout/SiteHeader'
import { Button } from '@/components/ui/Button'
import { useAuth } from '@/hooks/useAuth'

const FLOW = [
  {
    icon: PenLine,
    title: '随手记',
    desc: '灵感、摘录、链接、待办，都写成一张卡片。',
    rotate: 'sm:-rotate-2',
  },
  {
    icon: Tags,
    title: '贴标签',
    desc: '打标签、建关联，零散卡片自动连成一张网。',
    rotate: 'sm:rotate-1',
  },
  {
    icon: Sparkles,
    title: 'AI 联想',
    desc: 'AI 吸收你的卡片，帮你归纳、检索与回忆。',
    rotate: 'sm:-rotate-1',
  },
]

export function Home() {
  const navigate = useNavigate()
  const { user } = useAuth()

  const onStart = () => {
    if (user) navigate('/chat')
    else navigate('/login?redirect=' + encodeURIComponent('/chat'))
  }

  return (
    <div className="min-h-screen bg-paper-50 text-ink-900">
      <SiteHeader />
      <main className="mx-auto flex max-w-5xl flex-col items-center px-6 pt-[15vh] text-center">
        <span className="inline-flex items-center gap-2 rounded-full bg-accent-vermilion/10 px-4 py-1.5 text-sm font-medium tracking-wide text-accent-vermilion">
          <span className="text-xs">●</span>
          卡片 · 连接 · 知识库
        </span>

        {/* 品牌名 Ariadne 作主标题，副行点出「阿里阿德涅之线」的寓意 */}
        <h1 className="font-display mt-7 text-6xl font-black leading-none tracking-tight text-ink-950 md:text-8xl">
          Ariadne
        </h1>
        <p className="mt-4 font-serif text-sm italic text-ink-500 md:text-base">
          阿里阿德涅之线 — 一根线索，串起你所有的卡片
        </p>

        <h2 className="font-display mt-8 text-2xl font-bold leading-snug tracking-tight text-ink-900 md:text-4xl">
          把散落的卡片，<span className="text-accent-vermilion">连成一张知识网</span>
        </h2>

        <p className="mt-6 max-w-2xl text-base leading-relaxed text-ink-500 md:text-lg">
          灵感、摘录、想法，先写成一张卡片。Ariadne 替你把它们连成网，
          <span className="text-ink-700">AI 帮你整理、联想与检索，慢慢长成你的知识库。</span>
        </p>

        <div className="mt-10 flex flex-wrap justify-center gap-4">
          <Button size="lg" variant="primary" onClick={onStart}>
            立即开始
          </Button>
          <Button size="lg" variant="ghost">
            了解 Ariadne
          </Button>
        </div>

        {/* 用三张卡片把「收集 → 整理 → 联想」的产品理念直接呈现出来 */}
        <div className="mt-20 grid w-full gap-5 sm:grid-cols-3">
          {FLOW.map(({ icon: Icon, title, desc, rotate }, i) => (
            <div
              key={title}
              className={`group relative flex flex-col items-start rounded-2xl border border-ink-200 bg-paper-100 p-6 text-left shadow-sm transition-transform duration-200 hover:-translate-y-1 hover:rotate-0 ${rotate}`}
            >
              <span className="absolute right-5 top-5 font-display text-sm font-bold text-ink-300">
                0{i + 1}
              </span>
              <span className="inline-flex h-11 w-11 items-center justify-center rounded-xl bg-accent-vermilion/10 text-accent-vermilion">
                <Icon className="h-5 w-5" />
              </span>
              <h3 className="font-display mt-5 text-xl font-bold text-ink-950">{title}</h3>
              <p className="mt-2 text-sm leading-relaxed text-ink-500">{desc}</p>
            </div>
          ))}
        </div>
      </main>

      <footer className="mt-28 py-8 text-center text-xs text-ink-500">
        © 2026 Aurora · MIT License ·{' '}
        <a
          href="https://github.com/"
          target="_blank"
          rel="noreferrer noopener"
          className="hover:text-accent-vermilion"
        >
          GitHub
        </a>
      </footer>
    </div>
  )
}
