package messagehandler

import (
	"keep/cloud/pkg/common/modules"
	"keep/cloud/pkg/requestDispatcher/Router"
	beehiveContext "keep/pkg/util/core/context"
	logger "keep/pkg/util/loggerv1.0.1"
)

func MessageHandler() {
	for {
		select {
		case <-beehiveContext.Done():
			close(Router.SendChan)
			return
		default:
		}

		msg, err := beehiveContext.Receive(modules.RequestDispatcherModule)
		if err != nil {
			logger.Info("receive not Message format message")
			continue
		}

		switch msg.Router.Resource {
		case "$uestc/keep/requestDispatcher/router":
			Router.SendChan <- msg
		}
	}
}
