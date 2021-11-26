package metrics_server

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	for _, master := range []string{"192.168.1.140:6443", "192.168.1.141:6443"} {
		ctx, _ := context.WithTimeout(context.Background(), 5000*time.Millisecond)
		go func(master string) {
			select {
			case <-ctx.Done():
				fmt.Println("获取master:" + master + " 超时...")
			default:
				fmt.Println("开始睡眠....")
				time.Sleep(5 * time.Second)
			}
		}(master)
	}
	for {
		time.Sleep(100 * time.Millisecond)
	}
	//cancel()
}

func TestTimeOut(t *testing.T) {
	NewMetricServerImpl().CheckCadvisorStatus([]string{"192.168.1.140:6443", "192.168.1.141:6443"})
}
