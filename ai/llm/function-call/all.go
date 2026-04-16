package functioncall

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RunToolFunction(ctx *gin.Context, functionName string, functionArguments string) ([]byte, error) {
	switch functionName {
	case "get_weather":
		var arguments map[string]interface{}
		err := json.Unmarshal([]byte(functionArguments), &arguments)
		if err != nil {
			return nil, err
		}

		city, ok := arguments["city"]
		if !ok {
			return nil, fmt.Errorf("city is required")
		}
		cityStr, _ := city.(string)

		return GetWeather(cityStr), nil

	case "query_rag":
		return QueryRag(ctx, functionArguments)
	default:
		return nil, fmt.Errorf("function %s not found", functionName)
	}
}
