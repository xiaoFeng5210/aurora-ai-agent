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

	"aurora-agent/ai"
	functioncall "aurora-agent/ai/llm/function-call"
	utils "aurora-agent/utils"
)

const (
	GLM_MODEL_BASE_URL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
)

type ToolCallType string

const (
	ToolCallTypeFunction ToolCallType = "function"
)

type Function struct {
	Name        string `json:"name"`
	Arguments   string `json:"arguments"`
}

type ToolCall struct {
	Id        string `json:"id"`
	Index     int `json:"index"`
	Type      string `json:"type"`
	Function  Function `json:"function"`

}

type Delta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	ToolCalls  *[]ToolCall `json:"tool_calls"`
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
}

var (
	logger = utils.Logger
)

func InitModel() *GLM {
	apiKey := os.Getenv("GLM_API_KEY")
	glm := &GLM{
		APIKey:   apiKey,
		Model:    "glm-4.7",
		MaxToken: 65536,
		Stream:   true,
	}
	return glm
}

func (glm *GLM) ChatWithGLMInStream(messages []ai.Message) (bool, []ToolCall, error) {
	requestBody := map[string]interface{}{
		"model":       glm.Model,
		"messages":    messages,
		"max_tokens":  glm.MaxToken,
		"stream":      true,
		"temperature": 1.0,
		"thinking": map[string]interface{}{
			"type": "disabled",
		},
		"tools":       functioncall.WeatherTools,
	}

	body, _ := json.Marshal(requestBody)
	fmt.Printf("requestBody: %s\n", string(body))

	request, _ := http.NewRequest("POST", GLM_MODEL_BASE_URL, bytes.NewBuffer(body))
	request.Header.Set("Authorization", "Bearer "+glm.APIKey)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return false, nil, err
	}

	defer func() {
		if recover() != nil {
			log.Printf("Recovered from panic: %v", recover())
		}

		resp.Body.Close()
	}()

	return glm.AIStreamResponseHandler(resp.Body)
}

// ai流式回答处理
func (glm *GLM) AIStreamResponseHandler(body io.Reader) (bool, []ToolCall, error) {
	scanner := bufio.NewScanner(body)
	content := ""
	toolCalls := []ToolCall{}

	needToolCall := false  // 是否需要工具调用


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
			return false, nil, err
		}

		// 除了发送，我们自己也要组装content内容
		if streamResponse.Choices[0].Delta.Role == "assistant" {
			content += streamResponse.Choices[0].Delta.Content
		}

		if streamResponse.Choices[0].Delta.ToolCalls != nil {
			resToolCalls := streamResponse.Choices[0].Delta.ToolCalls
			for _, toolCall := range *resToolCalls {
				toolCalls = append(toolCalls, toolCall)
			}
			needToolCall = true
		}
	}
	return needToolCall, toolCalls, nil
}

