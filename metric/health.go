package metric

import "net/http"

func Health(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.WriteHeader(404)
	default:
		w.WriteHeader(200)
	}
}
