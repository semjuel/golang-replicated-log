package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"replicated_log/basic/api"
	"replicated_log/basic/model"
	"replicated_log/master/utils"
	"sync/atomic"
	"time"
)

func Replicate(message model.Message, target string, iterations *int32) {
	err := utils.IsAlive(target)
	if err != nil {
		retry(message, target)
	}

	err = send(message, target)
	if err != nil {
		retry(message, target)
	}

	atomic.AddInt32(iterations, 1)
}

func retry(message model.Message, target string) {
	var tempDelay time.Duration = 0 // how long to sleep on accept failure

	for {
		err := send(message, target)
		if err == nil {
			break
		}

		if tempDelay == 0 {
			tempDelay = 5 * time.Second
		} else {
			tempDelay *= 2
		}
		if max := 1 * time.Minute; tempDelay > max {
			tempDelay = max
		}
		log.Printf("GRPC send error: %v; retrying in %v", err, tempDelay)
		time.Sleep(tempDelay)
	}
}

func send(message model.Message, target string) error {
	// Run gRPC client.
	log.Printf("Start gRPC client - %s \n", target)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Printf("Did not connect: %s", err)
		return err
	}

	defer func() {
		cerr := conn.Close()
		if err == nil {
			err = cerr
		}
	}()

	c := api.NewLogMessageServiceClient(conn)

	response, err := c.Send(context.Background(), &api.LogMessage{Id: message.Id, Body: message.Text})
	if err != nil {
		log.Printf("Server: %s. Error when calling Send: %s", target, err)
		return err
	}

	log.Printf("Response from server %s: %s", target, response.Body)

	if response.Body != api.SUCCESS {
		return errors.New("failure response from grpc")
	}

	return nil
}
