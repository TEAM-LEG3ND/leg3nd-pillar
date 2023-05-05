package config

import (
	"fmt"
	"os"
)

func GetEnv(key string) (*string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return nil, fmt.Errorf("config for key %s do not exist", key)
	} else {
		return &val, nil
	}
}
