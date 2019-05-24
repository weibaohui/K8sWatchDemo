package utils

import (
	"gotest.tools/assert"
	"testing"
)

func TestCheck(t *testing.T) {

	list := map[string]struct {
		NetWork string
		Address string
		Result  bool
	}{
		"yy": {
			NetWork: "tcp",
			Address: "127.0.0.1:6443",
			Result:  true,
		},
		"xx": {
			NetWork: "udp",
			Address: "127.0.0.1:6443",
			Result:  true,
		},
	}

	for k, v := range list {
		result, err := Check(v.NetWork, v.Address)
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, v.Result, result, "%s 预期 %t,实际%t", k, v.Result, result)
	}
}
