package service

import (
	"fmt"
)

type Secondary struct {
	GRPS string
	HTTP string
}

func GetSecondaries() []Secondary {
	var s = make([]Secondary, 0)

	for i := 1; i <= 2; i++ {
		grpc := fmt.Sprintf("replicated-log-secondary-%d:800%d", i, i)

		port := 7085 + i
		http := fmt.Sprintf("http://replicated-log-secondary-%d:%d", i, port)

		s = append(s, Secondary{GRPS: grpc, HTTP: http})
	}

	return s
}
