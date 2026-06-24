import { useNavigate } from 'react-router-dom'
import { SiteHeader } from '@/components/layout/SiteHeader'
import { Button } from '@/components/ui/Button'
import { useAuth } from '@/hooks/useAuth'

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
      <main className="mx-auto flex max-w-5xl flex-col items-center px-6 pt-[18vh] text-center">
        <span className="inline-flex items-center gap-2 rounded-full bg-accent-vermilion/10 px-4 py-1.5 text-sm font-medium text-accent-vermilion">
          <span className="text-xs">●</span>
          Ariadne · 卡片笔记仓库
        </span>

        <h1 className="font-display mt-8 text-6xl font-black leading-[1.05] tracking-tight text-ink-950 md:text-8xl">
          把灵感、资料与想法，<br className="hidden md:block" />
          <span className="text-accent-vermilion">都装进卡片里</span>
        </h1>

        <p className="mt-8 max-w-3xl text-lg leading-relaxed text-ink-500 md:text-xl">
          Ariadne 是一款卡片笔记仓库，让你把灵感、摘录、任务、链接和零散记录
          统一沉淀成可打标签、可整理、可复用的卡片。AI 会持续吸收这些卡片，
          帮你归纳、联想与检索，逐步形成你的知识库与记忆库。
        </p>

        <div className="mt-10 flex flex-wrap justify-center gap-4">
          <Button size="lg" variant="primary" onClick={onStart}>
            立即开始
          </Button>
          <Button size="lg" variant="ghost">
            了解 Ariadne
          </Button>
        </div>
      </main>

      <footer className="mt-32 py-8 text-center text-xs text-ink-500">
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
