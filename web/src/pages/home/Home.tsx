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
          The Agent Harness
        </span>

        <h1 className="font-display mt-8 text-6xl font-black leading-[1.05] tracking-tight text-ink-950 md:text-8xl">
          模型是大脑，<span className="text-accent-vermilion">Aurora</span>{' '}
          <br className="hidden md:block" />
          是除此之外的一切。
        </h1>

        <p className="mt-8 max-w-2xl text-lg leading-relaxed text-ink-500">
          Aurora 为 AI Agent 提供完整的项目迭代环境 —— 从需求澄清到任务验收的结构化流水线 —— 让 Agent
          团队能交付项目，而不只是写代码。AI 提议，人类把关。
        </p>

        <div className="mt-10 flex flex-wrap justify-center gap-4">
          <Button size="lg" variant="primary" onClick={onStart}>
            立即开始
          </Button>
          <Button size="lg" variant="ghost">
            观看演示
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
