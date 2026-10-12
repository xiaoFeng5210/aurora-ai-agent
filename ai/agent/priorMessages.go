// 当前窗口上文 历史对话
package agent

import "aurora-agent/ai"

type PriorMessages struct {
	DocumentID int

	HistoryMessages []ai.Message
}
