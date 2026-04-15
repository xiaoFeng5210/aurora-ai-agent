package functioncall

import (
	qdrant_db "aurora-agent/database/qdrant"
	"aurora-agent/service/embedding"
	"encoding/json"
)

var RAGTools = []map[string]interface{}{
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "query_rag",
			"description": "查询RAG向量库的方法，首先你需要自主判断用户的问题，属于闲聊类的还是？如果你觉得用户的问题不属于闲聊倾向，对知识论文的问题，或者对故事问题的探索，或者对其他常用问题的询问，那么你需要调用这个方法查询RAG向量库。",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"prompt": map[string]interface{}{
						"type":        "string",
						"description": "直接将用户的问题或者相关关键词传递给这个参数，可以不进行任何处理，如果用户的问题很长，可以适当组织语言后传递给这个参数",
					},
				},
				"required": []interface{}{"prompt"},
			},
		},
	},
}


func QueryRag(prompt string) ([]byte, error) {
	queryVector, err := embedding.Embed(prompt, 1024)
	if err != nil {
		return nil, err
	}
	searchResult, err := qdrant_db.QueryRagVector(queryVector)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(searchResult)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}


