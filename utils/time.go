package utils

import (
	"fmt"
	"time"
)

func GetCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

func ValidateTimestamp(recvWindow int64, timestamp int64) error {
	serverTime := time.Now().UnixMilli()

	if recvWindow > 10000 {
		return fmt.Errorf("recvWindow too large, must be <= 10000")
	}

	if timestamp < serverTime-3000 || timestamp > serverTime+recvWindow {
		return fmt.Errorf("timestamp out of range, must be between %d and %d", serverTime-3000, serverTime+recvWindow)
	}

	return nil
}
