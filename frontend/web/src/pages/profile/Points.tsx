import useSWR from 'swr'
import { ArrowUpRight, Coins } from 'lucide-react'
import { getMyPoints } from '@/api/points'
import { Button } from '@/components/ui/Button'
import './points.css'

export function Points({ userID }: { userID: number }) {
  const { data, error, isLoading, mutate } = useSWR(
    ['my-points', userID], getMyPoints, { shouldRetryOnError: false },
  )
  const balance = data?.data?.balance ?? 0

  return (
    <section className="points-panel mt-6 sm:mt-8" aria-labelledby="points-heading">
      <div className="mb-5 flex items-center gap-2">
        <Coins aria-hidden="true" className="h-4 w-4 text-accent-vermilion" />
        <h2 id="points-heading" className="text-base font-semibold text-ink-950">个人积分</h2>
      </div>
      <div className="points-balance" aria-live="polite" aria-busy={isLoading}>
        <div className="flex items-center justify-between gap-4 border-b border-paper-50/20 pb-5">
          <p className="text-sm text-paper-200">当前可用积分</p>
          <span className="text-xs tracking-[0.18em] text-paper-300">ARIADNE / POINTS</span>
        </div>
        {error ? (
          <div className="py-9" role="alert">
            <p className="text-lg">积分加载失败</p>
            <p className="mt-2 text-sm text-paper-200">请稍后重试，查看最新余额。</p>
            <Button type="button" variant="ghost" onClick={() => void mutate()} className="mt-5 border-paper-300 text-paper-50 hover:bg-ink-800 hover:text-paper-50">重新加载</Button>
          </div>
        ) : (
          <div className="py-9 sm:py-12">
            <p className="points-amount font-display font-medium tabular-nums tracking-tight">
              {isLoading ? <span className="text-2xl tracking-normal text-paper-200">正在加载…</span> : balance.toLocaleString('zh-CN')}
              {!isLoading && <span className="ml-3 text-base font-normal tracking-normal text-paper-200">积分</span>}
            </p>
            <p className="mt-5 text-sm text-paper-200">{isLoading ? '正在获取你的积分余额' : balance === 0 ? '暂无可用积分' : '积分已入账，余额以当前账户为准'}</p>
          </div>
        )}
        <div className="flex items-center gap-2 border-t border-paper-50/20 pt-5 text-xs text-paper-200">
          <ArrowUpRight aria-hidden="true" className="h-4 w-4 shrink-0" />
          <span>获得积分后，余额会自动更新。</span>
        </div>
      </div>
    </section>
  )
}
