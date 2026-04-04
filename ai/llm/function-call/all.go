package functioncall

import (
	"encoding/json"
	"fmt"
)

func RunToolFunction(functionName string, functionArguments string) ([]byte, error) {
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
	default:
		return nil, fmt.Errorf("function %s not found", functionName)
	}
}
