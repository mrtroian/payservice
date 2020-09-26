package gateway

import (
	"sync"

	"github.com/mrtroian/payservice/internal/client"

	"github.com/pkg/errors"
)

var (
	collectorPool = sync.Pool{
		New: func() interface{} {
			return new(Collector)
		},
	}
)

type ICollector interface {
	Start(wg *sync.WaitGroup)
	Response() string
	Err() error
	Return()
}

type Collector struct {
	req      client.IRequest
	name     string
	response string
	err      error
}

// NewCollector: return new *Collector
// returns error on invalid arguments
func NewCollector(name string, req client.IRequest) (*Collector, error) {
	if len(name) <= 0 {
		return nil, errors.New("name is invalid")
	}

	c := collectorPool.Get().(*Collector)
	c.req = req
	c.name = name

	return c, nil
}

// (*Collector).Return: return object to the object pool
// usage of (*Collector) after Return invoke results in undefined behavior
func (c *Collector) Return() {
	c.req = nil
	c.name = ""
	c.response = ""
	c.err = nil

	collectorPool.Put(c)
}

func (c *Collector) Name() string {
	return c.name
}

// (*Collector).Err: return error
func (c *Collector) Err() error {
	return c.err
}

// (*Collector).Response: return response
func (c *Collector) Response() string {
	return c.response
}

// (*Collector).Start: start collection from url
func (c *Collector) Start(wg *sync.WaitGroup) {
	resp, err := c.req.Do()

	if err != nil {
		c.err = errors.Wrap(err, "collector."+c.name)
		wg.Done()
		return
	}

	c.response = resp
	wg.Done()
}
