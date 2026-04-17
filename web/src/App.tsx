import useSWR from 'swr'
import { fetcher } from './lib/fetcher'

interface PingResponse {
  message: string
}

function App() {
  const { data, error, isLoading, mutate } = useSWR<PingResponse>('/ping', fetcher)

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 text-slate-100">
      <div className="mx-auto flex max-w-3xl flex-col items-center gap-8 px-6 py-20">
        <h1 className="text-center text-5xl font-bold tracking-tight">
          Aurora <span className="text-indigo-400">AI Agent</span>
        </h1>
        <p className="text-center text-lg text-slate-300">
          React + TypeScript + Vite + TailwindCSS + SWR
        </p>

        <div className="w-full rounded-2xl border border-slate-700 bg-slate-800/60 p-6 shadow-xl backdrop-blur">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-xl font-semibold">/ping 接口测试</h2>
            <button
              onClick={() => mutate()}
              className="rounded-lg bg-indigo-500 px-4 py-2 text-sm font-medium text-white transition hover:bg-indigo-400 active:bg-indigo-600"
            >
              重新请求
            </button>
          </div>

          {isLoading && <p className="text-slate-400">加载中...</p>}
          {error && (
            <p className="text-rose-400">
              请求失败：{(error as Error).message}
            </p>
          )}
          {data && (
            <pre className="overflow-x-auto rounded-lg bg-slate-950/70 p-4 text-sm text-emerald-300">
              {JSON.stringify(data, null, 2)}
            </pre>
          )}
        </div>
      </div>
    </div>
  )
}

export default App
