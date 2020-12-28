package utils

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func IsAlive(target string) error {
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	url := fmt.Sprintf("%s/health", target)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("wrong response status code from `%s`", target))
	}

	return nil
}
