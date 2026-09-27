package utils

// 积分对应关系

import "math"

// 积分转化成人民币
func PointsToRMB(points int) float64 {
	return math.Round(float64(points) / 1000)
}

// token转化成积分
func TokensToPoints(tokens int) int {
	// 1积分=1000token
	return int(math.Round(float64(tokens / 1000)))
}

// token转化人民币
func TokensToRMB(tokens int) float64 {
	return PointsToRMB(TokensToPoints(tokens))
}
