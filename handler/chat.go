package handler

import (
	"aurora-agent/handler/dto"
	"aurora-agent/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)



func StreamChatWithGLMController(ctx *gin.Context) {
	documentID, err := strconv.Atoi(ctx.Param("document_id"))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, err)
		return
	}
	sseCh := make(chan string, 2)

	defer func() {
		close(sseCh)
	}()

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
	err = service.ChatWithGLMStream(documentID, req, func(event service.ChatStreamEvent) {
		if ctx.Request.Context().Err() != nil {
			return
		}
		writeErr = writeSSEEvent(ctx, event.Event, event.Data)
	})


	if err != nil {
		logger.Error("chat with glm agent failed", zap.Error(err))
		fmt.Println("chat with glm agent failed", err)
	}
	if writeErr != nil {
		logger.Error("sse connection failed", zap.Error(writeErr))
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
