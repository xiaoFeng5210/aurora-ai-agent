package dto

type AdjustPointsRequest struct {
	Amount int    `json:"amount" binding:"required,gt=0"`
	Remark string `json:"remark" binding:"required"`
}

type PointsBalanceResponse struct {
	Balance int `json:"balance"`
}
