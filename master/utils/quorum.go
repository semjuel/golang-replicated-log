package utils

import (
	"errors"
	"log"
	"replicated_log/master/service"
	"sync"
	"sync/atomic"
)
func HasQuorum() error {
	s := service.GetSecondaries()
	var alive int32 = 1
	wg := sync.WaitGroup{}
	for _, v := range s {
		wg.Add(1)
		go func(v service.Secondary, alive *int32) {
			err := IsAlive(v.HTTP)
			log.Printf("Health error = %v", err)
			if err == nil {
				atomic.AddInt32(alive, 1)
			}
			wg.Done()
		}(v, &alive)
	}
	wg.Wait()

	log.Printf("Alive nodes %v", alive)
	if (len(s) + 1) / 2 >= int(alive) {
		return errors.New("no quorum")
	}
	return nil
}
