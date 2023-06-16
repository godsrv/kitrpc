package middleware

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/juju/ratelimit"
)

type Bucket struct {
	*ratelimit.Bucket
}

func (b Bucket) Limiting() bool {
	return b.TakeAvailable(1) > 0
}

func TokenBucketLimitter(bkt Bucket) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Limiting() {
				return nil, errors.New("rate limit exceed!")
			}
			// 如果成功就到下一层
			return e(ctx, request)
		}
	}
}
