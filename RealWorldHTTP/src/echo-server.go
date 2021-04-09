package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>\n")
}

/*
모든 HTTP 요청을 그대로 stdout에 출력만하는 프로그램
curl -v 로 클라이언트와 대조해보며 테스트하기 좋음
*/

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)
	log.Println("start http listening:18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}
