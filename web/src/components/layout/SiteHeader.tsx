import { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import {
  ChevronDown,
  Container,
  Globe,
  Menu,
  Star,
  UserRound,
  X,
} from 'lucide-react'
import { useAuth } from '@/hooks/useAuth'
import { Button } from '@/components/ui/Button'
import logoUrl from '@/assets/logo.png'

const NAV = [
  { label: '功能', href: '#features' },
  { label: '展示', href: '#showcase' },
  // { label: 'Agent', href: '#agent' },
  { label: '博客', href: '#blog' },
]

export function SiteHeader() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [open, setOpen] = useState(false)

  useEffect(() => {
    setOpen(false)
  }, [location.pathname, location.hash])

  useEffect(() => {
    if (!open) return
    const prev = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => {
      document.body.style.overflow = prev
    }
  }, [open])

  const goProfile = () => {
    setOpen(false)
    navigate('/chat')
  }

  const go = (to: string) => () => {
    setOpen(false)
    navigate(to)
  }

  return (
    <header className="sticky top-0 z-40 border-b border-ink-200/60 bg-paper-50/80 backdrop-blur">
      <div className="mx-auto grid h-16 max-w-7xl grid-cols-[auto_1fr_auto] items-center gap-4 px-4 md:grid-cols-3 md:px-6">
        <div
          role="link"
          tabIndex={0}
          onClick={go('/')}
          onKeyDown={(e) => {
            if (e.key === 'Enter' || e.key === ' ') go('/')()
          }}
          className="flex h-full w-fit cursor-pointer items-center gap-2 select-none z-10"
        >
          <img src={logoUrl} alt="Aurora" className="h-8 w-8" />
          <span className="font-display text-xl font-extrabold text-accent-vermilion">Aurora</span>
          <span className="hidden text-ink-300 md:inline">·</span>
        </div>

        <nav className="hidden items-center justify-center gap-10 md:flex">
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

        <div className="hidden items-center justify-end gap-3 md:flex">
          {user ? (
            <Button
              size="sm"
              variant="ghost"
              onClick={goProfile}
              className="!rounded-full"
            >
              <UserRound className="h-4 w-4" />
              个人中心
            </Button>
          ) : (
            <>
              <div
                role="link"
                tabIndex={0}
                onClick={go('/login')}
                className="cursor-pointer text-sm text-ink-700 transition hover:text-accent-vermilion"
              >
                登录
              </div>
              <div
                role="link"
                tabIndex={0}
                onClick={go('/register')}
                className="cursor-pointer text-sm text-ink-700 transition hover:text-accent-vermilion"
              >
                注册
              </div>
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
            <span className="font-medium">GitHub</span>
            <span className="h-3 w-px bg-paper-50/30" />
            <Star className="h-3.5 w-3.5 fill-paper-50 text-paper-50" />
            <span className="text-xs">569</span>
          </a>
        </div>

        <button
          type="button"
          onClick={() => setOpen((v) => !v)}
          aria-label={open ? '关闭菜单' : '打开菜单'}
          aria-expanded={open}
          className="inline-flex h-10 w-10 items-center justify-center rounded-md border border-ink-200 bg-paper-50 text-ink-800 transition hover:border-ink-300 md:hidden"
        >
          {open ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
        </button>
      </div>

      {open ? (
        <div className="border-t border-ink-200/60 bg-paper-50 md:hidden">
          <div className="mx-auto max-w-7xl px-4 py-4">
            <nav className="flex flex-col">
              {NAV.map((item) => (
                <a
                  key={item.label}
                  href={item.href}
                  onClick={() => setOpen(false)}
                  className="border-b border-ink-200/50 py-3 text-base font-medium text-ink-800 transition hover:text-accent-vermilion"
                >
                  {item.label}
                </a>
              ))}
            </nav>

            <div className="mt-4 flex flex-col gap-3">
              {user ? (
                <Button variant="ghost" onClick={goProfile} className="!rounded-full">
                  <UserRound className="h-4 w-4" />
                  个人中心
                </Button>
              ) : (
                <div className="grid grid-cols-2 gap-3">
                  <div
                    role="link"
                    tabIndex={0}
                    onClick={go('/login')}
                    className="inline-flex cursor-pointer items-center justify-center rounded-md border border-ink-300 px-4 py-2 text-sm font-medium text-ink-900 transition hover:bg-paper-200/60"
                  >
                    登录
                  </div>
                  <div
                    role="link"
                    tabIndex={0}
                    onClick={go('/register')}
                    className="inline-flex cursor-pointer items-center justify-center rounded-md bg-accent-vermilion px-4 py-2 text-sm font-medium text-paper-50 transition hover:bg-accent-vermilion/90"
                  >
                    注册
                  </div>
                </div>
              )}

              <div className="flex items-center justify-between pt-2">
                <button
                  type="button"
                  className="inline-flex items-center gap-1.5 rounded-full border border-ink-200 bg-paper-50 px-3 py-1.5 text-sm text-ink-800 transition hover:border-ink-300"
                >
                  <Globe className="h-4 w-4" />
                  中文
                  <ChevronDown className="h-3.5 w-3.5 text-ink-500" />
                </button>

                <div className="flex items-center gap-3">
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
                    <span className="font-medium">GitHub</span>
                    <span className="h-3 w-px bg-paper-50/30" />
                    <Star className="h-3.5 w-3.5 fill-paper-50 text-paper-50" />
                    <span className="text-xs">569</span>
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </header>
  )
}
