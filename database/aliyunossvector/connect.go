package aliyunossvector

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/vectors"
)

const (
	defaultIndexName = "rag"
	defaultDimension = 1024
	defaultTopK      = 10
)

var (
	vectorClient *vectors.VectorsClient
	vectorConfig Config
	vectorOnce   sync.Once
	vectorErr    error
)

type Config struct {
	Region    string
	Endpoint  string
	AccountID string
	Bucket    string
	IndexName string
	Dimension int
	TopK      int
}

func Connect() error {
	vectorOnce.Do(func() {
		vectorConfig, vectorErr = loadConfig()
		if vectorErr != nil {
			return
		}

		accessKeyID := firstEnv("OSS_ACCESS_KEY_ID", "ALIBABA_CLOUD_ACCESS_KEY_ID", "ALIBABA_CLOUD_ACCESS_KEY")
		accessKeySecret := firstEnv("OSS_ACCESS_KEY_SECRET", "ALIBABA_CLOUD_ACCESS_KEY_SECRET")
		if accessKeyID == "" || accessKeySecret == "" {
			vectorErr = fmt.Errorf("OSS_ACCESS_KEY_ID or OSS_ACCESS_KEY_SECRET is not set")
			return
		}

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret)).
			WithRegion(vectorConfig.Region).
			WithAccountId(vectorConfig.AccountID)
		if vectorConfig.Endpoint != "" {
			cfg.WithEndpoint(vectorConfig.Endpoint)
		}

		vectorClient = vectors.NewVectorsClient(cfg)
	})
	return vectorErr
}

func CurrentConfig() Config {
	return vectorConfig
}

func loadConfig() (Config, error) {
	cfg := Config{
		Region:    strings.TrimSpace(os.Getenv("OSS_VECTOR_REGION")),
		Endpoint:  strings.TrimSpace(os.Getenv("OSS_VECTOR_ENDPOINT")),
		AccountID: strings.TrimSpace(os.Getenv("OSS_VECTOR_ACCOUNT_ID")),
		Bucket:    strings.TrimSpace(os.Getenv("OSS_VECTOR_BUCKET")),
		IndexName: strings.TrimSpace(os.Getenv("OSS_VECTOR_INDEX")),
		Dimension: envInt("OSS_VECTOR_DIMENSION", defaultDimension),
		TopK:      envInt("OSS_VECTOR_TOP_K", defaultTopK),
	}
	if cfg.IndexName == "" {
		cfg.IndexName = defaultIndexName
	}
	if cfg.TopK <= 0 {
		cfg.TopK = defaultTopK
	}
	if cfg.Dimension <= 0 {
		cfg.Dimension = defaultDimension
	}
	if cfg.Region == "" {
		return Config{}, fmt.Errorf("OSS_VECTOR_REGION is not set")
	}
	if cfg.AccountID == "" {
		return Config{}, fmt.Errorf("OSS_VECTOR_ACCOUNT_ID is not set")
	}
	if cfg.Bucket == "" {
		return Config{}, fmt.Errorf("OSS_VECTOR_BUCKET is not set")
	}
	return cfg, nil
}

func client() (*vectors.VectorsClient, Config, error) {
	if err := Connect(); err != nil {
		return nil, Config{}, err
	}
	return vectorClient, vectorConfig, nil
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func envInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
