package service

import (
	"testing"
)

func TestSendRagTask(t *testing.T) {
	msg := []byte("test message")
	SendRagTask(msg)
}

func TestConsumeRagTask(t *testing.T) {
	ConsumeRagTask()
}
