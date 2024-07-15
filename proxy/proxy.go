package proxy

import (
	"io"
	"mockambo/util"
	"net/http"
	"net/url"
)

func Proxy(req *http.Request, servers []string, server2 string) (*util.Response, error) {
	res := &util.Response{}
	out := util.ReplaceServerURL(req.URL.String(), servers, server2)
	u, err := url.Parse(out)
	if err != nil {
		return res, err
	}
	req2, err := http.NewRequest(req.Method, out, req.Body)
	if err != nil {
		return res, err
	}
	req2.Header = req.Header
	req2.Header.Set("host", u.Hostname())
	req2.Header.Del("Transfer-Encoding")
	req2.Header.Del("Accept-Encoding")
	r, err := httpClient.Do(req2)
	if err != nil {
		return res, err
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
	return res, err
}
