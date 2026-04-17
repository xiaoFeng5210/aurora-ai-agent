import { Link, useNavigate } from 'react-router-dom'
import { ChevronDown, Container, Globe, Music, Star } from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/Button'

const NAV = [
  { label: '功能', href: '#features' },
  { label: '展示', href: '#showcase' },
  { label: 'Agent', href: '#agent' },
  { label: '博客', href: '#blog' },
]

export function SiteHeader() {
  const { user } = useAuth()
  const navigate = useNavigate()

  return (
    <header className="sticky top-0 z-40 border-b border-ink-200/60 bg-paper-50/80 backdrop-blur">
      <div className="mx-auto grid h-16 max-w-7xl grid-cols-3 items-center px-6">
        <Link to="/" className="flex items-center gap-2">
          <Music className="h-6 w-6 text-accent-olive" strokeWidth={2.2} />
          <span className="font-display text-xl font-extrabold text-accent-vermilion">Aurora</span>
          <span className="text-ink-300">·</span>
        </Link>

        <nav className="flex items-center justify-center gap-10">
          {NAV.map((item) => (
            <a
              key={item.label}
              href={item.href}
              className="text-sm font-medium text-ink-700 transition hover:text-ink-950"
            >
              {item.label}
            </a>
          ))}
        </nav>

        <div className="flex items-center justify-end gap-3">
          {user ? (
            <>
              <span className="text-sm text-ink-700">
                你好，<span className="font-medium text-ink-950">{user.username}</span>
              </span>
              <Button size="sm" variant="primary" onClick={() => navigate('/chat')}>
                进入对话
              </Button>
            </>
          ) : (
            <>
              <Link
                to="/login"
                className="text-sm text-ink-700 transition hover:text-accent-vermilion"
              >
                登录
              </Link>
              <Link
                to="/register"
                className="text-sm text-ink-700 transition hover:text-accent-vermilion"
              >
                注册
              </Link>
            </>
          )}

          <button
            type="button"
            className="inline-flex items-center gap-1.5 rounded-full border border-ink-200 bg-paper-50 px-3 py-1.5 text-sm text-ink-800 transition hover:border-ink-300"
          >
            <Globe className="h-4 w-4" />
            中文
            <ChevronDown className="h-3.5 w-3.5 text-ink-500" />
          </button>

          <a
            href="https://hub.docker.com/"
            target="_blank"
            rel="noreferrer noopener"
            className="text-ink-500 transition hover:text-ink-800"
            aria-label="Docker"
          >
            <Container className="h-5 w-5" />
          </a>

          <a
            href="https://github.com/"
            target="_blank"
            rel="noreferrer noopener"
            className="inline-flex items-center gap-2 rounded-full bg-ink-950 px-3 py-1.5 text-sm text-paper-50 transition hover:bg-ink-900"
          >
            {/* TODO: <Github className="h-4 w-4" /> */}
            <span className="font-medium">GitHub</span>
            <span className="h-3 w-px bg-paper-50/30" />
            <Star className="h-3.5 w-3.5 fill-paper-50 text-paper-50" />
            <span className="text-xs">569</span>
          </a>
        </div>
      </div>
    </header>
  )
}
