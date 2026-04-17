export type SSEEventHandler = (event: string, data: unknown) => void

export async function parseSSEStream(
  reader: ReadableStreamDefaultReader<Uint8Array>,
  onEvent: SSEEventHandler,
  signal?: AbortSignal,
): Promise<void> {
  const decoder = new TextDecoder('utf-8')
  let buffer = ''

  const dispatch = (block: string) => {
    let event = 'message'
    const dataLines: string[] = []
    for (const rawLine of block.split('\n')) {
      const line = rawLine.replace(/\r$/, '')
      if (!line || line.startsWith(':')) continue
      const colon = line.indexOf(':')
      const field = colon === -1 ? line : line.slice(0, colon)
      const value = colon === -1 ? '' : line.slice(colon + 1).replace(/^ /, '')
      if (field === 'event') event = value
      else if (field === 'data') dataLines.push(value)
    }
    if (dataLines.length === 0) return
    const dataStr = dataLines.join('\n')
    let data: unknown = dataStr
    try {
      data = JSON.parse(dataStr)
    } catch {
      // 非 JSON 数据保留原文
    }
    onEvent(event, data)
  }

  try {
    while (true) {
      if (signal?.aborted) {
        await reader.cancel()
        return
      }
      const { value, done } = await reader.read()
      if (done) break
      buffer += decoder.decode(value, { stream: true })
      let sep: number
      while ((sep = buffer.indexOf('\n\n')) !== -1) {
        const block = buffer.slice(0, sep)
        buffer = buffer.slice(sep + 2)
        if (block.trim()) dispatch(block)
      }
    }
    if (buffer.trim()) dispatch(buffer)
  } finally {
    reader.releaseLock()
  }
}
