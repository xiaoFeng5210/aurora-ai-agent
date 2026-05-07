package bootstrap

import "aurora-agent/service"


func StartConsumer() {
	go func() {
		service.ConsumeRagVectorizeTask()
	}()
}
