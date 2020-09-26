package gateway

import (
	"fmt"
	"sync"

	"github.com/mrtroian/payservice/internal/client"

	"github.com/pkg/errors"
)

var manager *GatewayManager

type IManager interface {
	AddGateway(pg PaymentGateway)
	CollectByID(id int) (map[string]string, error)
	Map(f func(k, v string) error) error
}

type GatewayManager struct {
	paymentGateways map[string]string
	mtx             sync.RWMutex
}

type PaymentGateway struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// NewGatewayManager: return singleton *GatewayManager
func NewGatewayManager() *GatewayManager {
	if manager != nil {
		return manager
	}

	manager = &GatewayManager{
		paymentGateways: make(map[string]string),
		mtx:             sync.RWMutex{},
	}

	return manager
}

// (*GatewayManager).AddGateway: add gateway for collecting
// concurrent safe
func (gm *GatewayManager) AddGateway(pg PaymentGateway) {
	gm.mtx.Lock()
	gm.paymentGateways[pg.Name] = pg.URL
	gm.mtx.Unlock()
}

// (*GatewayManager).CollectByID: get all responses from payment gateways
func (gm *GatewayManager) CollectByID(id int) (map[string]string, error) {
	capacity := len(gm.paymentGateways)
	collectors := make([]*Collector, 0, capacity)
	wg := new(sync.WaitGroup)

	err := gm.Map(func(k, v string) error {
		id, err := translateProductID(id, k)

		if err != nil {
			gm.mtx.RUnlock()
			wg.Wait()
			return errors.Wrap(err, "gateway-manager")
		}
		url := fmt.Sprintf("%s%d", v, id)
		req := client.NewRequest(url)
		collector, err := NewCollector(k, req)

		if err != nil {
			gm.mtx.RUnlock()
			wg.Wait()
			return errors.Wrap(err, "gateway-manager: cannot collect")
		}

		wg.Add(1)
		collector.Start(wg)
		collectors = append(collectors, collector)
		return nil
	})
	wg.Wait()

	if err != nil {
		return nil, err
	}

	payMethods := make(map[string]string, capacity)

	for _, c := range collectors {
		if err := c.Err(); err != nil {
			return nil, errors.Wrap(err,
				"gateway-manager: collector returned an error")
		}

		payMethods[c.Name()] = c.Response()
		c.Return()
	}

	return payMethods, nil
}

// (*GatewayManager).Map: apply each key-value pair on f()
// concurrent safe
func (gm *GatewayManager) Map(f func(k, v string) error) error {
	gm.mtx.RLock()
	for k, v := range gm.paymentGateways {
		err := f(k, v)

		if err != nil {
			gm.mtx.RUnlock()
			return err
		}
	}
	gm.mtx.RUnlock()

	return nil
}
