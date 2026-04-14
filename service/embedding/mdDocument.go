package embedding

import (
	qdrant_db "aurora-agent/database/qdrant"
	"errors"
	"strings"
)

type MdDocument struct {
	data []byte  // markdown 文件数据
	content string // 转换后的文本内容
	chunks []string // 分块后的文本内容
	vectors []float32 // 向量化后
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
	return chunks
}

// 向量化
func (md *MdDocument) Embedding() ([]float32, error) {
	vectors, err := Embed(md.content, 2048)
	if err != nil {
		return nil, err
	}
	md.vectors = vectors
	return vectors, nil
}

func (md *MdDocument) UpsertQdrantVector() error {
	if md.content == "" {
		return errors.New("content is empty")
	}
	if len(md.chunks) > 0 && len(md.vectors) > 0 {
	  err := qdrant_db.UpsertRagVector(md.chunks, md.vectors, map[string]any{"text": md.content})
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("chunks or vectors is empty")
}
