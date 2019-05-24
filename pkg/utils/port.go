package utils

import (
	"fmt"
	"net"
	"time"
)

func Check(network, address string) (bool, error) {
	_, e := net.DialTimeout(network, address, time.Second*3)
	if e != nil {
		fmt.Println(e.Error())
		return false, e
	}
	return true, nil
}
