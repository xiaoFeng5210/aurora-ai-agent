package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"

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

func (glm *GLM) ChatWithGLMInStream() {
	glm.messages = []MessageRequest{
		{
			Role:    "system",
			Content: "你是一个有用的AI助手。",
		},
		{
			Role:    "user",
			Content: "上海天气怎么样？",
		},
	}

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
		"tool_stream": true,
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

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		segment := line[6:]

		if string(segment) == "[DONE]" {
			break
		}

		log.Printf("segment: %s\n", string(segment))

		var streamResponse StreamResponse
		if err := json.Unmarshal(segment, &streamResponse); err != nil {
			log.Fatalf("不支持json的字段: %s", segment)
			break
		}
	}
}

// buffer := make([]byte, 2048)

// 	for {
// 		n, err := resp.Body.Read(buffer)

// 		if err == io.EOF {
// 			break
// 		}

// 		var conversationId string = ""

// 		for _, segment := range bytes.Split(buffer[:n], []byte("data: ")) {
// 			segment = bytes.TrimSpace(segment)
// 			log.Printf("segment: %s\n", string(segment))
// 			streamResponse := StreamResponse{}
// 			if len(segment) == 0 {
// 				continue
// 			}

// 			if string(segment) == "[DONE]" {
// 				fmt.Printf("conversationId: %s\n", conversationId)
// 				break
// 			}

// 			err = json.Unmarshal(segment, &streamResponse)
// 			if err != nil {
// 				// log.Printf("不能序列化的数据: %s\n", string(segment))
// 				break
// 			}

// 			if conversationId == "" {
// 				conversationId = streamResponse.Id
// 			}

// 			// fmt.Print(streamResponse.Choices[0].Delta.Content)

// 		}

// 	}
