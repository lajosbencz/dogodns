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
