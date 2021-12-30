package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"keep/pkg/util/core/model"
	logger "keep/pkg/util/loggerv1.0.1"
	"net/http"

	"keep/edge/pkg/edgepublisher/tunnel"
)

func httpHandler(writer http.ResponseWriter, request *http.Request) {
	var msg model.Message
	if request.Method != "POST" {
		logger.Error("http请求方法错误")
		fmt.Fprintln(writer, "请求方法错误")
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("http接收数据出错")
		fmt.Fprintln(writer, "接收数据出错")
		return
	}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		logger.Error("http接收数据出错")
		fmt.Fprintln(writer, "接收数据出错")
		return
	}

	edgetunnel.WriteToCloud(&msg)
}

func ListenAndRoute() {
	http.HandleFunc("/healthagent", httpHandler)
	http.ListenAndServe(":8083", nil)
}
