package client

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/JrMarcco/easy-kit/slice"
	"google.golang.org/grpc/balancer"
)

type ContextKey string

const (
	readWeightStr  = "read_weight"
	writeWeightStr = "write_weight"
	groupStr       = "group"
	nodeStr        = "node"
)

const RequestType ContextKey = "request_type"

type groupKey struct{}

type rwServiceNode struct {
	mu                   sync.RWMutex
	subConn              balancer.SubConn
	readWeight           int32
	currentReadWeight    int32
	efficientReadWeight  int32
	writeWeight          int32
	currentWriteWeight   int32
	efficientWriteWeight int32
	group                string
}

var _ balancer.Picker = (*RWBalancer)(nil)

type RWBalancer struct {
	nodes []*rwServiceNode
}

func (r *RWBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(r.nodes) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	// select candidate nodes
	candidates := slice.FilterMap(r.nodes, func(_ int, src *rwServiceNode) (*rwServiceNode, bool) {
		src.mu.RLock()
		nodeGroup := src.group
		src.mu.RUnlock()

		return src, r.getGroup(info.Ctx) == nodeGroup
	})

	var totalWeight int32
	var selectedNode *rwServiceNode

	ctx := info.Ctx
	isWrite := r.isWrite(ctx)

	for _, node := range candidates {
		node.mu.Lock()
		if isWrite {
			totalWeight += node.efficientWriteWeight
			node.currentWriteWeight += node.efficientWriteWeight
			if selectedNode == nil || selectedNode.currentWriteWeight < node.currentWriteWeight {
				selectedNode = node
			}
		} else {
			totalWeight += node.efficientReadWeight
			node.currentReadWeight += node.efficientReadWeight
			if selectedNode == nil || selectedNode.currentReadWeight < node.currentReadWeight {
				selectedNode = node
			}
		}
		node.mu.Unlock()
	}

	if selectedNode == nil {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	selectedNode.mu.Lock()
	if isWrite {
		selectedNode.currentWriteWeight -= totalWeight
	} else {
		selectedNode.currentReadWeight -= totalWeight
	}
	selectedNode.mu.Unlock()

	return balancer.PickResult{
		SubConn: selectedNode.subConn,
		Done: func(info balancer.DoneInfo) {
			selectedNode.mu.Lock()
			defer selectedNode.mu.Unlock()

			isDecrementErr := info.Err == nil && (errors.Is(info.Err, context.DeadlineExceeded) || errors.Is(info.Err, io.EOF))
			if isWrite {
				if isDecrementErr && selectedNode.currentWriteWeight > 0 {
					selectedNode.currentWriteWeight--
				} else if info.Err == nil {
					selectedNode.currentWriteWeight++
				}
			} else {
				if isDecrementErr && selectedNode.currentReadWeight > 0 {
					selectedNode.currentReadWeight--
				} else if info.Err == nil {
					selectedNode.currentReadWeight++
				}
			}
		},
	}, nil
}

func (r *RWBalancer) getGroup(ctx context.Context) string {
	val := ctx.Value(groupStr)
	if val == nil {
		return ""
	}

	if valString, ok := val.(string); ok {
		return valString
	}

	return ""
}

func (r *RWBalancer) isWrite(ctx context.Context) bool {
	val := ctx.Value(RequestType)
	if val == nil {
		return false
	}

	if valInt, ok := val.(int); ok {
		return valInt == 1
	}

	return false
}
