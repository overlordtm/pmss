package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type staticTokenAuthHandler struct {
	token string
}

func NewStaticTokenAuthHandler(token string) staticTokenAuthHandler {
	return staticTokenAuthHandler{token}
}

func (s staticTokenAuthHandler) IsAuthorized(ctx context.Context, req *http.Request) (bool, error) {

	parts := strings.SplitN(req.Header.Get("Authorization"), " ", 2)

	if parts[0] != "Token" {
		return false, fmt.Errorf("invalid authorization header: %s", req.Header.Get("Authorization"))
	}

	return parts[1] == s.token, nil
}
