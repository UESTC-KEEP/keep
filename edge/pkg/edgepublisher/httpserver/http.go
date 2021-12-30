package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"keep/pkg/util/core/model"
	"net/http"

	"keep/edge/pkg/edgepublisher/tunnel"
)

func httpHandler(writer http.ResponseWriter, request *http.Request) {
	var msg *model.Message
	if request.Method != "POST" {
		fmt.Fprintln(writer, "请求方法错误")
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintln(writer, "接收数据出错")
		return
	}
	err = json.Unmarshal(body, msg)
	if err != nil {
		fmt.Fprintln(writer, "接收数据出错")
		return
	}

	edgetunnel.WriteToCloud(msg)
}

func ListenAndRoute() {
	http.HandleFunc("/", httpHandler)
	http.ListenAndServe(":8083", nil)
}
