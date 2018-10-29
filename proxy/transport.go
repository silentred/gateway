package proxy

import "net/http"

type ErrTransport struct {
	trans    http.RoundTripper
	LastErr  error
	Response *http.Response
}

func NewErrTransport(r http.RoundTripper) *ErrTransport {
	return &ErrTransport{r, nil, nil}
}

func (t *ErrTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Response, t.LastErr = t.trans.RoundTrip(req)
	return t.Response, t.LastErr
}
