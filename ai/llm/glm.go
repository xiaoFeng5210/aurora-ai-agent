package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	functioncall "aurora-agent/ai/llm/function-call"
)

const (
	GLM_MODEL_BASE_URL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
)

type MessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Delta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Delta Delta `json:"delta"`
}

type StreamResponse struct {
	Id      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

type GLM struct {
	APIKey   string
	Model    string
	MaxToken int
	Stream   bool
	messages []MessageRequest
}

var (
	GLM_SYSTEM_BASE_PROMPT = []MessageRequest{
		{
			Role:    "system",
			Content: "你是一个知识库助手，你可以根据用户的问题，从知识库中查询相关信息，并返回给用户。",
		},
	}
)

func InitModel(model string) *GLM {
	glm := &GLM{
		APIKey:   os.Getenv("GLM_API_KEY"),
		Model:    "glm-4.7",
		MaxToken: 65536,
		Stream:   true,
	}
	// 初始化时候我们先设置基础提示词
	glm.messages = GLM_SYSTEM_BASE_PROMPT
	return glm
}

func (glm *GLM) ChatWithGLMInStream() {
	requestBody := map[string]interface{}{
		"model":       glm.Model,
		"messages":    glm.messages,
		"max_tokens":  glm.MaxToken,
		"stream":      true,
		"temperature": 1.0,
		"thinking": map[string]interface{}{
			"type": "disabled",
		},
		"tools":       functioncall.WeatherTools,
		// "tool_stream": true,
	}

	body, _ := json.Marshal(requestBody)

	request, _ := http.NewRequest("POST", GLM_MODEL_BASE_URL, bytes.NewBuffer(body))
	request.Header.Set("Authorization", "Bearer "+glm.APIKey)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return
	}

	defer func() {
		if recover() != nil {
			log.Printf("Recovered from panic: %v", recover())
		}

		resp.Body.Close()
	}()

	glm.AIStreamResponseHandler(resp.Body)
}

// ai流式回答处理
func (glm *GLM) AIStreamResponseHandler(body io.Reader) {
	scanner := bufio.NewScanner(body)
	content := ""
	// tool_calls := []map[string]interface{}{}


	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}


		segment := line[6:] // 去掉data: 

		log.Println(string(segment))
		log.Println("--------------------------------")


		if string(segment) == "[DONE]" {
			fmt.Printf("content: %s\n", content)
			// TODO: 这里需要处理content内容
			break
		}
		

		var streamResponse StreamResponse
		if err := json.Unmarshal(segment, &streamResponse); err != nil {
			log.Fatalf("不支持json的字段: %s", segment)
			break
		}

		// 除了发送，我们自己也要组装content内容
		if streamResponse.Choices[0].Delta.Role == "assistant" {
			content += streamResponse.Choices[0].Delta.Content
		}
	}
}


// 发送当前用户的提示词
func (glm *GLM) SendUserPrompt(prompt string) {
	glm.messages = append(glm.messages, MessageRequest{
		Role:    "user",
		Content: prompt,
	})
}

