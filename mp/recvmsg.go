// Package mp provides functions for handling the received messages.
package mp

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"
)

// ValidateSignature is used to validate the signature in request to figure out
// whether the http request come from wechat mp platform.
func ValidateSignature(signature, token, timestamp, nonce string) bool {
	return signature == genSignature(token, timestamp, nonce)
}

func genSignature(token, timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}
