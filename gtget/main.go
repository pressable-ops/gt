package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pressable-ops/gt"

	"golang.org/x/crypto/ssh/terminal"
)

var uri = "/v0/transfers/343/"
var method = "GET"
var psk = os.Getenv("GT_PSK")
var debug = false
var debugBody = false

func init() {
	if psk == "" {
		psk = "465d5a695d317e6539543f3835414841692463213049764e7322784b4c"
	}
	flag.StringVar(&uri, "uri", uri, "The URI you want to sign")
	flag.StringVar(&method, "method", method, "GET|DELETE (POST will be used automatically if post key/val pairs are given")
	flag.StringVar(&psk, "psk", psk, "The GT preshared key [env: GT_PSK]")
	flag.BoolVar(&debug, "debug", debug, "Debug requests and response headers")
	flag.BoolVar(&debugBody, "debug-body", debugBody, "Debug request and response bodies")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gtget source code available at https://github.com/pressable-ops/gt\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "  [postkey=postval [postkey2=postval2 [...]]]\n")
	}
}

func maybeNewline(s string) string {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		return s + "\n"
	}
	return s
}

func main() {
	var req *http.Request
	var rsp *http.Response
	var err error
	var URL string
	var signature []byte
	flag.Parse()
	postArgs := flag.Args()
	if len(postArgs) > 0 {
		method = "POST"
	}
	signature, err = gt.Sign(psk, uri, postArgs)
	if err != nil {
		os.Stderr.WriteString(maybeNewline("Error generating signature: " + err.Error()))
		os.Exit(1)
	}
	if strings.Contains(uri, "?") {
		URL = fmt.Sprintf("https://transfer-api.wordpress.com%s&signature=%x", uri, signature)
	} else {
		URL = fmt.Sprintf("https://transfer-api.wordpress.com%s?signature=%x", uri, signature)
	}
	switch strings.ToUpper(method) {
	case "GET":
		req, err = http.NewRequest("GET", URL, nil)
	case "POST":
		var frm = url.Values{}
		for _, arg := range postArgs {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) < 2 {
				log.Fatalf("Incorrect post argument: %s", arg)
			}
			frm.Add(parts[0], parts[1])
		}
		req, err = http.NewRequest("POST", URL, strings.NewReader(frm.Encode()))
		if err != nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	default:
		req, err = http.NewRequest("DELETE", URL, nil)
	}
	if err != nil {
		os.Stderr.WriteString(maybeNewline("Error formulating request: " + err.Error()))
		os.Exit(1)
	}
	rsp, err = httpDo(req)
	if err != nil {
		os.Stderr.WriteString(maybeNewline("Error fetching url: " + err.Error()))
		os.Exit(1)
	}
	if rsp.StatusCode > 299 || rsp.StatusCode < 200 {
		os.Stderr.WriteString(maybeNewline("Error fetching url: " + rsp.Status))
		io.Copy(os.Stdout, rsp.Body)
		os.Stdout.WriteString(maybeNewline(""))
		os.Exit(rsp.StatusCode)
	}
	io.Copy(os.Stdout, rsp.Body)
	os.Stdout.WriteString(maybeNewline(""))
}
