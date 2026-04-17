import useSWR from 'swr'
import { fetcher } from './lib/fetcher'

interface PingResponse {
  message: string
}

function App() {
  const { data, error, isLoading, mutate } = useSWR<PingResponse>('/ping', fetcher)

  return (
    <div className="min-h-screen bg-paper-100 text-ink-900">
      <div className="mx-auto flex max-w-3xl flex-col items-center gap-8 px-6 py-20">
        <h1 className="font-display text-center text-6xl font-light tracking-tight text-ink-950">
          Aurora <span className="italic text-accent-vermilion">AI Agent</span>
        </h1>
        <p className="text-center text-lg text-ink-500">
          React + TypeScript + Vite + TailwindCSS + SWR
        </p>

        <div className="w-full rounded-lg border border-ink-200 bg-paper-50 p-6 shadow-[0_1px_0_0_rgba(58,46,31,0.04),0_8px_24px_-12px_rgba(58,46,31,0.15)]">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="font-display text-xl font-medium text-ink-900">/ping 接口测试</h2>
            <button
              onClick={() => mutate()}
              className="rounded-md bg-accent-vermilion px-4 py-2 text-sm font-medium text-paper-50 transition hover:bg-accent-vermilion/90 active:bg-accent-vermilion/80"
            >
              重新请求
            </button>
          </div>

          {isLoading && <p className="text-ink-500">加载中...</p>}
          {error && (
            <p className="text-danger">
              请求失败：{(error as Error).message}
            </p>
          )}
          {data && (
            <pre className="overflow-x-auto rounded-md border border-ink-200 bg-paper-200/60 p-4 text-sm text-ink-800">
              {JSON.stringify(data, null, 2)}
            </pre>
          )}
        </div>
      </div>
    </div>
  )
}

export default App
