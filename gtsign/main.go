package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/apokalyptik/gt"
)

var uri = "/v0/transfers/343/"
var psk = os.Getenv("GT_PSK")

func init() {
	if psk == "" {
		psk = "465d5a695d317e6539543f3835414841692463213049764e7322784b4c"
	}
	flag.StringVar(&uri, "uri", uri, "The URI you want to sign")
	flag.StringVar(&psk, "psk", psk, "The GT preshared key [env: GT_PSK]")
}

func main() {
	flag.Parse()
	signature, err := gt.Sign(psk, uri, flag.Args())
	if err != nil {
		os.Stderr.WriteString("Error generating signature: " + err.Error())
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			os.Stderr.WriteString("\n")
		}
		os.Exit(1)
	} else {
		os.Stdout.WriteString(fmt.Sprintf("%x", signature))
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			os.Stderr.WriteString("\n")
		}
		os.Exit(0)
	}
}
