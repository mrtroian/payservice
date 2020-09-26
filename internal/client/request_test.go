package client

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

type mockClient struct{}
type mockioReader struct {
	*bytes.Buffer
}

func (_ mockioReader) Close() error {
	return nil
}

func (_ mockClient) Do(req *http.Request) (*http.Response, error) {
	r := NewResponse()
	js, err := json.Marshal(r)

	if err != nil {
		log.Fatal(err)
	}

	m := mockioReader{bytes.NewBuffer(js)}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       m,
	}, nil
}

func TestRequest(t *testing.T) {
	r := NewRequest("")
	httpClient = mockClient{}

	if got, err := r.Do(); err == nil {
		t.Errorf(`(*Request).Do() = %s; want err`, got)
	}
}
