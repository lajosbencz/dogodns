package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func fail(msg any) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(-1)
}

var ErrDOInternal = fmt.Errorf("DO API Internal error")

func errDO(err error) {
	if err != ErrDOInternal {
		fail(err)
	}
}

func getPublicIP(url string, ipv6 bool) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	return string(bytes), nil
}
