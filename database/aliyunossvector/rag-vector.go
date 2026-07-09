package aliyunossvector

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/vectors"
)

type UpdateResult struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	RequestID  string `json:"request_id,omitempty"`
	Count      int    `json:"count"`
}

type ScoredVector struct {
	Key      string         `json:"key,omitempty"`
	Score    float64        `json:"score,omitempty"`
	Distance float64        `json:"distance,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Data     map[string]any `json:"data,omitempty"`
	Raw      map[string]any `json:"raw,omitempty"`
}

func EnsureRagResources() error {
	c, cfg, err := client()
	if err != nil {
		return err
	}
	ctx := context.Background()
	if _, err = c.GetVectorBucket(ctx, &vectors.GetVectorBucketRequest{
		Bucket: oss.Ptr(cfg.Bucket),
	}); err != nil {
		if !isServiceErrorCode(err, 404, "NoSuchBucket") {
			return err
		}
		if _, err = c.PutVectorBucket(ctx, &vectors.PutVectorBucketRequest{
			Bucket: oss.Ptr(cfg.Bucket),
		}); err != nil {
			return err
		}
	}

	if _, err = c.GetVectorIndex(ctx, &vectors.GetVectorIndexRequest{
		Bucket:    oss.Ptr(cfg.Bucket),
		IndexName: oss.Ptr(cfg.IndexName),
	}); err != nil {
		if !isServiceStatus(err, 404) {
			return err
		}
		_, err = c.PutVectorIndex(ctx, &vectors.PutVectorIndexRequest{
			Bucket:         oss.Ptr(cfg.Bucket),
			IndexName:      oss.Ptr(cfg.IndexName),
			DataType:       oss.Ptr("float32"),
			Dimension:      oss.Ptr(cfg.Dimension),
			DistanceMetric: oss.Ptr("cosine"),
			Metadata: map[string]any{
				"nonFilterableMetadataKeys": []string{"text"},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func UpsertCardVector(chunks []string, vectorData [][]float32, payload map[string]interface{}) (*UpdateResult, error) {
	c, cfg, err := client()
	if err != nil {
		return nil, err
	}
	if len(chunks) == 0 || len(vectorData) == 0 {
		return nil, fmt.Errorf("chunks or vectors is empty")
	}
	if len(chunks) != len(vectorData) {
		return nil, fmt.Errorf("chunks and vectors length mismatch")
	}

	userID := fmt.Sprintf("%v", payload["user_id"])
	cardID := fmt.Sprintf("%v", payload["card_id"])

	baseKey := ragBaseKey(userID, cardID)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	items := make([]map[string]any, 0, len(chunks))
	for idx := range chunks {
		items = append(items, map[string]any{
			"key": fmt.Sprintf("%s:%06d", baseKey, idx),
			"data": map[string]any{
				"float32": vectorData[idx],
			},
			"metadata": map[string]any{
				"text":        chunks[idx],
				"user_id":     userID,
				"card_id":     cardID,
				"card_title":  payload["card_title"],
				"chunk_index": idx,
				"updated_at":  now,
			},
		})
	}

	result, err := c.PutVectors(context.Background(), &vectors.PutVectorsRequest{
		Bucket:    oss.Ptr(cfg.Bucket),
		IndexName: oss.Ptr(cfg.IndexName),
		Vectors:   items,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateResult{
		StatusCode: result.StatusCode,
		Status:     result.Status,
		RequestID:  result.Headers.Get("X-Oss-Request-Id"),
		Count:      len(items),
	}, nil

}

func UpsertRagVector(chunks []string, vectorData [][]float32, payload map[string]any) (*UpdateResult, error) {
	c, cfg, err := client()
	if err != nil {
		return nil, err
	}
	if len(chunks) == 0 || len(vectorData) == 0 {
		return nil, fmt.Errorf("chunks or vectors is empty")
	}
	if len(chunks) != len(vectorData) {
		return nil, fmt.Errorf("chunks and vectors length mismatch")
	}

	userID := fmt.Sprintf("%v", payload["user_id"])
	filename := strings.TrimSpace(fmt.Sprintf("%v", payload["filename"]))
	if userID == "" || userID == "<nil>" {
		return nil, fmt.Errorf("user_id is required")
	}

	if filename != "" {
		if _, err = DeleteRagVectorByFilenames(toInt(userID), []string{filename}); err != nil {
			return nil, err
		}
	}

	baseKey := ragBaseKey(userID, filename)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	items := make([]map[string]any, 0, len(chunks))
	for idx := range chunks {
		items = append(items, map[string]any{
			"key": fmt.Sprintf("%s:%06d", baseKey, idx),
			"data": map[string]any{
				"float32": vectorData[idx],
			},
			"metadata": map[string]any{
				"text":        chunks[idx],
				"user_id":     userID,
				"filename":    filename,
				"chunk_index": idx,
				"updated_at":  now,
			},
		})
	}

	result, err := c.PutVectors(context.Background(), &vectors.PutVectorsRequest{
		Bucket:    oss.Ptr(cfg.Bucket),
		IndexName: oss.Ptr(cfg.IndexName),
		Vectors:   items,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateResult{
		StatusCode: result.StatusCode,
		Status:     result.Status,
		RequestID:  result.Headers.Get("X-Oss-Request-Id"),
		Count:      len(items),
	}, nil
}

func DeleteRagVectorByFilenames(userID int, filenames []string) (*UpdateResult, error) {
	c, cfg, err := client()
	if err != nil {
		return nil, err
	}
	filenameSet := make(map[string]struct{}, len(filenames))
	for _, filename := range filenames {
		filename = strings.TrimSpace(filename)
		if filename != "" {
			filenameSet[filename] = struct{}{}
		}
	}
	if len(filenameSet) == 0 {
		return nil, fmt.Errorf("filename is required")
	}

	keys, err := listKeysByMetadata(c, cfg, strconv.Itoa(userID), filenameSet)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return &UpdateResult{StatusCode: 204, Status: "NoContent", Count: 0}, nil
	}

	result, err := c.DeleteVectors(context.Background(), &vectors.DeleteVectorsRequest{
		Bucket:    oss.Ptr(cfg.Bucket),
		IndexName: oss.Ptr(cfg.IndexName),
		Keys:      keys,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateResult{
		StatusCode: result.StatusCode,
		Status:     result.Status,
		RequestID:  result.Headers.Get("X-Oss-Request-Id"),
		Count:      len(keys),
	}, nil
}

func QueryRagVector(queryVector []float32, userID int) ([]ScoredVector, error) {
	c, cfg, err := client()
	if err != nil {
		return nil, err
	}
	result, err := c.QueryVectors(context.Background(), &vectors.QueryVectorsRequest{
		Bucket:    oss.Ptr(cfg.Bucket),
		IndexName: oss.Ptr(cfg.IndexName),
		QueryVector: map[string]any{
			"float32": queryVector,
		},
		Filter: map[string]any{
			"$and": []map[string]any{
				{
					"user_id": map[string]any{
						"$in": []string{strconv.Itoa(userID)},
					},
				},
			},
		},
		ReturnMetadata: oss.Ptr(true),
		ReturnDistance: oss.Ptr(true),
		TopK:           oss.Ptr(cfg.TopK),
	})
	if err != nil {
		return nil, err
	}

	searchResult := make([]ScoredVector, 0, len(result.Vectors))
	for _, item := range result.Vectors {
		searchResult = append(searchResult, mapScoredVector(item))
	}
	return searchResult, nil
}

func listKeysByMetadata(c *vectors.VectorsClient, cfg Config, userID string, filenames map[string]struct{}) ([]string, error) {
	keys := make([]string, 0)
	var nextToken *string
	for {
		result, err := c.ListVectors(context.Background(), &vectors.ListVectorsRequest{
			Bucket:         oss.Ptr(cfg.Bucket),
			IndexName:      oss.Ptr(cfg.IndexName),
			MaxResults:     1000,
			NextToken:      nextToken,
			ReturnMetadata: oss.Ptr(true),
			ReturnData:     oss.Ptr(false),
		})
		if err != nil {
			return nil, err
		}
		for _, item := range result.Vectors {
			metadata, _ := item["metadata"].(map[string]any)
			if fmt.Sprintf("%v", metadata["user_id"]) != userID {
				continue
			}
			if _, ok := filenames[fmt.Sprintf("%v", metadata["filename"])]; !ok {
				continue
			}
			if key, ok := item["key"].(string); ok && key != "" {
				keys = append(keys, key)
			}
		}
		if result.NextToken == nil || *result.NextToken == "" {
			break
		}
		nextToken = result.NextToken
	}
	return keys, nil
}

func mapScoredVector(item map[string]any) ScoredVector {
	metadata, _ := item["metadata"].(map[string]any)
	data, _ := item["data"].(map[string]any)
	return ScoredVector{
		Key:      fmt.Sprintf("%v", item["key"]),
		Score:    toFloat64(item["score"]),
		Distance: toFloat64(item["distance"]),
		Metadata: metadata,
		Data:     data,
		Raw:      item,
	}
}

func ragBaseKey(userID string, filename string) string {
	sum := sha256.Sum256([]byte(userID + ":" + filename))
	return "rag:" + hex.EncodeToString(sum[:])[:24]
}

func isServiceErrorCode(err error, statusCode int, codes ...string) bool {
	var serviceErr *oss.ServiceError
	if !errors.As(err, &serviceErr) {
		return false
	}
	if serviceErr.StatusCode != statusCode {
		return false
	}
	for _, code := range codes {
		if serviceErr.Code == code {
			return true
		}
	}
	return false
}

func isServiceStatus(err error, statusCode int) bool {
	var serviceErr *oss.ServiceError
	if !errors.As(err, &serviceErr) {
		return false
	}
	return serviceErr.StatusCode == statusCode
}

func toFloat64(value any) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

func toInt(value string) int {
	i, _ := strconv.Atoi(value)
	return i
}
