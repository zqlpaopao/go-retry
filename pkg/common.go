package pkg

import "time"

const (
	RetryCount    = 3 //RetryCount retry times
	PoolCount     = 3
	RetryInterval = 3 * time.Second //RetryInterval retry interval

)
