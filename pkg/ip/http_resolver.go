package ip

import (
	"io"
	"net/http"
	"strings"

	"github.com/lajosbencz/dogodns/pkg/config"
)

type httpIPResolver string

func (r httpIPResolver) Resolve() (string, error) {
	resolverUri := string(r)
	resp, err := http.Get(resolverUri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 || strings.Index(resp.Header.Get("Content-Type"), "text/plain") != 0 {
		return "", ErrInvalidResponse
	}
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

func DefaultResolver(cfg config.Config) (Resolver, error) {
	return httpIPResolver(cfg.PIP), nil
}
