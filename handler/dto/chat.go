package dto

type ChatRequest struct {
	Model       string              `json:"model" binding:"required"`
	ModelId     int                 `json:"modelId" binding:"required"`
	Stream      bool                `json:"stream"`
	Thinking    ChatThinkingRequest `json:"thinking" binding:"required"`
	MaxTokens   int                 `json:"max_tokens" binding:"required"`
	Temperature float64             `json:"temperature"`
	TopP        float64             `json:"top_p"`
	Tools       []ChatToolRequest   `json:"tools"`
	Prompt      []ChatPromptRequest `json:"prompt" binding:"required"`
}

type ChatThinkingRequest struct {
	Type string `json:"type" binding:"required"`
}

type ChatToolRequest struct {
	Type      string               `json:"type" binding:"required"`
	WebSearch ChatWebSearchRequest `json:"web_search" binding:"required"`
}

type ChatWebSearchRequest struct {
	SearchEngine        string `json:"search_engine" binding:"required"`
	SearchRecencyFilter string `json:"search_recency_filter" binding:"required"`
	Count               int    `json:"count" binding:"required"`
	SearchIntent        bool   `json:"search_intent"`
	SearchDomainFilter  string `json:"search_domain_filter"`
	ContentSize         string `json:"content_size" binding:"required"`
}

type ChatPromptRequest struct {
	Role            string `json:"role" binding:"required"`
	Content         string `json:"content" binding:"required"`
	FileContentList []any  `json:"fileContentList"`
}

