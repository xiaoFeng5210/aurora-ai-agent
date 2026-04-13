package embedding

import "strings"

type MdDocument struct {
	data []byte // markdown 文件数据
	content string // 转换后的文本内容
	chunks []string // 分块后的文本内容
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

func (md *MdDocument) Embedding() ([]float64, error) {
	return Embed(md.content, 2048)
}
