package gateway

import (
	"fmt"
	"sync"
)

var (
	apDict = make(map[string]string)
	gpDict = make(map[string]string)
	ppDict = make(map[string]string)
	apmtx  = sync.RWMutex{}
	gpmtx  = sync.RWMutex{}
	ppmtx  = sync.RWMutex{}
)

// translateProductID: translate productID to ProductID in payment gateway
func translateProductID(id int, service string) (int, error) {
	switch service {
	case "ApplePay":
		return id, nil

	case "GooglePay":
		return id, nil

	case "PayPal":
		return id, nil

	}

	return -1, fmt.Errorf("translate: unknown payment service: %s", service)
}
