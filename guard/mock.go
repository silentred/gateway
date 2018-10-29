package guard

import (
	"net/http"
)

type MockGuard struct {
}

func (mg *MockGuard) Reject(r *http.Request) error {
	return nil
}

func (mg *MockGuard) String() string {
	return "mock-guard"
}
