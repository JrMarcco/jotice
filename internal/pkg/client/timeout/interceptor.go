package timeout

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const keyTimeout = "key_timeout"

func InjectInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if deadline, ok := ctx.Deadline(); ok {
			ctx = metadata.AppendToOutgoingContext(ctx, keyTimeout, fmt.Sprintf("%d", deadline.UnixMilli()))
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
