package embedding

import (
	qdrant_db "aurora-agent/database/qdrant"
	"errors"
	"strings"

	"github.com/qdrant/go-client/qdrant"
)

type MdDocument struct {
	data []byte  // markdown 文件数据
	content string // 转换后的文本内容
	chunks []string // 分块后的文本内容
	vectors [][]float32 // 向量化后
}

func (md *MdDocument) ConvertToText() string {
  if md.data != nil {
		md.content = string(md.data)
		return md.content
	}

	return ""
}

// 分块
func (md *MdDocument) Chunk() []string {
	chunks := []string{}
	
	start := 0
	delta := 500
	end := start + delta

	runeContent := []rune(strings.TrimSpace(md.content))
	for {
		if end >= len(runeContent) {
			chunks = append(chunks, string(runeContent[start:]))
			break
		} else {
			chunks = append(chunks, string(runeContent[start:end]))
			start = end
			end += delta
		}
	}
	md.chunks = chunks
	return chunks
}

// 向量化
func (md *MdDocument) Embedding() ([][]float32, error) {
	vectors := [][]float32{}
	for _, chunk := range md.chunks {
		vector, err := Embed(chunk, 1024)
		if err != nil {
			return nil, err
		}
		vectors = append(vectors, vector)
	}
	md.vectors = vectors
	return vectors, nil
}

func (md *MdDocument) UpsertQdrantVector() (*qdrant.UpdateResult, error) {
	if md.content == "" {
		return nil, errors.New("content is empty")
	}
	if len(md.chunks) > 0 && len(md.vectors) > 0 {
	  result,err := qdrant_db.UpsertRagVector(md.chunks, md.vectors, map[string]any{"text": md.content})
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("chunks or vectors is empty")
}
