package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")

	if auth == "" {
		return "", errors.New("no Authorization header present")
	}

	token := strings.Split(auth, " ")[1]

	return token, nil
}
