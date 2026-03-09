package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
			Content: "我想知道明朝皇帝最有名的一位是谁",
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

	buffer := make([]byte, 2048)

	for {
		n, err := resp.Body.Read(buffer)

		if err == io.EOF {
			break
		}

		var conversationId string = ""

		for _, segment := range bytes.Split(buffer[:n], []byte("data: ")) {
			segment = bytes.TrimSpace(segment)
			streamResponse := StreamResponse{}
			if len(segment) == 0 {
				continue
			}

			if string(segment) == "[DONE]" {
				fmt.Printf("conversationId: %s\n", conversationId)
				break
			}

			err = json.Unmarshal(segment, &streamResponse)
			if err != nil {
				log.Printf("不能序列化的数据: %s\n", string(segment))
				break
			}

			if conversationId == "" {
				conversationId = streamResponse.Id
			}

			fmt.Print(streamResponse.Choices[0].Delta.Content)

		}
	}

}
