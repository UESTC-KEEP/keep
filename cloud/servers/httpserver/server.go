package httpserver	

import (
	"fmt"
	"keep/conf"
	"net/http"

 )
  
 type HelloHandlerStruct struct {
	content string
 }
  
 func (handler *HelloHandlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, handler.content)
 }
  
 func StartHTTPServer() {
	port := conf.GetStringConfig("httpport")
	http.Handle("/", &HelloHandlerStruct{content: "Hello World"})
	http.ListenAndServe(":" + port, nil)
 }