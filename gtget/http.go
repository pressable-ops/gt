package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func maybeDebugResponse(rsp *http.Response) {
	if debug {
		if buf, err := httputil.DumpResponse(rsp, debugBody); err != nil {
			log.Println("Error dumping the response:", err.Error())
		} else {
			log.Println(string(buf))
		}
	}
}

func maybeDebugRequest(req *http.Request) {
	if debug {
		if buf, err := httputil.DumpRequestOut(req, debugBody); err != nil {
			log.Println("Error dumping the request:", err.Error())
		} else {
			log.Println(string(buf))
		}
	}
}

func httpDo(req *http.Request) (*http.Response, error) {
	maybeDebugRequest(req)
	rsp, err := http.DefaultClient.Do(req)
	maybeDebugResponse(rsp)
	return rsp, err
}
