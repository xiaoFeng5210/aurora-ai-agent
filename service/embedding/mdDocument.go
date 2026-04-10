package embedding

type MdDocument struct {
	data []byte // markdown 文件数据
	content string // 转换后的文本内容
}

func (md *MdDocument) ConvertToText() string {
  if md.data != nil {
		md.content = string(md.data)
		return md.content
	}

	return ""
}
