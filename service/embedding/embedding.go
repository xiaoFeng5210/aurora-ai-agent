package embedding

import (
	"aurora-agent/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
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

var logger *zap.Logger

// Embed 将文本向量化，dimensions 为 0 时使用模型默认维度，2048最高, 我们一般选取1024
func Embed(text string, dimensions int) ([]float32, error) {
	logger = utils.InitZap("log/zap")
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
		logger.Error("embedding api error", zap.Any("error", result.Error))
		return nil, fmt.Errorf("embedding api error %s: %s", result.Error.Code, result.Error.Message)
	}
	if len(result.Data) == 0 {
		logger.Error("embedding api returned empty data")
		return nil, fmt.Errorf("embedding api returned empty data")
	}

	return result.Data[0].Embedding, nil
}
