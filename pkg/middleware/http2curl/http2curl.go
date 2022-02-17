package http2curl

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// CurlCommand contains exec.Command compatible slice + helpers
type CurlCommand []string

// append appends a string to the CurlCommand
func (c *CurlCommand) append(newSlice ...string) {
	*c = append(*c, newSlice...)
}

// String returns a ready to copy/paste command
func (c *CurlCommand) String() string {
	return strings.Join(*c, " ")
}

func bashEscape(str string) string {
	return `'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}

// GetCurlCommand returns a CurlCommand corresponding to an http.Request
func GetCurlCommand(req *http.Request) (*CurlCommand, error) {
	command := CurlCommand{}

	command.append("curl")

	command.append("-X", bashEscape(req.Method))

	if req.Body != nil {
		var buff bytes.Buffer
		bodyReader, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("getCurlCommand: GetBody error: %w", err)
		}
		_, err = buff.ReadFrom(bodyReader)
		if err != nil {
			return nil, fmt.Errorf("getCurlCommand: buffer read from body erorr: %w", err)
		}
		if len(buff.String()) > 0 {
			bodyEscaped := bashEscape(buff.String())
			command.append("-d", bodyEscaped)
		}
	}

	var keys []string

	for k := range req.Header {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		command.append("-H", appendHeader(req, k))
	}

	command.append(bashEscape(req.URL.String()))

	return &command, nil
}

var redactedHeaders = []string{
	"Authorization",
	"authentication",
	"X-Ka-Service-Admin-Query",
	"X-Ka-Gateway-Secret",
	"Cookie",
	"X-Ka-Fkey",
}

func appendHeader(req *http.Request, k string) string {
	if insensitiveContains(redactedHeaders, k) {
		return bashEscape(fmt.Sprintf("%s: %s", k, "XXXXXX"))
	}
	return bashEscape(fmt.Sprintf("%s: %s", k, strings.Join(req.Header[k], " ")))
}

// insensitiveContains checks if a string slice contains an element, case
// insensitively
func insensitiveContains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
