package dto

type ChatRequest struct {
	Thinking    ChatThinkingRequest `json:"thinking" binding:"required"`
	MaxTokens   int                 `json:"max_tokens" binding:"required"`   // 最大token数
	Temperature float64             `json:"temperature"`                     // 温度
	TopP        float64             `json:"top_p"`                           // top_p
	Tools       []ChatToolRequest   `json:"tools"`                           // 工具
	Prompt      []ChatPromptRequest `json:"prompt" binding:"required"`        // 提示
}

type ChatThinkingRequest struct {
	// type: disabled, enabled
	Type string `json:"type" binding:"required"`
}

type ChatToolRequest struct {
	Type      string               `json:"type" binding:"required"`
	WebSearch ChatWebSearchRequest `json:"web_search"`
}

type ChatWebSearchRequest struct {
	SearchEngine        string `json:"search_engine" binding:"required"`
	SearchRecencyFilter string `json:"search_recency_filter" binding:"required"`
	Count               int    `json:"count" binding:"required"`
	SearchIntent        bool   `json:"search_intent"`
	SearchDomainFilter  string `json:"search_domain_filter"`
	ContentSize         string `json:"content_size" binding:"required"`
}


// role: user, assistant, system
// content: 内容
// fileContentList: 文件内容列表
type ChatPromptRequest struct {
	Role            string `json:"role" binding:"required"`
	Content         string `json:"content" binding:"required"`
	FileContentList []any  `json:"fileContentList"`
}

