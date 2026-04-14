package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	embeddingURL   = "https://open.bigmodel.cn/api/paas/v4/embeddings"
	embeddingModel = "embedding-3"
)

type embeddingRequest struct {
	Model      string `json:"model"`
	Input      string `json:"input"`
	Dimensions int    `json:"dimensions,omitempty"`
}

type embeddingData struct {
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

type embeddingResponse struct {
	Data  []embeddingData `json:"data"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Embed 将文本向量化，dimensions 为 0 时使用模型默认维度，2048最高, 我们一般选取1024
func Embed(text string, dimensions int) ([]float32, error) {
	apiKey := os.Getenv("GLM_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GLM_API_KEY is not set")
	}

	reqBody := embeddingRequest{
		Model:      embeddingModel,
		Input:      text,
		Dimensions: dimensions,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, embeddingURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result embeddingResponse
	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("embedding api error %s: %s", result.Error.Code, result.Error.Message)
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("embedding api returned empty data")
	}

	return result.Data[0].Embedding, nil
}
