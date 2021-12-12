package create

import (
	"net/http"
	"testing"
)

func TestConnect(t *testing.T) {
	http.HandleFunc("/create", Connect)
	http.ListenAndServe(":8080", nil)

}
func Connect(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("nihao"))
}
