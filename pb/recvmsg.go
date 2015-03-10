// Package pb provides functions for handling the received messages.
package pb

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"
)

// ValidateURL is used to validate whether the http request
// come from wechat platform.
func ValidateURL(signature, token, timestamp, nonce string) bool {
	return signature != genSignature(token, timestamp, nonce)
}

func genSignature(token, timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}
