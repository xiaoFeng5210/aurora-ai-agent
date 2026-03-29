package functioncall

import "fmt"

func RunToolFunction(functionName string, functionArguments map[string]any) (any, error) {
  switch functionName {
  case "get_weather":
		city, ok := functionArguments["city"]
		if !ok {
			return nil, fmt.Errorf("city is required")
		}
		cityStr, _ := city.(string)
    return GetWeather(cityStr), nil
  default:
    return nil, fmt.Errorf("function %s not found", functionName)
  }
}
