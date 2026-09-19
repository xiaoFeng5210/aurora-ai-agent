package handler

import (
	"aurora-agent/database/model"
	"aurora-agent/handler/dto"
	"aurora-agent/handler/vo"
	"aurora-agent/middleware"
	"aurora-agent/service/points"
	"aurora-agent/utils/enum"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetCurrentUserPoints(ctx *gin.Context) {
	balance, err := points.Query(ctx.GetInt(middleware.UID_IN_CTX))
	if err != nil {
		vo.RespondWithServiceError(ctx, err)
		return
	}
	vo.RespondSuccess(ctx, dto.PointsBalanceResponse{Balance: balance})
}

func GetUserPoints(ctx *gin.Context) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	balance, err := points.QueryByUserID(userID)
	if err != nil {
		logger.Error("get user points failed", zap.Error(err))
		vo.RespondWithServiceError(ctx, err)
		return
	}
	vo.RespondSuccess(ctx, dto.PointsBalanceResponse{Balance: balance})
}

func AddUserPoints(ctx *gin.Context) {
	adjustUserPoints(ctx, true)
}

func DeductUserPoints(ctx *gin.Context) {
	adjustUserPoints(ctx, false)
}

func adjustUserPoints(ctx *gin.Context, credit bool) {
	userID, ok := parseUserIDParam(ctx)
	if !ok {
		return
	}

	var req dto.AdjustPointsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return
	}
	req.Remark = strings.TrimSpace(req.Remark)
	if req.Remark == "" {
		vo.RespondError(ctx, http.StatusBadRequest, errors.New("remark is required"))
		return
	}

	var (
		balance model.PointsBalance
		err     error
	)
	if credit {
		balance, err = points.Add(userID, req.Amount, req.Remark, enum.TriggerMode_Admin)
	} else {
		balance, err = points.Deduct(userID, req.Amount, req.Remark, enum.TriggerMode_Admin)
	}
	if err != nil {
		logger.Error("adjust user points failed", zap.Error(err), zap.Bool("credit", credit))
		vo.RespondWithServiceError(ctx, err)
		return
	}
	vo.RespondSuccess(ctx, dto.PointsBalanceResponse{Balance: balance.BalanceAfter})
}

func parseUserIDParam(ctx *gin.Context) (int, bool) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || userID <= 0 {
		if err == nil {
			err = errors.New("invalid user id")
		}
		vo.RespondError(ctx, http.StatusBadRequest, err)
		return 0, false
	}
	return userID, true
}
