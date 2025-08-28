package sudir

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

const DefaultLength = 32

type CodeVerifier struct {
	Value string
}

func NewCodeVerifier(value string) (*CodeVerifier, error) {
	cv := &CodeVerifier{Value: value}

	if cv.Value != "" {
		return cv, nil
	}

	buf, err := cv.random(DefaultLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %v", err)
	}

	cv.Value = cv.encode(buf)
	return cv, nil
}

func (cv *CodeVerifier) CodeChallengeS256() string {
	h := sha256.New()
	h.Write([]byte(cv.Value))
	return cv.encode(h.Sum(nil))
}

func (cv *CodeVerifier) encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)
	return encoded
}

// https://tools.ietf.org/html/rfc7636#section-4.1
func (cv *CodeVerifier) random(length int) ([]byte, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const csLen = byte(len(charset))
	output := make([]byte, 0, length)
	for {
		buf := make([]byte, length)
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			return nil, fmt.Errorf("failed to read random bytes: %v", err)
		}
		for _, b := range buf {
			if b < (csLen * 4) {
				output = append(output, charset[b%csLen])

				if len(output) == length {
					return output, nil
				}
			}
		}
	}
}
