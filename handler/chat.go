package handler

import (
	"aurora-agent/handler/dto"
	"aurora-agent/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StreamChatWithGLM(ctx *gin.Context) {
	var req dto.ChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("X-Accel-Buffering", "no")
	ctx.Status(http.StatusOK)
	ctx.Writer.Flush()

	var writeErr error
	err := service.ChatWithGLMStream(req, func(event service.ChatStreamEvent) {
		if writeErr != nil || ctx.Request.Context().Err() != nil {
			return
		}

		writeErr = writeSSEEvent(ctx, event.Event, event.Data)
	})


	if err != nil {
		logger.Error("stream chat with glm failed", zap.Error(err))
		return
	}
	if writeErr != nil {
		logger.Error("write sse event failed", zap.Error(writeErr))
	}
}


func writeSSEEvent(ctx *gin.Context, event string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err = fmt.Fprintf(ctx.Writer, "event: %s\n", event); err != nil {
		return err
	}
	if _, err = fmt.Fprintf(ctx.Writer, "data: %s\n\n", payload); err != nil {
		return err
	}

	// 刷新缓冲区
	ctx.Writer.Flush()
	return nil
}
