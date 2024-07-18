package proxy

import (
	"io"
	"mockambo/exceptions"
	"mockambo/util"
	"net/http"
	"net/url"
)

// Proxy proxies a request against another server.
// We need the request, the list of declared servers, and which server we want to use as target
func Proxy(req *http.Request, servers []string, server2 string) (*util.Response, error) {
	res := &util.Response{}
	out := util.ReplaceServerURL(req.URL.String(), servers, server2)
	u, err := url.Parse(out)
	if err != nil {
		return res, exceptions.Wrap("parse_url", err)
	}
	req2, err := http.NewRequest(req.Method, out, req.Body)
	if err != nil {
		return res, exceptions.Wrap("new_request", err)
	}
	req2.Header = req.Header
	req2.Header.Set("host", u.Hostname())
	// Removing transport and encoding specific headers
	req2.Header.Del("Transfer-Encoding")
	req2.Header.Del("Accept-Encoding")
	r, err := httpClient.Do(req2)
	if err != nil {
		return res, exceptions.Wrap("perform_request", err)
	}
	defer func() {
		if r != nil && r.Body != nil {
			_ = r.Body.Close()
		}
	}()

	res.Status = r.StatusCode
	res.ContentType = r.Header.Get(util.HeaderContentType)
	res.Headers = r.Header
	res.Payload, err = io.ReadAll(r.Body)
	return res, exceptions.Wrap("read_body", err)
}
