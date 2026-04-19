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
          搭建属于你自己的 <br className="hidden md:block" />
          <span className="font-display font-black tracking-tight text-ink-700">
            AI
          </span>{' '}
          <span className="text-accent-vermilion">知识库</span>
        </h1>

        <p className="mt-8 max-w-2xl text-lg leading-relaxed text-ink-500">
          Aurora 让你把私有文档、笔记与资料沉淀为可检索、可对话的知识库 ——
          与模型协作回答问题、生成内容，真正拥有一份专属于你的 AI 大脑。
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
