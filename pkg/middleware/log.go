package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/StevenACoffman/commentary/pkg/middleware/http2curl"
)

type LoggingRoundTripper struct {
	next   http.RoundTripper
	logger io.Writer
}

func NewLoggingRoundTripper(
	next http.RoundTripper,
	w io.Writer,
) *LoggingRoundTripper {
	return &LoggingRoundTripper{
		next:   next,
		logger: w,
	}
}

func (lrt *LoggingRoundTripper) RoundTrip(
	req *http.Request,
) (resp *http.Response, err error) {
	defer func(begin time.Time) {
		var msg string
		body, getResponseBodyErr := GetResponseBody(resp)
		if getResponseBodyErr != nil {
			fmt.Println("unable to get response Body", getResponseBodyErr)
		}
		gotHTTPErr := resp != nil && (resp.StatusCode < 200 || resp.StatusCode >= 300)
		graphqlErr := GetGraphQLErrors(resp)
		// only log when there was a problem, so early return here when happy
		if !gotHTTPErr && graphqlErr == nil {
			return
		}
		msg = fmt.Sprintf(
			"method=%s host=%s path=%s status_code=%d took=%s\n",
			req.Method,
			req.URL.Host,
			req.URL.Path,
			resp.StatusCode,
			time.Since(begin),
		)
		if err != nil {
			fmt.Fprintf(lrt.logger, "%s : %+v\n", msg, err)
		} else {
			fmt.Fprintf(lrt.logger, "%s\n", msg)
		}
		command, _ := http2curl.GetCurlCommand(req)
		fmt.Println(command)

		fmt.Println("body:\n", body, "\nend body\n ")
		if graphqlErr != nil {
			fmt.Printf("graphql errors:%+v\n", graphqlErr)
		}
	}(time.Now())

	return lrt.next.RoundTrip(req)
}

// GetResponseBody will read the response body without clobbering it
// so it can be re-read elsewhere
func GetResponseBody(r *http.Response) (string, error) {
	if r == nil || r.Body == nil {
		return "", nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return string(body), nil
}

func GetGraphQLErrors(r *http.Response) error {
	if r == nil || r.Body == nil {
		return nil
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	var dataAndErrors response
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&dataAndErrors)
	if err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if len(dataAndErrors.Errors) > 0 {
		return dataAndErrors.Errors
	}
	return nil
}

type response struct {
	Data   interface{}   `json:"data"`
	Errors gqlerror.List `json:"errors"`
}
