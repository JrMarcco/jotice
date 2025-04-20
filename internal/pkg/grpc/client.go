package grpc

import "github.com/JrMarcco/easy-kit/xmap"

type Client[T any] struct {
	clientM xmap.Map[string, T]
}
