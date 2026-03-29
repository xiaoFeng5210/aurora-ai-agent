package functioncall

import "encoding/json"

var WeatherTools = []map[string]interface{}{
	{
		"type": "function",
		"function": map[string]interface{}{
			"name":        "get_weather",
			"description": "获取指定城市的当前天气信息...",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"city": map[string]interface{}{
						"type":        "string",
						"description": "城市名称，例如：北京、上海...",
					},
				},
				"required": []interface{}{"city"},
			},
		},
	},
}

func GetWeather(city string) []byte {
	weather_data := map[string]string{
		"city":        city,
		"temperature": "22°C",
		"condition":   "晴天",
		"humidity":    "65%",
		"wind_speed":  "5 km/h",
	}

	jsonData, err := json.Marshal(weather_data)
	if err != nil {
		return nil
	}

	return jsonData
}
