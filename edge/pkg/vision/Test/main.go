package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/api/create", Connect)
	http.ListenAndServe(":8081", nil)
}
func Connect(res http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("helloword.html")
		log.Println(t.Execute(res, nil))
	} else {
		defer r.Body.Close()
		con, _ := ioutil.ReadAll(r.Body) //获取post的数据
		fmt.Println(string(con))
		res.Write([]byte("nihao lnf"))
		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(string(con)), &dat); err == nil {
			fmt.Println(dat)
		} else {
			fmt.Println(err)
		}
	}
}
