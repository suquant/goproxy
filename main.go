package main

import (
	"log"
	"net/http"
	"encoding/json"
	"strings"
	"fmt"
	"io"
	"os"
)

type PreparedHttpRequest struct {
	Method string `json:"method"`
	Url string `json:"url"`
	Body string `json:"body"`
	Header http.Header `json:"header"`
}

func copyHeaders(dst, src http.Header) {
	for k, _ := range dst {
		dst.Del(k)
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var preparedRequest PreparedHttpRequest
	err := decoder.Decode(&preparedRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("JSON decode error: %s", err.Error()), 500)
		return
	}
	
	reqBodyReader := strings.NewReader(preparedRequest.Body)
	origReq, err := http.NewRequest(preparedRequest.Method, preparedRequest.Url, reqBodyReader)
	if preparedRequest.Header != nil {
		copyHeaders(origReq.Header, preparedRequest.Header)	
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Request initialize error: %s", err.Error()), 500)
		return
	}
	
	client := &http.Client{}
	clientResp, err := client.Do(origReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Remote response error: %s", err.Error()), 500)
		return
	}
	
	w.WriteHeader(clientResp.StatusCode)
	io.Copy(w, clientResp.Body)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
