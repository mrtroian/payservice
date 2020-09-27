package methods_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"

	"github.com/mrtroian/payservice/internal/client"
	"github.com/mrtroian/payservice/internal/configuration"
)

var (
	config      *configuration.Config
	url         string
	httpsClient = client.NewClient()
)

func init() {
	var err error

	config, err = configuration.GetConfig()

	if err != nil {
		log.Fatalln(err)
	}
	url = fmt.Sprintf(
		"https://%s:%d%s",
		config.Host,
		config.Port,
		config.Endpoint,
	)
}

func TestNoID(t *testing.T) {
	req, err := http.NewRequest("GET", url, &bytes.Buffer{})

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	resp, err := httpsClient.Do(req)

	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fail()
		return
	}
	t.Log(resp)
}

func TestRangeIntID(t *testing.T) {
	wg := new(sync.WaitGroup)

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, id int) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", url, id), &bytes.Buffer{})

			if err != nil {
				t.Log(err)
				t.Fail()
				return
			}

			resp, err := httpsClient.Do(req)

			if err != nil {
				t.Log(err)
				t.Fail()
				return
			}

			if resp.StatusCode != http.StatusOK {
				t.Fail()
				return
			}
			t.Log(resp)
			wg.Done()
		}(wg, i)
	}
	wg.Wait()
}
