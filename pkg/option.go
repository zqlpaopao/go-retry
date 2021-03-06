package pkg

import (
	"time"
)

type Option interface {
	apply(opt *option)
}

type DelayTypeFunc func(executeCount uint) time.Duration

//RetryableFunc Method of performing retry
type RetryableFunc func() bool

//onRetryCallbackFun Function to execute each retry
type onRetryCallbackFun func(uint)

//onCompleteCallbackFun Retry the function that completed execution
type onCompleteCallbackFun func(uint, bool, ...interface{})

type OpFunc func(*option)

//retry options
type option struct {
	asyncTag      bool
	retryCount    uint
	poolCount     int
	retryInterval time.Duration
	delayType     DelayTypeFunc
	pool          *pool
}

//apply assignment function entity
func (o OpFunc) apply(opt *option) {
	o(opt)
}

//NewOption make option
func NewOption(opt ...Option) *option {
	o := &option{
		asyncTag:      false,
		retryCount:    RetryCount,
		retryInterval: RetryInterval,
		delayType:     WithDefaultDelayType,
		poolCount:     PoolCount,
		pool:          &pool{},
		//poolCount:  NewPool(PoolCount),
	}
	return o.WithOptions(opt...)
}

//clone  new object
func (o *option) clone() *option {
	cp := *o
	return &cp
}

//WithOptions Execute assignment function entity
func (o option) WithOptions(opt ...Option) *option {
	c := o.clone()
	for _, v := range opt {
		v.apply(c)
	}
	if !c.asyncTag {
		return c
	}
	c.pool = NewPool(c.poolCount)
	go c.pool.Run()
	return c
}

//WithRetryCount Set the number of retries. The default is retryCount = 3
func WithRetryCount(retryCount uint) OpFunc {
	return func(o *option) {
		o.retryCount = retryCount
	}
}

//WithAsyncTag is use go execute
func WithAsyncTag(asyncTag bool) OpFunc {
	return func(o *option) {
		o.asyncTag = asyncTag
	}
}

//WithPoolCount go pool number
func WithPoolCount(poolCount int) OpFunc {
	return func(o *option) {
		o.poolCount = poolCount
	}
}

//WithRetryInterval Set the retry interval. The default is retryInterval = 3
func WithRetryInterval(retryInterval time.Duration) OpFunc {
	return func(o *option) {
		o.retryInterval = retryInterval
	}
}

//WithDelayType Set the growth mode of retry time. The default value is WithDefaultDelayType
func WithDelayType(delayType DelayTypeFunc) OpFunc {
	return func(o *option) {
		o.delayType = delayType
	}
}

//WithDefaultDelayType default
func WithDefaultDelayType(executeCount uint) time.Duration {
	return 1
}

//WithBackOffDelayType  Incremental retry interval, exponential
//2-4-8-16-32-...
func WithBackOffDelayType(executeCount uint) time.Duration {
	return 1 << (executeCount - 1)
}
