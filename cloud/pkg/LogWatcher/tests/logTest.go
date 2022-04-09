package main

import (
	"github.com/UESTC-KEEP/keep/cloud/pkg/LogWatcher"
)

func main() {
	LogWatcher.GetAndPusherKafka()
}
