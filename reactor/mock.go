package reactor

import (
	"log"
	"net/http"
)

type MockReactor struct {
}

func (mg *MockReactor) String() string {
	return "mock-reactor"
}

func (mg *MockReactor) Reject(r *http.Request) error {
	log.Println("mockReactor.Reject")
	return nil
}

func (mg *MockReactor) Modify(r *http.Response) error {
	log.Println("mockReactor.ObserveModify")
	return nil
}

func (mg *MockReactor) ObserveError(r *http.Request, resp *http.Response, err error) {
	log.Println("mockReactor.ObserveError")
}
