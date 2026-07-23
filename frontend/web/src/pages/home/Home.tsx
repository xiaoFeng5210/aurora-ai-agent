import { useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  FileText,
  Lightbulb,
  MessageCircle,
  Network,
  PenLine,
  Sparkles,
  Tags,
  Upload,
} from 'lucide-react'
import { SiteHeader } from '@/components/layout/SiteHeader'
import { Button } from '@/components/ui/Button'
import { useAuth } from '@/hooks/useAuth'
import logoUrl from '@/assets/logo.png'

const FILES = [
  { name: '产品手册.pdf', size: '2.4 MB' },
  { name: '读书笔记.md', size: '18 KB' },
  { name: '会议纪要.docx', size: '96 KB' },
]

const ABILITIES = [
  {
    icon: MessageCircle,
    title: '问答',
    desc: '提问时，它带着你的卡片与资料回答，而不是泛泛而谈。',
  },
  {
    icon: Network,
    title: '归纳',
    desc: '把散落的输入串成脉络，主题与结构自己长出来。',
  },
  {
    icon: Lightbulb,
    title: '联想',
    desc: '新想法落笔时，相关的旧卡片自动浮现。',
  },
]

const FLOW = [
  {
    icon: PenLine,
    no: '壹',
    title: '随手记',
    desc: '灵感、摘录、链接、待办，都写成一张卡片。',
  },
  {
    icon: Tags,
    no: '贰',
    title: '贴标签',
    desc: '打标签、建关联，零散卡片自动连成一张网。',
  },
  {
    icon: Sparkles,
    no: '叁',
    title: 'AI 联想',
    desc: 'AI 吸收你的卡片，帮你归纳、检索与回忆。',
  },
]

export function Home() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const contentRef = useRef<HTMLElement>(null)

  const onStart = () => {
    if (user) navigate('/cards')
    else navigate('/login?redirect=' + encodeURIComponent('/cards'))
  }

  const onLearnMore = () => {
    contentRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }

  return (
    <div className="relative min-h-screen overflow-x-clip bg-paper-50 text-ink-900">
      <SiteHeader />

      <main className="relative mx-auto max-w-5xl px-6">
        {/* ──────── Hero ──────── */}
        <section className="flex flex-col items-center pb-24 pt-[14vh] text-center">
          <div className="flex items-center gap-4 text-xs tracking-[0.35em] text-ink-500">
            <span aria-hidden className="h-px w-8 bg-ink-300" />
            卡片 · 知识库 · AI
            <span aria-hidden className="h-px w-8 bg-ink-300" />
          </div>

          {/* 品牌名 Ariadne 作主标题，副行点出「阿里阿德涅之线」的寓意 */}
          <h1 className="font-display mt-8 text-6xl font-semibold leading-none tracking-tight text-ink-950 md:text-8xl">
            Ariadne
          </h1>
          <p className="mt-5 font-serif text-sm italic tracking-wide text-ink-500 md:text-base">
            阿里阿德涅之线 — 一根线索，串起你所有的卡片
          </p>

          <h2 className="font-display mt-14 text-2xl font-semibold leading-snug tracking-tight text-ink-900 md:text-4xl">
            把散落的卡片，<span className="text-accent-vermilion">连成一张知识网</span>
          </h2>

          <p className="mt-6 max-w-2xl text-base leading-loose text-ink-500 md:text-lg">
            零碎的灵感，写成一张卡片；完整的资料，收进知识库。
            <span className="text-ink-700">它们都汇入同一个 AI —— 帮你整理、联想与检索，慢慢长成你的知识网。</span>
          </p>

          <div className="mt-12 flex flex-wrap justify-center gap-4">
            <Button size="lg" variant="primary" onClick={onStart} className="tracking-wide">
              立即开始
            </Button>
            <Button size="lg" variant="ghost" onClick={onLearnMore} className="tracking-wide">
              了解 Ariadne
            </Button>
          </div>
        </section>

        <Divider />

        {/* ──────── 两种收集方式 ──────── */}
        <section ref={contentRef} className="scroll-mt-24 pb-28 pt-20">
          <p className="text-center text-xs tracking-[0.35em] text-ink-500">收集</p>
          <h2 className="font-display mt-4 text-center text-2xl font-semibold tracking-tight text-ink-950 md:text-3xl">
            碎片与成文，各归其位
          </h2>
          <p className="mt-3 text-center text-sm leading-relaxed text-ink-500">
            三言两语的灵感写成卡片，完整的资料收进知识库。
          </p>

          <div className="mt-14 grid gap-14 md:grid-cols-2 md:gap-10">
            {/* 卡片 · 碎片灵感 */}
            <div>
              <div className="flex items-center gap-2.5">
                <PenLine className="h-4 w-4 text-ink-400" strokeWidth={1.5} />
                <h3 className="font-display text-lg font-semibold text-ink-950">
                  卡片 · 记下灵光一闪
                </h3>
              </div>
              <p className="mt-2 text-sm leading-relaxed text-ink-500">
                灵感、摘录、待办，三言两语即成一卡。打上标签，彼此连成网。
              </p>

              <div className="mt-6 rounded-xl border border-ink-200 bg-paper-100/60 p-5">
                <div className="flex items-center justify-between text-xs text-ink-400">
                  <span>今天 14:20</span>
                  <span className="rounded-full border border-ink-200 px-2 py-0.5">灵感</span>
                </div>
                <p className="mt-3 text-sm leading-relaxed text-ink-800">
                  把读到的句子抄下来，旁边写上自己的话 —— 笔记的意义在于对话。
                </p>
                <div className="mt-4 flex gap-2 text-xs text-ink-500">
                  <span className="rounded-full bg-paper-200/70 px-2 py-0.5">摘抄</span>
                  <span className="rounded-full bg-paper-200/70 px-2 py-0.5">写作</span>
                </div>
              </div>
            </div>

            {/* 知识库 · 文件资料 */}
            <div>
              <div className="flex items-center gap-2.5">
                <FileText className="h-4 w-4 text-ink-400" strokeWidth={1.5} />
                <h3 className="font-display text-lg font-semibold text-ink-950">
                  知识库 · 收纳完整资料
                </h3>
              </div>
              <p className="mt-2 text-sm leading-relaxed text-ink-500">
                PDF、Word、Markdown，直接上传。AI 读懂它们，成为可检索的知识。
              </p>

              <div className="mt-6 overflow-hidden rounded-xl border border-ink-200 bg-paper-100/60">
                <ul className="divide-y divide-ink-200/70 text-sm">
                  {FILES.map((file) => (
                    <li key={file.name} className="flex items-center gap-3 px-4 py-3">
                      <FileText className="h-4 w-4 shrink-0 text-ink-400" strokeWidth={1.5} />
                      <span className="flex-1 truncate text-ink-800">{file.name}</span>
                      <span className="text-xs text-ink-400">{file.size}</span>
                    </li>
                  ))}
                </ul>
                <div className="flex items-center justify-center gap-2 border-t border-dashed border-ink-300 px-4 py-3 text-xs text-ink-500">
                  <Upload className="h-3.5 w-3.5" strokeWidth={1.5} />
                  拖拽或点击，上传文件
                </div>
              </div>
            </div>
          </div>
        </section>

        <Divider />

        {/* ──────── 汇入同一个 AI ──────── */}
        <section className="pb-28 pt-20 text-center">
          <p className="text-xs tracking-[0.35em] text-ink-500">汇聚</p>
          <h2 className="font-display mt-4 text-2xl font-semibold tracking-tight text-ink-950 md:text-3xl">
            汇入同一个 AI
          </h2>
          <p className="mx-auto mt-3 max-w-xl text-sm leading-relaxed text-ink-500">
            你只管收集，连接交给 AI —— 它记得你的卡片，也读过你的文件。
          </p>

          <div className="mt-14 flex flex-col items-center">
            <div className="grid w-full max-w-xl grid-cols-2 gap-4 text-sm">
              <span className="rounded-full border border-ink-200 bg-paper-100/60 px-4 py-2 text-ink-700">
                卡片 · 碎片灵感
              </span>
              <span className="rounded-full border border-ink-200 bg-paper-100/60 px-4 py-2 text-ink-700">
                知识库 · 文件资料
              </span>
            </div>

            {/* 两条曲线汇入一点 */}
            <svg
              viewBox="0 0 640 64"
              className="w-full max-w-xl text-ink-300"
              aria-hidden
              fill="none"
            >
              <path
                d="M160 0 C 160 36, 320 28, 320 64"
                stroke="currentColor"
                strokeWidth="1"
                vectorEffect="non-scaling-stroke"
              />
              <path
                d="M480 0 C 480 36, 320 28, 320 64"
                stroke="currentColor"
                strokeWidth="1"
                vectorEffect="non-scaling-stroke"
              />
            </svg>

            <div className="flex h-16 w-16 items-center justify-center rounded-full border border-accent-vermilion/40 bg-paper-50 shadow-sm">
              <img src={logoUrl} alt="Ariadne" className="h-8 w-8" />
            </div>
            <p className="font-display mt-4 text-lg font-semibold text-ink-950">Ariadne AI</p>
            <p className="mt-1 text-xs tracking-wide text-ink-500">记住卡片，读懂文件</p>

            <div className="mt-12 grid w-full gap-10 text-left sm:grid-cols-3 sm:gap-6">
              {ABILITIES.map(({ icon: Icon, title, desc }) => (
                <div key={title} className="group">
                  <span
                    aria-hidden
                    className="block h-px w-10 bg-ink-300 transition-all duration-300 group-hover:w-16 group-hover:bg-accent-vermilion"
                  />
                  <div className="mt-5 flex items-center gap-2.5">
                    <Icon
                      className="h-4 w-4 text-ink-400 transition-colors duration-300 group-hover:text-accent-vermilion"
                      strokeWidth={1.5}
                    />
                    <h3 className="font-display text-base font-semibold text-ink-950">{title}</h3>
                  </div>
                  <p className="mt-2 text-sm leading-relaxed text-ink-500">{desc}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        <Divider />

        {/* ──────── 用「壹 / 贰 / 叁」三步呈现「收集 → 整理 → 联想」的产品理念 ──────── */}
        <section className="pb-28 pt-20">
          <h2 className="font-display text-center text-2xl font-semibold tracking-tight text-ink-950 md:text-3xl">
            从一张卡片开始
          </h2>
          <p className="mt-3 text-center text-sm leading-relaxed text-ink-500">
            收集、整理、联想 —— 三步，让碎片长成体系
          </p>

          <div className="mt-14 grid sm:grid-cols-3 sm:divide-x sm:divide-ink-200">
            {FLOW.map(({ icon: Icon, no, title, desc }, i) => (
              <div
                key={title}
                className={`group flex flex-col items-start py-10 sm:px-10 sm:py-2 ${
                  i === 0 ? 'sm:pl-0' : ''
                } ${i === FLOW.length - 1 ? 'sm:pr-0' : ''} ${
                  i > 0 ? 'border-t border-ink-200 sm:border-t-0' : ''
                }`}
              >
                <span
                  aria-hidden
                  className="block h-px w-10 bg-ink-300 transition-all duration-300 group-hover:w-16 group-hover:bg-accent-vermilion"
                />
                <div className="mt-6 flex items-baseline gap-3">
                  <span className="font-display text-lg font-medium text-accent-vermilion">
                    {no}
                  </span>
                  <h3 className="font-display text-xl font-semibold text-ink-950">{title}</h3>
                  <Icon
                    className="h-4 w-4 self-center text-ink-400 transition-colors duration-300 group-hover:text-accent-vermilion"
                    strokeWidth={1.5}
                  />
                </div>
                <p className="mt-3 text-sm leading-relaxed text-ink-500">{desc}</p>
              </div>
            ))}
          </div>
        </section>

        <Divider />

        {/* ──────── 行动召唤 ──────── */}
        <section className="pb-28 pt-20 text-center">
          <h2 className="font-display text-2xl font-semibold tracking-tight text-ink-950 md:text-3xl">
            现在，放下第一张卡片
          </h2>
          <p className="mt-3 text-sm leading-relaxed text-ink-500">
            剩下的连接，交给 Ariadne。
          </p>
          <div className="mt-8">
            <Button size="lg" variant="primary" onClick={onStart} className="tracking-wide">
              立即开始
            </Button>
          </div>
        </section>
      </main>

      <footer className="relative border-t border-ink-200/70 py-8 text-center text-xs text-ink-500">
        © 2026 Aurora · MIT License ·{' '}
        <a
          href="https://github.com/"
          target="_blank"
          rel="noreferrer noopener"
          className="transition-colors hover:text-accent-vermilion"
        >
          GitHub
        </a>
      </footer>
    </div>
  )
}

function Divider() {
  return (
    <div aria-hidden className="flex items-center gap-3">
      <span className="h-px flex-1 bg-ink-200" />
      <span className="h-1.5 w-1.5 rotate-45 bg-accent-vermilion/70" />
      <span className="h-px flex-1 bg-ink-200" />
    </div>
  )
}
